// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package utils

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"

	k8sConst "github.com/cilium/cilium/pkg/k8s/apis/cilium.io"
	slim_metav1 "github.com/cilium/cilium/pkg/k8s/slim/k8s/apis/meta/v1"
	"github.com/cilium/cilium/pkg/labels"
	"github.com/cilium/cilium/pkg/logging"
	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/policy/api"
)

const (
	// subsysK8s is the value for logfields.LogSubsys
	subsysK8s = "k8s"
	// podPrefixLbl is the value the prefix used in the label selector to
	// represent pods on the default namespace.
	podPrefixLbl = labels.LabelSourceK8sKeyPrefix + k8sConst.PodNamespaceLabel

	// podAnyPrefixLbl is the value of the prefix used in the label selector to
	// represent pods in the default namespace for any source type.
	podAnyPrefixLbl = labels.LabelSourceAnyKeyPrefix + k8sConst.PodNamespaceLabel

	// podK8SNamespaceLabelsPrefix is the prefix use in the label selector for namespace labels.
	podK8SNamespaceLabelsPrefix = labels.LabelSourceK8sKeyPrefix + k8sConst.PodNamespaceMetaLabelsPrefix
	// podAnyNamespaceLabelsPrefix is the prefix use in the label selector for namespace labels
	// for any source type.
	podAnyNamespaceLabelsPrefix = labels.LabelSourceAnyKeyPrefix + k8sConst.PodNamespaceMetaLabelsPrefix

	// podInitLbl is the label used in a label selector to match on
	// initializing pods.
	podInitLbl = labels.LabelSourceReservedKeyPrefix + labels.IDNameInit

	// ResourceTypeCiliumNetworkPolicy is the resource type used for the
	// PolicyLabelDerivedFrom label
	ResourceTypeCiliumNetworkPolicy = "CiliumNetworkPolicy"

	// ResourceTypeCiliumClusterwideNetworkPolicy is the resource type used for the
	// PolicyLabelDerivedFrom label
	ResourceTypeCiliumClusterwideNetworkPolicy = "CiliumClusterwideNetworkPolicy"
)

var (
	// log is the k8s package logger object.
	log = logging.DefaultLogger.WithField(logfields.LogSubsys, subsysK8s)
)

// GetPolicyLabels returns a LabelArray for the given namespace and name.
func GetPolicyLabels(ns, name string, uid types.UID, derivedFrom string) labels.LabelArray {
	// Keep labels sorted by the key.
	labelsArr := labels.LabelArray{
		labels.NewLabel(k8sConst.PolicyLabelDerivedFrom, derivedFrom, labels.LabelSourceK8s),
		labels.NewLabel(k8sConst.PolicyLabelName, name, labels.LabelSourceK8s),
	}

	// For clusterwide policy namespace will be empty.
	if ns != "" {
		nsLabel := labels.NewLabel(k8sConst.PolicyLabelNamespace, ns, labels.LabelSourceK8s)
		labelsArr = append(labelsArr, nsLabel)
	}

	srcLabel := labels.NewLabel(k8sConst.PolicyLabelUID, string(uid), labels.LabelSourceK8s)
	return append(labelsArr, srcLabel)
}

// getEndpointSelector converts the provided labelSelector into an EndpointSelector,
// adding the relevant matches for namespaces based on the provided options.
// If no namespace is provided then it is assumed that the selector is global to the cluster
// this is when translating selectors for CiliumClusterwideNetworkPolicy.
func getEndpointSelector(namespace string, labelSelector *slim_metav1.LabelSelector, addK8sPrefix, matchesInit bool) api.EndpointSelector {
	es := api.NewESFromK8sLabelSelector("", labelSelector)

	// The k8s prefix must not be added to reserved labels.
	if addK8sPrefix && es.HasKeyPrefix(labels.LabelSourceReservedKeyPrefix) {
		return es
	}

	// The user can explicitly specify the namespace in the
	// FromEndpoints selector. If omitted, we limit the
	// scope to the namespace the policy lives in.
	//
	// Policies applying on initializing pods are a special case.
	// Those pods don't have any labels, so they don't have a namespace label either.
	// Don't add a namespace label to those endpoint selectors, or we wouldn't be
	// able to match on those pods.
	if !es.HasKey(podPrefixLbl) && !es.HasKey(podAnyPrefixLbl) {
		if namespace == "" {
			// For a clusterwide policy if a namespace is not specified in the labels we add
			// a selector to only match endpoints that contains a namespace label.
			// This is to make sure that we are only allowing traffic for cilium managed k8s endpoints
			// and even if a wildcard is provided in the selector we don't proceed with a truly
			// empty(allow all) endpoint selector for the policy.
			if !matchesInit {
				es.AddMatchExpression(podPrefixLbl, slim_metav1.LabelSelectorOpExists, []string{})
			}
		} else if !es.HasKeyPrefix(podK8SNamespaceLabelsPrefix) && !es.HasKeyPrefix(podAnyNamespaceLabelsPrefix) {
			es.AddMatch(podPrefixLbl, namespace)
		}
	}

	return es
}

