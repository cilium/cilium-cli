package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/utils/features"
	v1 "k8s.io/api/core/v1"
)

func BPFMasquerade() check.Scenario {
	return &bpfMasquerade{}
}

type bpfMasquerade struct{}

func (s *bpfMasquerade) Name() string {
	return "bpf-masquerade"
}

func (s *bpfMasquerade) Run(ctx context.Context, t *check.Test) {
	var i int
	ct := t.Context()

	for _, client := range ct.ClientPods() {
		client := client // copy to avoid memory aliasing when using reference

		for _, echo := range ct.ExternalEchoPods() {
			echo := echo // copy to avoid memory aliasing when using reference
			if !echo.Pod.Spec.HostNetwork {
				continue
			}

			baseURL := fmt.Sprintf("%s://%s:%d/echo", echo.Scheme(), echo.Pod.Status.HostIP, 8080)
			ep := check.HTTPEndpoint(echo.Name(), baseURL)
			t.NewAction(s, fmt.Sprintf("curl-%d", i), &client, ep, features.IPFamilyAny).Run(func(a *check.Action) {
				a.ExecInPod(ctx, ct.CurlCommandWithOutput(ep, features.IPFamilyAny))
				out := a.CmdOutput()
				m := map[string]any{}
				if err := json.Unmarshal([]byte(out), &m); err != nil {
					a.Fail(err)
					return
				}

				remote, ok := m["remote_ip"]
				if !ok {
					a.Failf("echo response did not contain remote addr: %s", out)
				}

				var nodeaddr string
				node := ct.Nodes()[client.NodeName()]
				for _, addr := range node.Status.Addresses {
					if addr.Type == v1.NodeInternalIP {
						ip := net.ParseIP(addr.Address)
						if ip == nil {
							continue
						}
						// todo: ipv6...
						if ip.To4() != nil {
							nodeaddr = addr.Address
							break
						}
					}
				}

				remoteip := strings.Split(remote.(string), ":")[0]

				if remoteip != nodeaddr {
					a.Failf("request traffic should have been masqueraded: remote=%s nodeaddr=%s", remote, nodeaddr)
				}
				/*a.ValidateFlows(ctx, client, a.GetEgressRequirements(check.FlowParameters{
					// Because the HostPort request is NATed, we might only
					// observe flows after DNAT has been applied (e.g. by
					// HostReachableServices),
					AltDstIP:   echo.Address(features.IPFamilyAny),
					AltDstPort: echo.Port(),
				}))*/
			})

			//i++
		}
	}
}
