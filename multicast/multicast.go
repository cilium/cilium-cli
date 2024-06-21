// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package multicast

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/k8s"
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/node/addressing"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Multicast struct {
	client *k8s.Client
	params Parameters
}

type Parameters struct {
	CiliumNamespace  string
	Writer           io.Writer
	WaitDuration     time.Duration
	MulticastGroupIP string
	All              bool
}

func NewMulticast(client *k8s.Client, p Parameters) *Multicast {
	return &Multicast{
		client: client,
		params: p,
	}
}

func (m *Multicast) getCiliumNode(ctx context.Context, nodeName string) (v2.CiliumNode, error) {
	ciliumNodes, err := m.client.ListCiliumNodes(ctx)
	if err != nil {
		return v2.CiliumNode{}, err
	}
	var ciliumNode v2.CiliumNode
	for _, node := range ciliumNodes.Items {
		if node.Name == nodeName {
			ciliumNode = node
		}
	}
	return ciliumNode, nil
}

func (m *Multicast) getCiliumInternalIP(nodeName string) (v2.NodeAddress, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.params.WaitDuration)
	defer cancel()
	ciliumNode, err := m.getCiliumNode(ctx, nodeName)
	if err != nil {
		return v2.NodeAddress{}, fmt.Errorf("unable to get cilium node: %w", err)
	}
	addrs := ciliumNode.Spec.Addresses
	var ciliumInternalIP v2.NodeAddress
	for _, addr := range addrs {
		if addr.AddrType() == addressing.NodeCiliumInternalIP {
			ip, err := netip.ParseAddr(addr.IP)
			if err != nil {
				continue
			}
			if ip.Is4() {
				ciliumInternalIP = addr
			}
		}
	}
	if ciliumInternalIP.IP == "" {
		return v2.NodeAddress{}, fmt.Errorf("ciliumInternalIP not found")
	}
	return ciliumInternalIP, nil
}

// ListGroup lists multicast groups in every node
func (m *Multicast) ListGroups() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.params.WaitDuration)
	defer cancel()

	ciliumPodsList, err := m.client.ListPods(ctx, m.params.CiliumNamespace, metav1.ListOptions{LabelSelector: "k8s-app=cilium"})
	if err != nil {
		return err
	}
	ciliumPods := ciliumPodsList.Items

	var wg sync.WaitGroup
	errCh := make(chan error, len(ciliumPods))
	wg.Add(len(ciliumPods))

	for _, ciliumPod := range ciliumPods {
		go func(pod corev1.Pod) {
			defer wg.Done()
			// List multicast groups
			cmd := []string{"cilium-dbg", "bpf", "multicast", "group", "list"}
			output, err := m.client.ExecInPod(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
			if err != nil {
				errCh <- err
				return
			}
			outputString := "Node: " + pod.Spec.NodeName + "\n" + output.String()
			fmt.Fprintln(m.params.Writer, outputString)
		}(ciliumPod)
	}

	wg.Wait()
	close(errCh)

	var errRet error
	for fetchdata := range errCh {
		if fetchdata != nil {
			errRet = errors.Join(errRet, fetchdata)
		}
	}
	return errRet
}