func parseToCiliumIngressCommonRule(namespace string, es api.EndpointSelector, ing api.IngressCommonRule) api.IngressCommonRule {
	matchesInit := matchesPodInit(es)
	var retRule api.IngressCommonRule

	if ing.FromEndpoints != nil {
		retRule.FromEndpoints = make([]api.EndpointSelector, len(ing.FromEndpoints))
		for j, ep := range ing.FromEndpoints {
			retRule.FromEndpoints[j] = getEndpointSelector(namespace, ep.LabelSelector, true, matchesInit)
		}
	}

	if ing.FromNodes != nil {
		retRule.FromNodes = make([]api.EndpointSelector, len(ing.FromNodes))
		for j, node := range ing.FromNodes {
			es = api.NewESFromK8sLabelSelector("", node.LabelSelector)
			es.AddMatchExpression(labels.LabelSourceReservedKeyPrefix+labels.IDNameRemoteNode, slim_metav1.LabelSelectorOpExists, []string{})
			retRule.FromNodes[j] = es
		}
	}

	if ing.FromCIDR != nil {
		retRule.FromCIDR = make([]api.CIDR, len(ing.FromCIDR))
		copy(retRule.FromCIDR, ing.FromCIDR)
	}

	if ing.FromCIDRSet != nil {
		retRule.FromCIDRSet = make([]api.CIDRRule, len(ing.FromCIDRSet))
		copy(retRule.FromCIDRSet, ing.FromCIDRSet)
	}

	if ing.FromRequires != nil {
		retRule.FromRequires = make([]api.EndpointSelector, len(ing.FromRequires))
		for j, ep := range ing.FromRequires {
			retRule.FromRequires[j] = getEndpointSelector(namespace, ep.LabelSelector, false, matchesInit)
		}
	}

	if ing.FromEntities != nil {
		retRule.FromEntities = make([]api.Entity, len(ing.FromEntities))
		copy(retRule.FromEntities, ing.FromEntities)
	}

	if ing.FromGroups != nil {
		retRule.FromGroups = make([]api.Groups, len(ing.FromGroups))
		copy(retRule.FromGroups, ing.FromGroups)
	}

	return retRule
}

func parseToCiliumIngressRule(namespace string, es api.EndpointSelector, inRules []api.IngressRule) []api.IngressRule {
	var retRules []api.IngressRule

	if inRules != nil {
		retRules = make([]api.IngressRule, len(inRules))
		for i, ing := range inRules {
			if ing.ToPorts != nil {
				retRules[i].ToPorts = make([]api.PortRule, len(ing.ToPorts))
				copy(retRules[i].ToPorts, ing.ToPorts)
			}
			if ing.ICMPs != nil {
				retRules[i].ICMPs = make(api.ICMPRules, len(ing.ICMPs))
				copy(retRules[i].ICMPs, ing.ICMPs)
			}
			retRules[i].IngressCommonRule = parseToCiliumIngressCommonRule(namespace, es, ing.IngressCommonRule)
			retRules[i].Authentication = ing.Authentication.DeepCopy()
			retRules[i].SetAggregatedSelectors()
		}
	}
	return retRules
}

func parseToCiliumIngressDenyRule(namespace string, es api.EndpointSelector, inRules []api.IngressDenyRule) []api.IngressDenyRule {
	var retRules []api.IngressDenyRule

	if inRules != nil {
		retRules = make([]api.IngressDenyRule, len(inRules))
		for i, ing := range inRules {
			if ing.ToPorts != nil {
				retRules[i].ToPorts = make([]api.PortDenyRule, len(ing.ToPorts))
				copy(retRules[i].ToPorts, ing.ToPorts)
			}
			if ing.ICMPs != nil {
				retRules[i].ICMPs = make(api.ICMPRules, len(ing.ICMPs))
				copy(retRules[i].ICMPs, ing.ICMPs)
			}
			retRules[i].IngressCommonRule = parseToCiliumIngressCommonRule(namespace, es, ing.IngressCommonRule)
			retRules[i].SetAggregatedSelectors()
		}
	}
	return retRules
}

