// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package builder

import (
	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/connectivity/tests"
	"github.com/cilium/cilium-cli/utils/features"
)

type hostFirewallEgressToFqdns struct{}

func (t hostFirewallEgressToFqdns) build(ct *check.ConnectivityTest, templates map[string]string) {
	// This policy only allows port 80 to domain-name, default one.one.one.one., DNS proxy enabled.
	newTest("host-firewall-egress-to-fqdns", ct).
		WithCondition(func() bool { return ct.Params().IncludeUnsafeTests }).
		WithFeatureRequirements(
			features.RequireEnabled(features.L7Proxy),
			features.RequireEnabled(features.HostFirewall)).
		WithCiliumClusterwidePolicy(templates["hostFirewallEgressToFQDNsPolicyYAML"]).
		WithScenarios(
			tests.HostToWorld(),
			tests.HostToWorld2()).
		WithExpectations(func(a *check.Action) (egress, ingress check.Result) {
			extTarget := ct.Params().ExternalTarget
			if a.Destination().Address(features.GetIPFamily(extTarget)) == extTarget {
				return check.ResultDNSOK, check.ResultNone
			}

			return check.ResultDNSOKDropCurlTimeout, check.ResultNone
		})
}
