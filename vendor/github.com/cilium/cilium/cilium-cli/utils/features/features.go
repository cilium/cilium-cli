// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package features

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/blang/semver/v4"
	v1 "k8s.io/api/core/v1"

	ciliumdef "github.com/cilium/cilium/pkg/defaults"
	"github.com/cilium/cilium/pkg/versioncheck"
)

const (
	CNIChaining        Feature = "cni-chaining"
	MonitorAggregation Feature = "monitor-aggregation"
	L7Proxy            Feature = "l7-proxy"
	HostFirewall       Feature = "host-firewall"
	ICMPPolicy         Feature = "icmp-policy"
	PortRanges         Feature = "port-ranges"
	L7PortRanges       Feature = "l7-port-ranges"
	Tunnel             Feature = "tunnel"
	TunnelPort         Feature = "tunnel-port"
	EndpointRoutes     Feature = "endpoint-routes"

	KPRMode                 Feature = "kpr-mode"
	KPRExternalIPs          Feature = "kpr-external-ips"
	KPRHostPort             Feature = "kpr-hostport"
	KPRSocketLB             Feature = "kpr-socket-lb"
	KPRSocketLBHostnsOnly   Feature = "kpr-socket-lb-hostns-only"
	KPRNodePort             Feature = "kpr-nodeport"
	KPRNodePortAcceleration Feature = "kpr-nodeport-acceleration"
	KPRSessionAffinity      Feature = "kpr-session-affinity"

	BPFLBExternalClusterIP Feature = "bpf-lb-external-clusterip"

	HostPort Feature = "host-port"

	NodeWithoutCilium Feature = "node-without-cilium"

	HealthChecking Feature = "health-checking"

	EncryptionPod        Feature = "encryption-pod"
	EncryptionNode       Feature = "encryption-node"
	EncryptionStrictMode Feature = "enable-encryption-strict-mode"

	IPv4 Feature = "ipv4"
	IPv6 Feature = "ipv6"

	Flavor Feature = "flavor"

	// The following settings control Policy Secrets tests.
	//
	// Cilium can be in three states for Policy Secrets:
	//
	// * Policy Secrets can be read by the agent from anywhere in the cluster (via either
	//   direct read or from the configured secret namespace via secret
	//   synchonrization by the Cilium Operator).
	// * Policy Secrets can be read by the agent, but only from the configured Secrets
	//   namespace. This is an advanced use case, and is included for migration purposes.
	// * Policy Secrets cannot be read.

	// PolicySecretsOnlyFromSecretsNamespace sets if Cilium  will look only
	// in the configured secrets namespace for Policy Secrets, or if it will look
	// in the entire cluster.
	//
	// If it's `true`, then Cilium will only read Secrets from the configured namespace.
	//
	// If it's `false`, then the Cilium agent will be granted Read access to _all_ Secrets
	// in the cluster.
	//
	// This feature replaces the existing `tls.secretsBackend: k8s` one. SecretsBackend
	// will be removed in a future release.
	//
	// This feature has Helm automation to mirror the setting of secretsBackend in the meantime.
	PolicySecretsOnlyFromSecretsNamespace Feature = "policy-secrets-only-from-secrets-namespace"

	// PolicySecretSync controls whether the Cilium Operator will synchronize Secrets referenced
	// in Network Policy into the configured Secrets namespace.
	//
	// This has important interactions with
	PolicySecretSync Feature = "enable-policy-secrets-sync"
	// For connectivity tests, we only care if Secrets can be read from the cluster
	// _somehow_, whether that is via direct read or secret sync is not important.
	// So, this feature tracks if we can read Policy secrets _somehow_.
	PolicySecretsReadable Feature = "policy-secrets-readable"

	CNP  Feature = "cilium-network-policy"
	CCNP Feature = "cilium-clusterwide-network-policy"
	KNP  Feature = "k8s-network-policy"

	// Whether or not CIDR selectors can match node IPs
	CIDRMatchNodes Feature = "cidr-match-nodes"

	AuthSpiffe Feature = "mutual-auth-spiffe"

	IngressController Feature = "ingress-controller"

	EgressGateway Feature = "enable-egress-gateway"
	GatewayAPI    Feature = "enable-gateway-api"

	EnableEnvoyConfig Feature = "enable-envoy-config"

	WireguardEncapsulate Feature = "wireguard-encapsulate"

	CiliumIPAMMode Feature = "ipam"

	IPsecEnabled                  Feature = "enable-ipsec"
	ClusterMeshEnableEndpointSync Feature = "clustermesh-enable-endpoint-sync"

	PolicyDefaultLocalCLuster Feature = "policy-default-local-cluster"

	LocalRedirectPolicy Feature = "enable-local-redirect-policy"

	BGPControlPlane Feature = "enable-bgp-control-plane"

	NodeLocalDNS Feature = "node-local-dns"

	Multicast Feature = "multicast-enabled"

	L7LoadBalancer Feature = "loadbalancer-l7"
)

