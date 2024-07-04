// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package builder

import (
	_ "embed"

	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/connectivity/deploy"
	"github.com/cilium/cilium-cli/connectivity/tests"
	"github.com/cilium/cilium-cli/utils/features"
)

var (
	//go:embed manifests/node-local-dns-lrp.yaml
	nodeDNSLocalRedirectPolicyYAML string

	//go:embed manifests/client-egress-node-local-dns.yaml
	clientEgressNodeLocalDNSYAML string
)

type localRedirectPolicyWithNodeDNS struct{}

func (t localRedirectPolicyWithNodeDNS) build(ct *check.ConnectivityTest, template map[string]string) {
	newTest("local-redirect-policy-with-node-dns", ct).
		WithSetupFunc(deploy.NodeDNS).
		WithCiliumPolicy(template["clientEgressNodeLocalDNSYAML"]).
		WithCiliumLocalRedirectPolicy(check.CiliumLocalRedirectPolicyParams{
			Policy:                  nodeDNSLocalRedirectPolicyYAML,
			NameSpace:               "kube-system",
			Name:                    "nodelocaldns",
			SkipRedirectFromBackend: false,
		}).
		WithFeatureRequirements(features.RequireEnabled(features.LocalRedirectPolicy)).
		WithFeatureRequirements(features.RequireEnabled(features.KPRSocketLB)).
		WithScenarios(
			tests.LRPWithNodeDNS(),
		)
}
