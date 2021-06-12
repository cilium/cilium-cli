// Copyright 2020-2021 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/cilium/cilium-cli/connectivity/check"
)

// RestartDNS sends multiple HTTP requests from a randomly selected client Pod
// to all Services in the test context. Cilium agent restart on the client's node
// is executed while the test is ongoing.
func RestartDNS(name string) check.Scenario {
	return &restartDNS{
		name: name,
	}
}

// restartDNS implements a Scenario.
type restartDNS struct {
	name string
}

func (s *restartDNS) Name() string {
	tn := "restart-dns"
	if s.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, s.name)
}

func (s *restartDNS) Run(ctx context.Context, t *check.Test) {
	for _, src := range t.Context().ClientPods() {
		for _, dst := range t.Context().ClientPods() {
			if src.Pod.Status.PodIP == dst.Pod.Status.PodIP {
				// Currently we only get flows once per IP,
				// skip pings to self.
				continue
			}
			ping := func(i int, maxFail int) (success bool) {
				t.NewAction(s, fmt.Sprintf("ping-%d", i), &src, &dst).Run(func(a *check.Action) {
					if i < maxFail {
						a.AllowFail()
					}
					a.ExecInPod(ctx, pingName(dst))
					success = a.Succeeded()
					egressFlowRequirements := a.GetEgressRequirements(check.FlowParameters{
						Protocol:    check.ICMP,
						DNSRequired: true,
					})
					// Do not validate allowed failures
					if !success && i < maxFail {
						return
					}
					a.ValidateFlows(ctx, src.Name(), src.Address(), egressFlowRequirements)
				})
				return success
			}

			// Run 20 queries before restarting
			for i := -20; i < 0; i++ {
				ping(i, 0)
				time.Sleep(125 * time.Millisecond)
			}

			// Restart Cilium PODs and keep pinging
			t.RestartCiliumPods(src)

			end := time.Now().Add(40 * time.Second)
			maxFail := 10
			dropped := false
			failed := false
			streak := 0
			// Run loop until Cilium Agents are available again
			for i := 0; i < 200; i++ {
				success := ping(i, maxFail)
				if success {
					streak++
					// Do not allow new failures after fails have turned to success
					if dropped {
						maxFail = 0
						// Stop the test after one minute if ping has recovered and Cilium API is up
						if time.Now().After(end) && streak > 10 && t.CiliumRunning() {
							break
						}
					}
					time.Sleep(250 * time.Millisecond)
				} else {
					dropped = true
					streak = 0
					if maxFail == 0 {
						failed = true
					}
				}
			}

			if dropped && !failed {
				t.ClearFail()
			}

			// Run only once
			return
		}
	}
}