// Feature is the name of a Cilium Feature (e.g. l7-proxy, cni chaining mode etc)
type Feature string

// Status describes the status of a Feature. Some features are either
// turned on or off (c.f. Enabled), while others additionally might include a
// Mode string which provides more information about in what mode a
// particular Feature is running ((e.g. when running with CNI chaining,
// Enabled will be true, and the Mode string will additionally contain the name
// of the chained CNI).
type Status struct {
	Enabled bool
	Mode    string
}

func (s Status) String() string {
	str := "Disabled"
	if s.Enabled {
		str = "Enabled"
	}

	if len(s.Mode) == 0 {
		return str
	}

	return fmt.Sprintf("%s:%s", str, s.Mode)
}

// Set contains the Status of a collection of Features.
type Set map[Feature]Status

// MatchRequirements returns true if the Set fs satisfies all the
// requirements in reqs. Returns true for empty requirements list.
func (fs Set) MatchRequirements(reqs ...Requirement) (bool, string) {
	for _, req := range reqs {
		status := fs[req.Feature]
		if req.requiresEnabled && (req.enabled != status.Enabled) {
			return false, fmt.Sprintf("Feature %s is disabled", req.Feature)
		}
		if req.requiresMode && (req.mode != status.Mode) {
			return false, fmt.Sprintf("requires Feature %s mode %s, got %s", req.Feature, req.mode, status.Mode)
		}
		if req.requireModeIsNot && (req.mode == status.Mode) {
			return false, fmt.Sprintf("requires Feature %s not equal to %s", req.Feature, req.mode)
		}
	}

	return true, ""
}

// IPFamilies returns the list of enabled IP families.
func (fs Set) IPFamilies() []IPFamily {
	var families []IPFamily

	if match, _ := fs.MatchRequirements(RequireEnabled(IPv4)); match {
		families = append(families, IPFamilyV4)
	}

	if match, _ := fs.MatchRequirements(RequireEnabled(IPv6)); match {
		families = append(families, IPFamilyV6)
	}

	return families
}

// deriveFeatures derives additional features based on the status of other features
func (fs Set) DeriveFeatures() error {
	fs[HostPort] = Status{
		// HostPort support can either be enabled via KPR, or via CNI chaining with portmap plugin
		Enabled: (fs[CNIChaining].Enabled && fs[CNIChaining].Mode == "portmap" &&
			// cilium/cilium#12541: Host firewall doesn't work with portmap CNI chaining
			!fs[HostFirewall].Enabled) ||
			fs[KPRHostPort].Enabled,
	}

	return nil
}

// Requirement defines a test requirement. A given Set may or
// may not satisfy this requirement
type Requirement struct {
	Feature Feature

	requiresEnabled bool
	enabled         bool

	requiresMode     bool
	requireModeIsNot bool
	mode             string
}

// RequireEnabled constructs a Requirement which expects the
// Feature to be enabled
func RequireEnabled(feature Feature) Requirement {
	return Requirement{
		Feature:         feature,
		requiresEnabled: true,
		enabled:         true,
	}
}

// RequireDisabled constructs a Requirement which expects the
// Feature to be disabled
func RequireDisabled(feature Feature) Requirement {
	return Requirement{
		Feature:         feature,
		requiresEnabled: true,
		enabled:         false,
	}
}

// RequireMode constructs a Requirement which expects the Feature
// to be in the given mode
func RequireMode(feature Feature, mode string) Requirement {
	return Requirement{
		Feature:      feature,
		requiresMode: true,
		mode:         mode,
	}
}

// RequiredModeIsNot constructs a Requirement which expects the Feature to not
// be in the given mode
//
// When evaluating a set of requirements with MatchRequirements,
// having a RequireMode requirement of the same feature and mode will cause
// conflicting results.
func RequireModeIsNot(feature Feature, mode string) Requirement {
	return Requirement{
		Feature:          feature,
		requireModeIsNot: true,
		mode:             mode,
	}
}

// ExtractFromVersionedConfigMap extracts features based on Cilium version and cilium-config
// ConfigMap.
func (fs Set) ExtractFromVersionedConfigMap(ciliumVersion semver.Version, cm *v1.ConfigMap) {
	fs[PortRanges] = ExtractPortRanges(ciliumVersion)
	fs[L7PortRanges] = ExtractL7PortRanges(ciliumVersion)
}

func ExtractPortRanges(ciliumVersion semver.Version) Status {
	enabled := versioncheck.MustCompile(">=1.16.0")(ciliumVersion)
	return Status{
		Enabled: enabled,
	}
}

func ExtractL7PortRanges(ciliumVersion semver.Version) Status {
	enabled := versioncheck.MustCompile(">=1.17.0")(ciliumVersion)
	return Status{
		Enabled: enabled,
	}
}