// ListSubscriber lists multicast subscribers in every node for the specified multicast group or all multicast groups
func (m *Multicast) ListSubscribers() error {
	if m.params.MulticastGroupIP == "" && !m.params.All {
		return fmt.Errorf("group-ip or all flag must be specified")
	} else if m.params.MulticastGroupIP != "" && m.params.All {
		return fmt.Errorf("only one of group-ip or all flag must be specified")
	}

	var target string
	if m.params.All {
		target = "all"
	} else {
		target = m.params.MulticastGroupIP
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.params.WaitDuration)
	defer cancel()

	ciliumPodsList, err := m.client.ListPods(ctx, m.params.CiliumNamespace, metav1.ListOptions{LabelSelector: "k8s-app=cilium"})
	if err != nil {
		return err
	}
	ciliumPods := ciliumPodsList.Items

	var wg sync.WaitGroup
	errCh := make(chan error, len(ciliumPods))
	wg.Add(len(ciliumPods))

	for _, ciliumPod := range ciliumPods {
		go func(pod corev1.Pod) {
			defer wg.Done()
			// List multicast subscribers
			cmd := []string{"cilium-dbg", "bpf", "multicast", "subscriber", "list", target}
			output, stdErr, err := m.client.ExecInPodWithStderr(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
			if err != nil {
				if strings.Contains(stdErr.String(), "does not exist") {
					fmt.Fprintf(m.params.Writer, "Multicast group %s does not exist in %s\n", target, pod.Spec.NodeName)
					return
				}
				errCh <- err
				return
			}
			outputString := "Node: " + pod.Spec.NodeName + "\n" + output.String()
			fmt.Fprintln(m.params.Writer, outputString)
		}(ciliumPod)
	}

	wg.Wait()
	close(errCh)

	var errRet error
	for fetchdata := range errCh {
		if fetchdata != nil {
			errRet = errors.Join(errRet, fetchdata)
		}
	}
	return errRet
}

func (m *Multicast) populateMaps(ciliumPods []corev1.Pod, ipToPodMap map[v2.NodeAddress]string, ipToNodeMap map[v2.NodeAddress]string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(ciliumPods))
	wg.Add(len(ciliumPods))

	for _, ciliumPod := range ciliumPods {
		go func(pod corev1.Pod) {
			defer wg.Done()
			ciliumInternalIP, err := m.getCiliumInternalIP(pod.Spec.NodeName)
			if err != nil {
				errCh <- err
				return
			}
			ipToPodMap[ciliumInternalIP] = pod.Name
			ipToNodeMap[ciliumInternalIP] = pod.Spec.NodeName
		}(ciliumPod)
	}

	wg.Wait()
	close(errCh)

	var errRet error
	for fetchdata := range errCh {
		if fetchdata != nil {
			errRet = errors.Join(errRet, fetchdata)
		}
	}
	return errRet
}