func parseToCiliumEgressCommonRule(namespace string, es api.EndpointSelector, egr api.EgressCommonRule) api.EgressCommonRule {
	matchesInit := matchesPodInit(es)
	var retRule api.EgressCommonRule
	if egr.ToEndpoints != nil {
		retRule.ToEndpoints = make([]api.EndpointSelector, len(egr.ToEndpoints))
		for j, ep := range egr.ToEndpoints {
			endpointSelector := getEndpointSelector(namespace, ep.LabelSelector, true, matchesInit)
			endpointSelector.Generated = ep.Generated
			retRule.ToEndpoints[j] = endpointSelector
		}
	}

	if egr.ToCIDR != nil {
		retRule.ToCIDR = make([]api.CIDR, len(egr.ToCIDR))
		copy(retRule.ToCIDR, egr.ToCIDR)
	}

	if egr.ToCIDRSet != nil {
		retRule.ToCIDRSet = make(api.CIDRRuleSlice, len(egr.ToCIDRSet))
		copy(retRule.ToCIDRSet, egr.ToCIDRSet)
	}

	if egr.ToRequires != nil {
		retRule.ToRequires = make([]api.EndpointSelector, len(egr.ToRequires))
		for j, ep := range egr.ToRequires {
			retRule.ToRequires[j] = getEndpointSelector(namespace, ep.LabelSelector, false, matchesInit)
		}
	}

	if egr.ToServices != nil {
		retRule.ToServices = make([]api.Service, len(egr.ToServices))
		copy(retRule.ToServices, egr.ToServices)
	}

	if egr.ToEntities != nil {
		retRule.ToEntities = make([]api.Entity, len(egr.ToEntities))
		copy(retRule.ToEntities, egr.ToEntities)
	}

	if egr.ToNodes != nil {
		retRule.ToNodes = make([]api.EndpointSelector, len(egr.ToNodes))
		for j, node := range egr.ToNodes {
			es = api.NewESFromK8sLabelSelector("", node.LabelSelector)
			es.AddMatchExpression(labels.LabelSourceReservedKeyPrefix+labels.IDNameRemoteNode, slim_metav1.LabelSelectorOpExists, []string{})
			retRule.ToNodes[j] = es
		}
	}

	if egr.ToGroups != nil {
		retRule.ToGroups = make([]api.Groups, len(egr.ToGroups))
		copy(retRule.ToGroups, egr.ToGroups)
	}

	return retRule
}

func parseToCiliumEgressRule(namespace string, es api.EndpointSelector, inRules []api.EgressRule) []api.EgressRule {
	var retRules []api.EgressRule

	if inRules != nil {
		retRules = make([]api.EgressRule, len(inRules))
		for i, egr := range inRules {
			if egr.ToPorts != nil {
				retRules[i].ToPorts = make([]api.PortRule, len(egr.ToPorts))
				copy(retRules[i].ToPorts, egr.ToPorts)
			}

			if egr.ICMPs != nil {
				retRules[i].ICMPs = make(api.ICMPRules, len(egr.ICMPs))
				copy(retRules[i].ICMPs, egr.ICMPs)
			}

			if egr.ToFQDNs != nil {
				retRules[i].ToFQDNs = make([]api.FQDNSelector, len(egr.ToFQDNs))
				copy(retRules[i].ToFQDNs, egr.ToFQDNs)
			}

			retRules[i].EgressCommonRule = parseToCiliumEgressCommonRule(namespace, es, egr.EgressCommonRule)
			retRules[i].Authentication = egr.Authentication
			retRules[i].SetAggregatedSelectors()
		}
	}
	return retRules
}