func ExtractTunnelFeatureFromConfigMap(cm *v1.ConfigMap) (Status, Status) {
	getTunnelPortFeature := func(tunnelProtocol string) Status {
		tunnelPort, ok := cm.Data["tunnel-port"]
		switch {
		case !ok && tunnelProtocol == "vxlan":
			tunnelPort = fmt.Sprintf("%d", ciliumdef.TunnelPortVXLAN)
		case !ok && tunnelProtocol == "geneve":
			tunnelPort = fmt.Sprintf("%d", ciliumdef.TunnelPortGeneve)
		}
		return Status{
			Enabled: ok,
			Mode:    tunnelPort,
		}
	}

	mode := "tunnel"
	if v, ok := cm.Data["routing-mode"]; ok {
		mode = v
	}

	tunnelProto := "vxlan"
	if v, ok := cm.Data["tunnel-protocol"]; ok {
		tunnelProto = v
	}

	return Status{
		Enabled: mode != "native",
		Mode:    tunnelProto,
	}, getTunnelPortFeature(tunnelProto)
}

// ExtractFromConfigMap extracts features from the Cilium ConfigMap.
// Note that there is no rule regarding if the default value is reflected
// in the ConfigMap or not.
func (fs Set) ExtractFromConfigMap(cm *v1.ConfigMap) {
	// CNI chaining.
	// Note: This value might be overwritten by extractFeaturesFromCiliumStatus
	// if this information is present in `cilium status`
	mode := "none"
	if v, ok := cm.Data["cni-chaining-mode"]; ok {
		mode = v
	}
	fs[CNIChaining] = Status{
		Enabled: mode != "none",
		Mode:    mode,
	}

	fs[IPv4] = Status{
		Enabled: cm.Data["enable-ipv4"] == "true",
	}
	fs[IPv6] = Status{
		Enabled: cm.Data["enable-ipv6"] == "true",
	}

	fs[EndpointRoutes] = Status{
		Enabled: cm.Data["enable-endpoint-routes"] == "true",
	}

	fs[AuthSpiffe] = Status{
		Enabled: cm.Data["mesh-auth-mutual-enabled"] == "true",
	}

	fs[IngressController] = Status{
		Enabled: cm.Data["enable-ingress-controller"] == "true",
	}

	fs[EgressGateway] = Status{
		Enabled: cm.Data[string(EgressGateway)] == "true" || cm.Data["enable-ipv4-egress-gateway"] == "true",
	}

	fs[CIDRMatchNodes] = Status{
		Enabled: strings.Contains(cm.Data["policy-cidr-match-mode"], "nodes"),
	}

	fs[GatewayAPI] = Status{
		Enabled: cm.Data[string(GatewayAPI)] == "true",
	}

	fs[EnableEnvoyConfig] = Status{
		Enabled: cm.Data[string(EnableEnvoyConfig)] == "true",
	}

	fs[WireguardEncapsulate] = Status{
		Enabled: cm.Data[string(WireguardEncapsulate)] == "true",
	}

	fs[CiliumIPAMMode] = Status{
		Mode: cm.Data[string(CiliumIPAMMode)],
	}

	fs[IPsecEnabled] = Status{
		Enabled: cm.Data[string(IPsecEnabled)] == "true",
	}

	fs[ClusterMeshEnableEndpointSync] = Status{
		Enabled: cm.Data[string(ClusterMeshEnableEndpointSync)] == "true",
	}

	fs[PolicyDefaultLocalCLuster] = Status{
		Enabled: cm.Data[string(PolicyDefaultLocalCLuster)] == "true",
	}

	fs[LocalRedirectPolicy] = Status{
		Enabled: cm.Data[string(LocalRedirectPolicy)] == "true",
	}

	fs[BPFLBExternalClusterIP] = Status{
		Enabled: cm.Data[string(BPFLBExternalClusterIP)] == "true",
	}

	fs[BGPControlPlane] = Status{
		Enabled: cm.Data[string(BGPControlPlane)] == "true",
	}

	fs[Multicast] = Status{
		Enabled: cm.Data[string(Multicast)] == "true",
	}

	fs[EncryptionStrictMode] = Status{
		Enabled: cm.Data[string(EncryptionStrictMode)] == "true",
	}

	// This could be enabled via ClusterRole check as well, so only
	// check if it's false.
	if !fs[PolicySecretsOnlyFromSecretsNamespace].Enabled {
		fs[PolicySecretsOnlyFromSecretsNamespace] = Status{
			Enabled: cm.Data[string(PolicySecretSync)] == "true",
		}
	}

	fs[PolicySecretSync] = Status{
		Enabled: cm.Data[string(PolicySecretSync)] == "true",
	}

	fs[L7LoadBalancer] = Status{
		Enabled: cm.Data[string(L7LoadBalancer)] == "envoy",
	}

	fs[Tunnel], fs[TunnelPort] = ExtractTunnelFeatureFromConfigMap(cm)
}

func (fs Set) ExtractFromNodes(nodesWithoutCilium map[string]struct{}) {
	fs[NodeWithoutCilium] = Status{
		Enabled: len(nodesWithoutCilium) != 0,
		Mode:    strings.Join(slices.Collect(maps.Keys(nodesWithoutCilium)), ","),
	}
}
