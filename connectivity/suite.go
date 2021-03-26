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

package connectivity

import (
	"context"

	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/connectivity/tests"
)

var (
	//l3Policies = check.PolicyContext{}
)

func Run(ctx context.Context, k *check.K8sConnectivityCheck) error {
	return k.Run(ctx,
		// This should fail, failures not yet handled properly:
		(&tests.PodToPod{Variant: "-client-egress-only-dns-SHOULD-FAIL"}).WithPolicy(clientEgressOnlyDNSPolicyYaml),
		// Policy installed with 'WithPolicy()' is automatically removed, so this should succeed:
		&tests.PodToPod{},
		// This policy allows port 8080 from client to echo, so this should succeed
		(&tests.PodToPod{Variant: "-client-egress-to-echo"}).WithPolicy(clientEgressToEchoPolicyYaml),
		// This should also succeed, no policy applied
		&tests.PodToPod{Variant: "-2"},
		// Apply policy that is kept around until explicitly removed
		&check.CiliumNetworkPolicy{Title: "client-egress-only-dns", Policy: clientEgressOnlyDNSPolicyYaml},
		&tests.PodToPod{Variant: "-client-egress-only-dns-SHOULD-FAIL-too"},
		&check.CiliumNetworkPolicy{}, // delete all applied policies
		// This should also succeed, no policy applied
		&tests.PodToPod{Variant: "-3"},
		&tests.PodToService{},
		&tests.PodToNodePort{},
		&tests.PodToLocalNodePort{},
		&tests.PodToWorld{},
		&tests.PodToHost{},
	)
}

var (
	clientEgressOnlyDNSPolicyYaml = `
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  namespace: cilium-test
  name: client-egress-only-dns
spec:
  endpointSelector:
    matchLabels:
      kind: client
  egress:
  - toPorts:
    - ports:
      - port: "53"
        protocol: ANY
    toEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: kube-system
        k8s:k8s-app: kube-dns
`
	clientEgressToEchoPolicyYaml = `apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  namespace: cilium-test
  name: client-egress-to-echo
spec:
  endpointSelector:
    matchLabels:
      kind: client
  egress:
  - toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
    toEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: cilium-test
        k8s:kind: echo
  - toPorts:
    - ports:
      - port: "53"
        protocol: ANY
    toEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: kube-system
        k8s:k8s-app: kube-dns
`
)