func parseToCiliumEgressDenyRule(namespace string, es api.EndpointSelector, inRules []api.EgressDenyRule) []api.EgressDenyRule {
	var retRules []api.EgressDenyRule

	if inRules != nil {
		retRules = make([]api.EgressDenyRule, len(inRules))
		for i, egr := range inRules {
			if egr.ToPorts != nil {
				retRules[i].ToPorts = make([]api.PortDenyRule, len(egr.ToPorts))
				copy(retRules[i].ToPorts, egr.ToPorts)
			}

			if egr.ICMPs != nil {
				retRules[i].ICMPs = make(api.ICMPRules, len(egr.ICMPs))
				copy(retRules[i].ICMPs, egr.ICMPs)
			}

			retRules[i].EgressCommonRule = parseToCiliumEgressCommonRule(namespace, es, egr.EgressCommonRule)
			retRules[i].SetAggregatedSelectors()
		}
	}
	return retRules
}

func matchesPodInit(epSelector api.EndpointSelector) bool {
	if epSelector.LabelSelector == nil {
		return false
	}
	return epSelector.HasKey(podInitLbl)
}

// namespacesAreValid checks the set of namespaces from a rule returns true if
// they are not specified, or if they are specified and match the namespace
// where the rule is being inserted.
func namespacesAreValid(namespace string, userNamespaces []string) bool {
	return len(userNamespaces) == 0 ||
		(len(userNamespaces) == 1 && userNamespaces[0] == namespace)
}

// ParseToCiliumRule returns an api.Rule with all the labels parsed into cilium
// labels. If the namespace provided is empty then the rule is cluster scoped, this
// might happen in case of CiliumClusterwideNetworkPolicy which enforces a policy on the cluster
// instead of the particular namespace.
func ParseToCiliumRule(namespace, name string, uid types.UID, r *api.Rule) *api.Rule {
	retRule := &api.Rule{}
	if r.EndpointSelector.LabelSelector != nil {
		retRule.EndpointSelector = api.NewESFromK8sLabelSelector("", r.EndpointSelector.LabelSelector)
		// The PodSelector should only reflect to the same namespace
		// the policy is being stored, thus we add the namespace to
		// the MatchLabels map. Additionally, Policy repository relies
		// on this fact to properly choose correct network policies for
		// a given Security Identity.
		//
		// Policies applying to all namespaces are a special case.
		// Such policies can match on any traffic from Pods or Nodes,
		// so it wouldn't make sense to inject a namespace match for
		// those policies.
		if namespace != "" {
			userNamespace, present := r.EndpointSelector.GetMatch(podPrefixLbl)
			if present && !namespacesAreValid(namespace, userNamespace) {
				log.WithFields(logrus.Fields{
					logfields.K8sNamespace:              namespace,
					logfields.CiliumNetworkPolicyName:   name,
					logfields.K8sNamespace + ".illegal": userNamespace,
				}).Warn("CiliumNetworkPolicy contains illegal namespace match in EndpointSelector." +
					" EndpointSelector always applies in namespace of the policy resource, removing illegal namespace match'.")
			}
			retRule.EndpointSelector.AddMatch(podPrefixLbl, namespace)
		}
	} else if r.NodeSelector.LabelSelector != nil {
		retRule.NodeSelector = api.NewESFromK8sLabelSelector("", r.NodeSelector.LabelSelector)
	}

	retRule.Ingress = parseToCiliumIngressRule(namespace, r.EndpointSelector, r.Ingress)
	retRule.IngressDeny = parseToCiliumIngressDenyRule(namespace, r.EndpointSelector, r.IngressDeny)
	retRule.Egress = parseToCiliumEgressRule(namespace, r.EndpointSelector, r.Egress)
	retRule.EgressDeny = parseToCiliumEgressDenyRule(namespace, r.EndpointSelector, r.EgressDeny)

	retRule.Labels = ParseToCiliumLabels(namespace, name, uid, r.Labels)

	retRule.Description = r.Description
	retRule.EnableDefaultDeny = r.EnableDefaultDeny

	return retRule
}

// ParseToCiliumLabels returns all ruleLbls appended with a specific label that
// represents the given namespace and name along with a label that specifies
// these labels were derived from a CiliumNetworkPolicy.
func ParseToCiliumLabels(namespace, name string, uid types.UID, ruleLbs labels.LabelArray) labels.LabelArray {
	resourceType := ResourceTypeCiliumNetworkPolicy
	if namespace == "" {
		resourceType = ResourceTypeCiliumClusterwideNetworkPolicy
	}

	policyLbls := GetPolicyLabels(namespace, name, uid, resourceType)
	return append(policyLbls, ruleLbs...).Sort()
}
