package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cilium/cilium-cli/connectivity/check"
)

// TCP Network Performance
func TCPPodtoPod(n string, ct *check.ConnectivityTest) check.Scenario {
	return &tcpPodToPod{
		name:             n,
		connectivityTest: ct,
	}
}

type tcpPodToPod struct {
	name             string
	connectivityTest *check.ConnectivityTest
}

func (s *tcpPodToPod) Name() string {
	tn := "perf[tcp]-pod-to-pod"
	if s.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, s.name)
}

func (s *tcpPodToPod) Run(ctx context.Context, t *check.Test) {
	for _, c := range t.Context().PerfClientPods() {
		c := c
		for _, server := range t.Context().PerfServerPod() {
			t.NewAction(s, "iperf-tcp", &c, server).Run(func(a *check.Action) {
				iperfStream(ctx, server.Pod.Status.PodIP, c.Pod.Name, a, s.connectivityTest, 1, false)
			})
		}
	}
}

// UDP Network pPerformance
func UDPPodtoPod(n string, ct *check.ConnectivityTest) check.Scenario {
	return &udpPodtoPod{
		name:             n,
		connectivityTest: ct,
	}
}

type udpPodtoPod struct {
	name             string
	connectivityTest *check.ConnectivityTest
}

func (s *udpPodtoPod) Name() string {
	tn := "perf[udp]-pod-to-pod"
	if s.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, s.name)
}

func (s *udpPodtoPod) Run(ctx context.Context, t *check.Test) {
	for _, c := range t.Context().PerfClientPods() {
		c := c
		for _, server := range t.Context().PerfServerPod() {
			t.NewAction(s, "iperf-udp", &c, server).Run(func(a *check.Action) {
				iperfStream(ctx, server.Pod.Status.PodIP, c.Pod.Name, a, s.connectivityTest, 1, true)
			})
		}
	}
}

func iperfStream(ctx context.Context, sip string, podname string, a *check.Action, ct *check.ConnectivityTest, samples int, udp bool) {
	// Allow the user to override the number of samples to capture
	env, _ := strconv.Atoi(os.Getenv("samples"))
	if samples < env {
		samples = env
	}
	for i := 0; i < samples; i++ {
		r := make(map[string]string)
		r["test"] = "stream"
		exec := []string{"/usr/bin/iperf3", "-c", sip, "-t", "60", "-Z", "--no-delay", "-J"}
		if udp {
			exec = []string{"/usr/bin/iperf3", "-u", "-c", sip, "-t", "5", "-Z", "--no-delay", "-J"}
		}
		a.ExecInPod(ctx, exec)
		var payload map[string]interface{}
		err := json.Unmarshal([]byte(a.CmdOutput()), &payload)
		if err != nil {
			a.Fatal("unable to parse output from iperf3")
		}
		r["iteration"] = fmt.Sprintf("%d", i)
		if !udp {
			r["protocol"] = "tcp"
			r["bps"] = fmt.Sprintf("%f", (payload["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["bits_per_second"]).(float64)/1000000.0)
			r["rt"] = fmt.Sprintf("%f", payload["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["retransmits"])
			ct.PerfResults[fmt.Sprintf("%s-tcp-stream-%d", podname, i)] = r
		} else {
			r["protocol"] = "udp"
			r["bps"] = fmt.Sprintf("%f", (payload["end"].(map[string]interface{})["sum"].(map[string]interface{})["bits_per_second"]).(float64)/1000000)
			ct.PerfResults[fmt.Sprintf("%s-udp-stream-%d", podname, i)] = r
		}

	}
}
