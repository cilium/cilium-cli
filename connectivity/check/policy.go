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

package check

import (
	"context"
	"fmt"
	"strings"

	"github.com/cilium/cilium/pkg/k8s/client/clientset/versioned/scheme"

	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CiliumNetworkPolicy implements ConnectivityTest interface so that
// policies can be applied between tests and policy apply failures can
// be reported like any other test results.
type CiliumNetworkPolicy struct {
	Title  string
	Policy string
}

// Name returns the absolute name of the policy
func (p *CiliumNetworkPolicy) Name() string {
	return p.Title
}

// Run applies the policy, use nil policy to delete all policies
func (p *CiliumNetworkPolicy) Run(ctx context.Context, c TestContext) {
	failures := c.ApplyPolicyYaml(ctx, p.Policy)
	c.Report(TestResult{
		Name:     p.Name(),
		Failures: failures,
		Warnings: 0,
	})
}

type WithPolicy struct {
	err  error
	cnps []*ciliumv2.CiliumNetworkPolicy
}

func (wp *WithPolicy) Parse(policy string) {
	wp.cnps, wp.err = ParsePolicyYaml(policy)
}

// Apply returns the number of failures and a cancel function that deletes the applied policies
func (wp *WithPolicy) Apply(ctx context.Context, c TestContext) (failures int, cancel func()) {
	if wp.err != nil {
		return 1, func() {}
	}

	for _, cnp := range wp.cnps {
		failures += c.ApplyCNP(ctx, cnp)
	}
	return failures, func() {
		for _, cnp := range wp.cnps {
			c.DeleteCNP(ctx, cnp)
		}
	}
}

// Run applies the policy, use empty policy to delete all policies
func (k *K8sConnectivityCheck) ApplyPolicyYaml(ctx context.Context, policy string) int {
	if policy == "" {
		// Delete all policies
		return k.ApplyCNP(ctx, nil)
	}
	failures := 0
	cnps, err := ParsePolicyYaml(policy)
	if err != nil {
		k.Log("❌ %s", err)
		failures++
	} else {
		for _, cnp := range cnps {
			failures += k.ApplyCNP(ctx, cnp)
		}
	}
	return failures
}

// ParsePolicyYaml decodes policy yaml into a slice of CiliumNetworkPolicies
func ParsePolicyYaml(policy string) (cnps []*ciliumv2.CiliumNetworkPolicy, err error) {
	if policy == "" {
		return nil, nil
	}
	yamls := strings.Split(policy, "---")
	for _, yaml := range yamls {
		if yaml == "\n" || yaml == "" {
			continue
		}
		obj, groupVersionKind, err := scheme.Codecs.UniversalDeserializer().Decode([]byte(yaml), nil, nil)
		if err != nil {
			return nil, fmt.Errorf("Resource decode error (%s) in: %s", err, yaml)
		}
		switch groupVersionKind.Kind {
		case "CiliumNetworkPolicy":
			cnp, ok := obj.(*ciliumv2.CiliumNetworkPolicy)
			if !ok {
				return nil, fmt.Errorf("Object cast to CiliumNetworkPolicy failed: %s", yaml)
			}
			cnps = append(cnps, cnp)
		default:
			return nil, fmt.Errorf("Unknown policy type '%s' in: %s", groupVersionKind.Kind, yaml)
		}
	}
	return cnps, nil
}

// DeleteCNP deletes a CNP
func (k *K8sConnectivityCheck) DeleteCNP(ctx context.Context, cnp *ciliumv2.CiliumNetworkPolicy) {
	name := cnp.Namespace + "/" + cnp.Name
	if err := k.deleteCNP(ctx, cnp); err != nil {
		k.Log("❌ [%s] policy delete failed: %s", name, err)
	}
	delete(k.policies, name)
}

// ApplyCNP returns the number of failures
func (k *K8sConnectivityCheck) ApplyCNP(ctx context.Context, cnp *ciliumv2.CiliumNetworkPolicy) int {
	failures := 0
	if cnp == nil {
		k.Header("🔌 Deleting all previously applied policies...")
		for _, cnp := range k.policies {
			k.DeleteCNP(ctx, cnp)
		}
	} else {
		name := cnp.Namespace + "/" + cnp.Name
		k.Header("🔌 [%s] Applying CiliumNetworkPolicy...", name)
		k8sCNP, err := k.updateOrCreateCNP(ctx, cnp)
		if err == nil {
			k.Log("✅ [%s] CiliumNetworkPolicy applied", name)
			k.policies[name] = k8sCNP
		} else {
			k.Log("❌ policy apply failed: %s", err)
			failures++
		}
	}
	return failures
}

func (k *K8sConnectivityCheck) updateOrCreateCNP(ctx context.Context, cnp *ciliumv2.CiliumNetworkPolicy) (*ciliumv2.CiliumNetworkPolicy, error) {
        k8sCNP, err := k.clients.src.GetCiliumNetworkPolicy(ctx, cnp.Namespace, cnp.Name, metav1.GetOptions{})
        if err == nil {
                k8sCNP.ObjectMeta.Labels = cnp.ObjectMeta.Labels
                k8sCNP.Spec = cnp.Spec
                k8sCNP.Specs = cnp.Specs
                k8sCNP.Status = ciliumv2.CiliumNetworkPolicyStatus{}
                return k.clients.src.UpdateCiliumNetworkPolicy(ctx, k8sCNP, metav1.UpdateOptions{})
        }
        return k.clients.src.CreateCiliumNetworkPolicy(ctx, cnp, metav1.CreateOptions{})
}

func (k *K8sConnectivityCheck) deleteCNP(ctx context.Context, cnp *ciliumv2.CiliumNetworkPolicy) error {
        return k.clients.src.DeleteCiliumNetworkPolicy(ctx, cnp.Namespace, cnp.Name, metav1.DeleteOptions{})
}

func (k *K8sConnectivityCheck) updateOrCreateCCNP(ctx context.Context, ccnp *ciliumv2.CiliumClusterwideNetworkPolicy) (*ciliumv2.CiliumClusterwideNetworkPolicy, error) {
        k8sCCNP, err := k.clients.src.GetCiliumClusterwideNetworkPolicy(ctx, ccnp.Name, metav1.GetOptions{})
        if err == nil {
                k8sCCNP.ObjectMeta.Labels = ccnp.ObjectMeta.Labels
                k8sCCNP.Spec = ccnp.Spec
                k8sCCNP.Specs = ccnp.Specs
                k8sCCNP.Status = ciliumv2.CiliumNetworkPolicyStatus{}
                return k.clients.src.UpdateCiliumClusterwideNetworkPolicy(ctx, k8sCCNP, metav1.UpdateOptions{})
        }
        return k.clients.src.CreateCiliumClusterwideNetworkPolicy(ctx, ccnp, metav1.CreateOptions{})
}

func (k *K8sConnectivityCheck) deleteCCNP(ctx context.Context, ccnp *ciliumv2.CiliumNetworkPolicy) error {
        return k.clients.src.DeleteCiliumClusterwideNetworkPolicy(ctx, ccnp.Name, metav1.DeleteOptions{})
}
