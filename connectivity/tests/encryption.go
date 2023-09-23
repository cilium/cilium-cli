// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cilium/cilium/pkg/defaults"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cilium/cilium-cli/connectivity/check"
)

type requestType int

const (
	requestHTTP requestType = iota
	requestICMPEcho
)

// getInterNodeIface determines on which netdev iface to capture pkts.
// We run "ip route get $DST_IP" from the client pod's node to see to
// which interface the traffic is routed to. Additionally, we translate
// the interface name to the tunneling interface name, if the route goes
// through "cilium_host" and tunneling is enabled.
func getInterNodeIface(ctx context.Context, t *check.Test, clientHost *check.Pod, dstIP string) string {
	cmd := []string{
		"/bin/sh", "-c",
		fmt.Sprintf("ip -o route get %s | grep -oE 'dev [^ ]*' | cut -d' ' -f2",
			dstIP),
	}
	t.Debugf("Running %s", strings.Join(cmd, " "))
	dev, err := clientHost.K8sClient.ExecInPod(ctx, clientHost.Pod.Namespace,
		clientHost.Pod.Name, "", cmd)
	if err != nil {
		t.Fatalf("Failed to get IP route: %s", err)
	}

	device := strings.TrimRight(dev.String(), "\n\r")

	// When tunneling is enabled, and the traffic is routed to the cilium IP space
	// we want to capture on the tunnel interface.
	if tunnelFeat, ok := t.Context().Feature(check.FeatureTunnel); ok &&
		tunnelFeat.Enabled && device == defaults.HostDevice {
		return "cilium_" + tunnelFeat.Mode // E.g. cilium_vxlan
	}

	return device
}

// getSourceAddress determines the source IP address we want to use for
// capturing packet. If direct routing is used, the source IP is the client IP.
func getSourceAddress(ctx context.Context, t *check.Test, client,
	clientHost *check.Pod, ipFam check.IPFamily, dstIP string,
) string {
	if tunnelStatus, ok := t.Context().Feature(check.FeatureTunnel); ok &&
		!tunnelStatus.Enabled {
		return client.Address(ipFam)
	}

	cmd := []string{
		"/bin/sh", "-c",
		fmt.Sprintf("ip -o route get %s | grep -oE 'src [^ ]*' | cut -d' ' -f2",
			dstIP),
	}
	t.Debugf("Running %s", strings.Join(cmd, " "))
	srcIP, err := clientHost.K8sClient.ExecInPod(ctx, clientHost.Pod.Namespace,
		clientHost.Pod.Name, "", cmd)
	if err != nil {
		t.Fatalf("Failed to get IP route: %s", err)
	}

	return strings.TrimRight(srcIP.String(), "\n\r")
}

// PodToPodEncryption is a test case which checks the following:
//   - There is a connectivity between pods on different nodes when any
//     encryption mode is on (either WireGuard or IPsec).
//   - No unencrypted packet is leaked.
//
// The checks are implemented by curl'ing a server pod from a client pod, and
// then inspecting tcpdump captures from the client pod's node.
func PodToPodEncryption() check.Scenario {
	return &podToPodEncryption{}
}

type podToPodEncryption struct{}

func (s *podToPodEncryption) Name() string {
	return "pod-to-pod-encryption"
}

func (s *podToPodEncryption) Run(ctx context.Context, t *check.Test) {
	ct := t.Context()
	client := ct.RandomClientPod()

	var server check.Pod
	for _, pod := range ct.EchoPods() {
		// Make sure that the server pod is on another node than client
		if pod.Pod.Status.HostIP != client.Pod.Status.HostIP {
			server = pod
			break
		}
	}

	// clientHost is a pod running on the same node as the client pod, just in
	// the host netns.
	clientHost := ct.HostNetNSPodsByNode()[client.Pod.Spec.NodeName]

	t.ForEachIPFamily(func(ipFam check.IPFamily) {
		testNoTrafficLeak(ctx, t, client, &server, &clientHost, requestHTTP, ipFam, false)
	})
}