// AddAllNodes add CiliumInternalIPs of all nodes to the specified multicast group as subscribers in every cilium-agent
func (m *Multicast) AddAllNodes() error {
	if m.params.MulticastGroupIP == "" {
		return fmt.Errorf("group-ip must be specified")
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.params.WaitDuration)
	defer cancel()

	ciliumPodsList, err := m.client.ListPods(ctx, m.params.CiliumNamespace, metav1.ListOptions{LabelSelector: "k8s-app=cilium"})
	if err != nil {
		return err
	}
	ciliumPods := ciliumPodsList.Items

	//Create a map of ciliumInternalIPs of all nodes
	ipToPodMap := make(map[v2.NodeAddress]string)
	ipToNodeMap := make(map[v2.NodeAddress]string)

	if err := m.populateMaps(ciliumPods, ipToPodMap, ipToNodeMap); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(ciliumPods))
	wg.Add(len(ciliumPods))

	for _, ciliumPod := range ciliumPods {
		go func(pod corev1.Pod) {
			defer wg.Done()
			//If there are not specified multicast group, create it
			cmd := []string{"cilium-dbg", "bpf", "multicast", "subscriber", "list", m.params.MulticastGroupIP}
			_, stdErr, err := m.client.ExecInPodWithStderr(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
			if err != nil {
				if !strings.Contains(stdErr.String(), "does not exist") {
					errCh <- err
					return
				}
				cmd = []string{"cilium-dbg", "bpf", "multicast", "group", "add", m.params.MulticastGroupIP}
				_, err := m.client.ExecInPod(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
				if err != nil {
					errCh <- err
					fmt.Fprintf(m.params.Writer, "Unable to create multicast group %s in %s\n", m.params.MulticastGroupIP, pod.Name)
					return
				}
			}

			//Add all ciliumInternalIPs of all nodes to the multicast group as subscribers
			cnt := 0
			var nodeLists []string
			var displayOutput string
			for ip, podName := range ipToPodMap {
				if ip.IP != "" && pod.Name != podName { //My node itself does not need to be in a multicast group.
					cmd = []string{"cilium-dbg", "bpf", "multicast", "subscriber", "add", m.params.MulticastGroupIP, ip.IP}
					_, stdErr, err := m.client.ExecInPodWithStderr(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
					if err == nil {
						cnt++
						nodeLists = append(nodeLists, ipToNodeMap[ip])
					} else if !strings.Contains(stdErr.String(), "already exists") {
						errCh <- err
						fmt.Fprintf(m.params.Writer, "Unable to add node %s to multicast group %s in %s by fatal error\n", ip.IP, m.params.MulticastGroupIP, pod.Spec.NodeName)
						return
					}
				}
			}
			if cnt == 0 {
				fmt.Fprintf(m.params.Writer, "Unable to add any node to multicast group %s in %s\n", m.params.MulticastGroupIP, pod.Spec.NodeName)
				return
			}
			if cnt == 1 {
				displayOutput = "Added a node ("
			} else {
				displayOutput = fmt.Sprintf("Added %d nodes (", cnt)
			}
			for i, node := range nodeLists {
				if i == len(nodeLists)-1 {
					displayOutput += node
				} else {
					displayOutput += node + ", "
				}
			}
			displayOutput += fmt.Sprintf(") to multicast group %s in %s\n", m.params.MulticastGroupIP, pod.Spec.NodeName)
			fmt.Fprint(m.params.Writer, displayOutput)
		}(ciliumPod)
	}

	wg.Wait()
	close(errCh)
	var errRet error
	for fetchdata := range errCh {
		if fetchdata != nil {
			errRet = errors.Join(errRet, fetchdata)
		}
	}
	return errRet
}

// DelAllNodes delete CiliumInternalIPs of all nodes from the specified multicast group's subscribers in every cilium-agent
func (m *Multicast) DelAllNodes() error {
	if m.params.MulticastGroupIP == "" {
		return fmt.Errorf("group-ip must be specified")
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.params.WaitDuration)
	defer cancel()

	ciliumPodsList, err := m.client.ListPods(ctx, m.params.CiliumNamespace, metav1.ListOptions{LabelSelector: "k8s-app=cilium"})
	if err != nil {
		return err
	}
	ciliumPods := ciliumPodsList.Items

	//Create a map of ciliumInternalIPs of all nodes
	ipToPodMap := make(map[v2.NodeAddress]string)
	ipToNodeMap := make(map[v2.NodeAddress]string)
	for _, ciliumPod := range ciliumPods {
		ciliumInternalIP, err := m.getCiliumInternalIP(ciliumPod.Spec.NodeName)
		if err != nil {
			return err
		}
		ipToPodMap[ciliumInternalIP] = ciliumPod.Name
		ipToNodeMap[ciliumInternalIP] = ciliumPod.Spec.NodeName
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(ciliumPods))
	wg.Add(len(ciliumPods))

	for _, ciliumPod := range ciliumPods {
		go func(pod corev1.Pod) {
			defer wg.Done()
			//Delete all ciliumInternalIPs of all nodes from the multicast group's subscribers
			cmd := []string{"cilium-dbg", "bpf", "multicast", "group", "delete", m.params.MulticastGroupIP}
			_, stdErr, err := m.client.ExecInPodWithStderr(ctx, pod.Namespace, pod.Name, defaults.AgentContainerName, cmd)
			if err != nil {
				if !strings.Contains(stdErr.String(), "does not exist") {
					errCh <- err
					fmt.Fprintf(m.params.Writer, "Unable to delete multicast group %s in %s by fatal error\n", m.params.MulticastGroupIP, pod.Spec.NodeName)
					return
				}
				fmt.Fprintf(m.params.Writer, "Multicast group %s does not exist in %s\n", m.params.MulticastGroupIP, pod.Spec.NodeName)
				return
			}
			fmt.Fprintf(m.params.Writer, "Deleted multicast group %s in %s\n", m.params.MulticastGroupIP, pod.Spec.NodeName)
		}(ciliumPod)
	}

	wg.Wait()
	close(errCh)

	var errRet error
	for fetchdata := range errCh {
		if fetchdata != nil {
			errRet = errors.Join(errRet, fetchdata)
		}
	}
	return errRet
}
