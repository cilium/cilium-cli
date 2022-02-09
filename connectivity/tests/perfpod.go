package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cilium/cilium-cli/connectivity/check"
)

// TCP Network perforamnce between pods
func TCPPodtoPod(n string) check.Scenario {
	return &tcpPodToPod{
		name: n,
	}
}

type tcpPodToPod struct {
	name string
}

func (s *tcpPodToPod) Name() string {
	tn := "perf[tcp]-pod-to-pod"
	if s.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, s.name)
}

type Result struct {
	iteration int
	bps       interface{}
	rt        interface{}
}

func iperf(sip string, ctx context.Context, podname string, a *check.Action, samples int, udp bool) map[string]Result {
	results := make(map[string]Result)
	// Allow the user to override the number of samples to capture
	env, _ := strconv.Atoi(os.Getenv("samples"))
	if samples < env {
		samples = env
	}
	for i := 0; i < samples; i++ {
		var r Result
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
		r.iteration = i
		if !udp {
			r.bps = payload["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["bits_per_second"]
			r.rt = payload["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["retransmits"]
			results[fmt.Sprintf("%s-tcp-stream-%d", podname, i)] = r
		} else {
			r.bps = payload["end"].(map[string]interface{})["sum"].(map[string]interface{})["bits_per_second"]
			results[fmt.Sprintf("%s-udp-stream-%d", podname, i)] = r
		}

	}
	return results
}

func (s *tcpPodToPod) Run(ctx context.Context, t *check.Test) {
	for _, c := range t.Context().PerfClientPods() {
		for _, server := range t.Context().PerfServerPod() {
			t.NewAction(s, "iperf-tcp", &c, server).Run(func(a *check.Action) {
				results := iperf(server.Pod.Status.PodIP, ctx, c.Pod.Name, a, 1, false)
				for test, result := range results {

					fmt.Printf("\n📄 Results for : %s\n", test)
					fmt.Printf("✅ TCP Retransmissions (lower is better) : %d\n", int(result.rt.(float64)))
					fmt.Printf("✅ TCP Stream Performance - Iteration %d : %f(Mbps)\n", result.iteration, (result.bps).(float64)/1000000.0)
				}
			})
		}
	}
}

// UDP Network performance between
// TCP Network perforamnce between pods
func UDPPodtoPod(n string) check.Scenario {
	return &udpPodtoPod{
		name: n,
	}
}

type udpPodtoPod struct {
	name string
}

func (s *udpPodtoPod) Name() string {
	tn := "perf-pod-to-pod"
	if s.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, s.name)
}

func (s *udpPodtoPod) Run(ctx context.Context, t *check.Test) {
	for _, c := range t.Context().PerfClientPods() {
		for _, server := range t.Context().PerfServerPod() {
			t.NewAction(s, "iperf-udp", &c, server).Run(func(a *check.Action) {
				results := iperf(server.Pod.Status.PodIP, ctx, c.Pod.Name, a, 1, true)
				for test, result := range results {
					fmt.Printf("\n📄 Results for : %s\n", test)
					fmt.Printf("✅ UDP Stream Performance - Iteration %d : %f(Mbps)\n", result.iteration, (result.bps).(float64)/1000000.0)
				}
			})
		}
	}
	t.Info("Message")
}