// PodToPodEncryption is a test case which checks the following:
//   - There is a connectivity between pods on different nodes when any
//     encryption mode is on (either WireGuard or IPsec).
//   - No unencrypted packet is leaked.
//
// The checks are implemented by curl'ing a server pod from a client pod, and
// then inspecting tcpdump captures from the client pod's node.
func PodToPodStrictEncryption() check.Scenario {
	return &podToPodStrictEncryption{}
}

type podToPodStrictEncryption struct{}

func (s *podToPodStrictEncryption) Name() string {
	return "pod-to-pod-strict-encryption"
}

func (s *podToPodStrictEncryption) Run(ctx context.Context, t *check.Test) {
	ct := t.Context()
	client := ct.RandomClientPod()

	var server check.Pod
	for _, pod := range ct.EchoPods() {
		// Make sure that the server pod is on another node than client
		if pod.Pod.Status.HostIP != client.Pod.Status.HostIP {
			server = pod
			break
		}
	}

	// clientHost is a pod running on the same node as the client pod, just in
	// the host netns.
	clientHost := ct.HostNetNSPodsByNode()[client.Pod.Spec.NodeName]

	t.ForEachIPFamily(func(ipFam check.IPFamily) {
		testNoTrafficLeakStrict(ctx, t, client, &server, &clientHost, requestHTTP, ipFam)
	})
}

// startTcpdump starts tcpdump in the background, and returns a cancel function
// to stop it, and a channel which is closed when tcpdump has exited.
// It writes captured pkts to /tmp/$TEST_NAME.pcap.
func startTcpdump(ctx context.Context, t *check.Test,
	client, server, clientHost *check.Pod, reqType requestType, ipFam check.IPFamily,
) (context.CancelFunc, chan struct{}) {
	dstAddr := server.Address(ipFam)
	iface := getInterNodeIface(ctx, t, clientHost, dstAddr)
	srcAddr := getSourceAddress(ctx, t, client, clientHost, ipFam, dstAddr)

	bgStdout := &safeBuffer{}
	bgStderr := &safeBuffer{}
	bgExited := make(chan struct{})
	killCmdCtx, killCmd := context.WithCancel(context.Background())
	// Start kubectl exec in bg (=goroutine)
	go func() {
		protoFilter := ""
		switch reqType {
		case requestHTTP:
			protoFilter = "tcp"
		case requestICMPEcho:
			protoFilter = "icmp"
			if ipFam == check.IPFamilyV6 {
				protoFilter = "icmp6"
			}
		}
		// Run tcpdump with -w instead of directly printing captured pkts. This
		// is to avoid a race after sending ^C (triggered by bgCancel()) which
		// might terminate the tcpdump process before it gets a chance to dump
		// its captures.
		cmd := []string{
			"tcpdump", "-i", iface, "--immediate-mode", "-w", fmt.Sprintf("/tmp/%s.pcap", t.Name()),
			// Capture egress traffic.
			// Unfortunately, we cannot use "host %s and host %s" filter here,
			// as IPsec recirculates replies to the iface netdev, which would
			// make tcpdump to capture the pkts (false positive).
			fmt.Sprintf("src host %s and dst host %s and %s", srcAddr, dstAddr, protoFilter),
		}
		t.Debugf("Running in bg: %s", strings.Join(cmd, " "))
		err := clientHost.K8sClient.ExecInPodWithWriters(ctx, killCmdCtx,
			clientHost.Pod.Namespace, clientHost.Pod.Name, "", cmd, bgStdout, bgStderr)
		if err != nil && !errors.Is(err, context.Canceled) {
			t.Fatalf("Failed to execute tcpdump: %s", err)
		}
		close(bgExited)
	}()

	// Wait until tcpdump is ready to capture pkts
	timeout := time.After(5 * time.Second)
	for found := false; !found; {
		select {
		case <-timeout:
			t.Fatalf("Failed to wait for tcpdump to be ready")
		default:
			line, err := bgStdout.ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				t.Fatalf("Failed to read kubectl exec's stdout: %s", err)
			}
			if strings.Contains(line, fmt.Sprintf("listening on %s", iface)) {
				found = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	return killCmd, bgExited
}

// checkPcapForLeak checks whether there is any unencrypted pkt captured in the
// pcap file. If so, it fails the test.
func checkPcapForLeak(ctx context.Context, t *check.Test, clientHost *check.Pod) {
	// Redirect stderr to /dev/null, as tcpdump logs to stderr, and ExecInPod
	// will return an error if any char is written to stderr. Anyway, the count
	// is written to stdout.
	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("tcpdump -r /tmp/%s.pcap --count 2>/dev/null", t.Name())}
	count, err := clientHost.K8sClient.ExecInPod(ctx, clientHost.Pod.Namespace, clientHost.Pod.Name, "", cmd)
	if err != nil {
		t.Fatalf("Failed to retrieve tcpdump pkt count: %s", err)
	}
	if !strings.HasPrefix(count.String(), "0 packets") {
		t.Failf("Captured unencrypted pkt (count=%s)", strings.TrimRight(count.String(), "\n\r"))

		// If debug mode is enabled, dump the captured pkts
		if t.Context().Params().Debug {
			cmd := []string{"/bin/sh", "-c", fmt.Sprintf("tcpdump -r /tmp/%s.pcap 2>/dev/null", t.Name())}
			out, err := clientHost.K8sClient.ExecInPod(ctx, clientHost.Pod.Namespace, clientHost.Pod.Name, "", cmd)
			if err != nil {
				t.Fatalf("Failed to retrieve tcpdump output: %s", err)
			}
			t.Debugf("Captured pkts:\n%s", out.String())
		}
	}
}

func testNoTrafficLeakStrict(ctx context.Context, t *check.Test,
	client, server, clientHost *check.Pod, reqType requestType, ipFam check.IPFamily,
) {
	deleteCES := func(endpointName string) {
		cesList, err := clientHost.K8sClient.CiliumClientset.CiliumV2alpha1().CiliumEndpointSlices().List(ctx, metav1.ListOptions{})
		if err != nil {
			t.Fatalf("Failed to list CiliumEndpointSlices: %s", err)
		}
		var cesToDelete string
		for _, ces := range cesList.Items {
			for _, ep := range ces.Endpoints {
				if ep.Name == endpointName {
					cesToDelete = ces.Name
					break
				}
			}
		}
		if cesToDelete == "" {
			t.Fatalf("Failed to find CiliumEndpointSlice for pod %s", server.Pod.Name)
		}
		if err := clientHost.K8sClient.CiliumClientset.CiliumV2alpha1().CiliumEndpointSlices().Delete(ctx, cesToDelete, metav1.DeleteOptions{}); err != nil {
			t.Fatalf("Failed to delete CiliumEndpointSlice %s: %s", cesToDelete, err)
		}
	}
	setCiliumOperatorScale := func(replicas int32) int32 {
		scale, err := clientHost.K8sClient.Clientset.AppsV1().Deployments(t.Context().Params().CiliumNamespace).GetScale(ctx, "cilium-operator", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Failed to get cilium-operator scale: %s", err)
		}
		savedReplicas := scale.Spec.Replicas
		scale.Spec.Replicas = replicas
		if _, err := clientHost.K8sClient.Clientset.AppsV1().Deployments(t.Context().Params().CiliumNamespace).UpdateScale(ctx, "cilium-operator", scale, metav1.UpdateOptions{}); err != nil {
			t.Fatalf("Failed to scale cilium-operator: %s", err)
		}
		return savedReplicas
	}
	waitForIPCacheEntry := func(clientPod, dstPod *check.Pod) {
		timeout := time.After(20 * time.Second)
		var ciliumHostPodName string
		for _, pod := range t.Context().CiliumPods() {
			if pod.NodeName() == clientPod.Pod.Spec.NodeName {
				ciliumHostPodName = pod.Pod.Name
				break
			}
		}
		for found := false; !found; {
			select {
			case <-timeout:
				t.Fatalf("Failed to wait for ipcache to contain pod %s's IP", dstPod.Pod.Name)
			default:
				cmd := []string{"/bin/sh", "-c", "cilium bpf ipcache list"}
				out, err := clientPod.K8sClient.ExecInPod(ctx, t.Context().Params().CiliumNamespace, ciliumHostPodName, "", cmd)
				if err != nil {
					t.Fatalf("Failed to retrieve ipcache output: %s", err)
				}
				if strings.Contains(out.String(), dstPod.Pod.Status.PodIP) {
					found = true
					break
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	}

	// Disable endpoint propagation by scaling down the cilium-operator in the cilium namespace
	savedScale := setCiliumOperatorScale(0)

	// Delete CES of the server pod
	deleteCES(server.Pod.Name)
	deleteCES(client.Pod.Name)

	// Run the test
	testNoTrafficLeak(ctx, t, client, server, clientHost, reqType, ipFam, true)

	// Restore the cilium-operator scale
	_ = setCiliumOperatorScale(savedScale)

	// wait for the ipcache to contain the server pod's IP
	waitForIPCacheEntry(client, server)

	// Run the test
	testNoTrafficLeak(ctx, t, client, server, clientHost, reqType, ipFam, false)
}

func testNoTrafficLeak(ctx context.Context, t *check.Test, client, server, clientHost *check.Pod,
	reqType requestType, ipFam check.IPFamily, expectFail bool,
) {
	// Setup
	killCmd, bgExited := startTcpdump(ctx, t, client, server, clientHost, reqType, ipFam)

	var cmd []string
	// Run the test
	switch reqType {
	case requestHTTP:
		// Curl the server from the client to generate some traffic
		cmd = t.Context().CurlCommand(server, ipFam)
	case requestICMPEcho:
		// Ping the server from the client to generate some traffic
		cmd = t.Context().PingCommand(server, ipFam)
	default:
		t.Fatalf("Invalid request type: %d", reqType)
	}

	_, err := client.K8sClient.ExecInPod(ctx, client.Pod.Namespace, client.Pod.Name, "", cmd)
	if expectFail && err == nil {
		t.Failf("Expected curl to fail, but it succeeded")
	} else if !expectFail && err != nil {
		t.Fatalf("Failed to curl server: %s", err)
	}

	// Wait until tcpdump has exited
	killCmd()
	<-bgExited

	// Assert no traffic leak
	checkPcapForLeak(ctx, t, clientHost)
}

// bytes.Buffer from the stdlib is non-thread safe, thus our custom
// implementation. Unfortunately, we cannot use io.Pipe, as Write() blocks until
// Read() has read all content, which makes it deadlock-prone when used with
// ExecInPodWithWriters() running in a separate goroutine.
type safeBuffer struct {
	sync.Mutex
	b bytes.Buffer
}

func (b *safeBuffer) Read(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	return b.b.Read(p)
}

func (b *safeBuffer) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	return b.b.Write(p)
}

func (b *safeBuffer) String() string {
	b.Lock()
	defer b.Unlock()
	return b.b.String()
}

func (b *safeBuffer) ReadString(d byte) (string, error) {
	b.Lock()
	defer b.Unlock()
	return b.b.ReadString(d)
}

func NodeToNodeEncryption() check.Scenario {
	return &nodeToNodeEncryption{}
}

type nodeToNodeEncryption struct{}

func (s *nodeToNodeEncryption) Name() string {
	return "node-to-node-encryption"
}

func (s *nodeToNodeEncryption) Run(ctx context.Context, t *check.Test) {
	client := t.Context().RandomClientPod()

	var server check.Pod
	for _, pod := range t.Context().EchoPods() {
		// Make sure that the server pod is on another node than client
		if pod.Pod.Status.HostIP != client.Pod.Status.HostIP {
			server = pod
			break
		}
	}

	// clientHost is a pod running on the same node as the client pod, just in
	// the host netns.
	clientHost := t.Context().HostNetNSPodsByNode()[client.Pod.Spec.NodeName]
	// serverHost is a pod running in a remote node's host netns.
	serverHost := t.Context().HostNetNSPodsByNode()[server.Pod.Spec.NodeName]

	t.ForEachIPFamily(func(ipFam check.IPFamily) {
		// Test pod-to-remote-host (ICMP Echo instead of HTTP because a remote host
		// does not have a HTTP server running)
		testNoTrafficLeak(ctx, t, client, &serverHost, &clientHost, requestICMPEcho, ipFam, false)
		// Test host-to-remote-host
		testNoTrafficLeak(ctx, t, &clientHost, &serverHost, &clientHost, requestICMPEcho, ipFam, false)
		// Test host-to-remote-pod
		testNoTrafficLeak(ctx, t, &clientHost, &server, &clientHost, requestHTTP, ipFam, false)
	})
}
