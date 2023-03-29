// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package sysdump

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/blang/semver/v4"
	"github.com/cilium/cilium/api/v1/models"
	ciliumv2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumv2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	tetragonv1alpha1 "github.com/cilium/tetragon/pkg/k8s/apis/cilium.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/k8s"

	apiserverv1 "github.com/openshift/api/apiserver/v1"
	openshiftAppsv1 "github.com/openshift/api/apps/v1"
	cloudnetworkv1 "github.com/openshift/api/cloudnetwork/v1"
	configv1 "github.com/openshift/api/config/v1"
	imagev1 "github.com/openshift/api/image/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	networkv1 "github.com/openshift/api/network/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	projectv1 "github.com/openshift/api/project/v1"
	quotav1 "github.com/openshift/api/quota/v1"
	routev1 "github.com/openshift/api/route/v1"
	samplesv1 "github.com/openshift/api/samples/v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type SysdumpSuite struct{}

var _ = check.Suite(&SysdumpSuite{})

func (b *SysdumpSuite) TestSysdumpCollector(c *check.C) {
	client := fakeClient{
		nodeList: &corev1.NodeList{
			Items: []corev1.Node{
				{ObjectMeta: metav1.ObjectMeta{Name: "node-a"}},
			},
		},
	}
	options := Options{
		OutputFileName: "my-sysdump-<ts>",
		Writer:         io.Discard,
	}
	startTime := time.Unix(946713600, 0)
	timestamp := startTime.Format(timeFormat)
	collector, err := NewCollector(&client, options, startTime, "cilium-cli-version")
	c.Assert(err, check.IsNil)
	c.Assert(path.Base(collector.sysdumpDir), check.Equals, "my-sysdump-"+timestamp)
	tempFile := collector.AbsoluteTempPath("my-file-<ts>")
	c.Assert(tempFile, check.Equals, path.Join(collector.sysdumpDir, "my-file-"+timestamp))
	_, err = os.Stat(path.Join(collector.sysdumpDir, sysdumpLogFile))
	c.Assert(err, check.IsNil)
}

func (b *SysdumpSuite) TestNodeList(c *check.C) {
	options := Options{
		Writer: io.Discard,
	}
	client := fakeClient{
		nodeList: &corev1.NodeList{
			Items: []corev1.Node{
				{ObjectMeta: metav1.ObjectMeta{Name: "node-a"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "node-b"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "node-c"}},
			},
		},
	}
	collector, err := NewCollector(&client, options, time.Now(), "cilium-cli-version")
	c.Assert(err, check.IsNil)
	c.Assert(collector.NodeList, check.DeepEquals, []string{"node-a", "node-b", "node-c"})

	options = Options{
		Writer:   io.Discard,
		NodeList: "node-a,node-c",
	}
	collector, err = NewCollector(&client, options, time.Now(), "cilium-cli-version")
	c.Assert(err, check.IsNil)
	c.Assert(collector.NodeList, check.DeepEquals, []string{"node-a", "node-c"})
}

func (b *SysdumpSuite) TestAddTasks(c *check.C) {
	options := Options{
		Writer: io.Discard,
	}
	client := fakeClient{
		nodeList: &corev1.NodeList{
			Items: []corev1.Node{
				{ObjectMeta: metav1.ObjectMeta{Name: "node-a"}},
			},
		},
	}
	collector, err := NewCollector(&client, options, time.Now(), "cilium-cli-version")
	c.Assert(err, check.IsNil)
	collector.AddTasks([]Task{{}, {}, {}})
	c.Assert(len(collector.additionalTasks), check.Equals, 3)
	collector.AddTasks([]Task{{}, {}, {}})
	c.Assert(len(collector.additionalTasks), check.Equals, 6)
}

func (b *SysdumpSuite) TestExtractGopsPID(c *check.C) {
	var pid string
	var err error

	normalOutput := `
25863 0     gops          unknown Go version /usr/bin/gops
25852 25847 cilium        unknown Go version /usr/bin/cilium
10    1     cilium-agent* unknown Go version /usr/bin/cilium-agent
1     0     custom        go1.16.3           /usr/local/bin/custom
	`
	pid, err = extractGopsPID(normalOutput)
	c.Assert(err, check.IsNil)
	c.Assert(pid, check.Equals, "10")

	missingAgent := `
25863 0     gops          unknown Go version /usr/bin/gops
25852 25847 cilium        unknown Go version /usr/bin/cilium
10    1     cilium-agent unknown Go version /usr/bin/cilium-agent
1     0     custom        go1.16.3           /usr/local/bin/custom
	`
	pid, err = extractGopsPID(missingAgent)
	c.Assert(err, check.NotNil)
	c.Assert(pid, check.Equals, "")

	multipleAgents := `
25863 0     gops*          unknown Go version /usr/bin/gops
25852 25847 cilium*        unknown Go version /usr/bin/cilium
10    1     cilium-agent unknown Go version /usr/bin/cilium-agent
1     0     custom        go1.16.3           /usr/local/bin/custom
	`
	pid, err = extractGopsPID(multipleAgents)
	c.Assert(err, check.IsNil)
	c.Assert(pid, check.Equals, "25863")

	noOutput := ``
	_, err = extractGopsPID(noOutput)
	c.Assert(err, check.NotNil)

}

func (b *SysdumpSuite) TestExtractGopsProfileData(c *check.C) {
	gopsOutput := `
	Profiling CPU now, will take 30 secs...
	Profile dump saved to: /tmp/cpu_profile3302111893
	`
	wantFilepath := "/tmp/cpu_profile3302111893"

	gotFilepath, err := extractGopsProfileData(gopsOutput)
	c.Assert(err, check.IsNil)
	c.Assert(gotFilepath, check.Equals, wantFilepath)

}

func TestKVStoreTask(t *testing.T) {
	assert := assert.New(t)
	client := &fakeClient{
		nodeList: &corev1.NodeList{
			Items: []corev1.Node{{ObjectMeta: metav1.ObjectMeta{Name: "node-a"}}},
		},
		execs: make(map[execRequest]execResult),
	}
	addKVStoreGet := func(c *fakeClient, ciliumPaths ...string) {
		for _, path := range ciliumPaths {
			c.expectExec("ns0", "pod0", defaults.AgentContainerName,
				[]string{"cilium", "kvstore", "get", "cilium/" + path, "--recursive", "-o", "json"},
				[]byte("{}"), nil, nil)
		}
	}
	addKVStoreGet(client, "state/identities", "state/ip", "state/nodes", "state/cnpstatuses", ".heartbeat", "state/services")
	options := Options{
		OutputFileName: "my-sysdump-<ts>",
		Writer:         io.Discard,
	}
	collector, err := NewCollector(client, options, time.Now(), "cilium-cli-version")
	assert.NoError(err)
	collector.submitKVStoreTasks(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod0",
			Namespace: "ns0",
		},
	})
	fd, err := os.Open(path.Join(collector.sysdumpDir, "kvstore-heartbeat.json"))
	assert.NoError(err)
	data, err := io.ReadAll(fd)
	assert.NoError(err)
	assert.Equal([]byte("{}"), data)
}

func TestListTetragonTracingPolicies(t *testing.T) {
	assert := assert.New(t)
	client := &fakeClient{}

	tracingPolicies, err := client.ListTetragonTracingPolicies(context.Background(), metav1.ListOptions{})
	assert.NoError(err)
	assert.GreaterOrEqual(len(tracingPolicies.Items), 0)
}

func TestListCiliumEndpointSlices(t *testing.T) {
	assert := assert.New(t)
	client := &fakeClient{}

	endpointSlices, err := client.ListCiliumEndpointSlices(context.Background(), metav1.ListOptions{})
	assert.NoError(err)
	assert.GreaterOrEqual(len(endpointSlices.Items), 0)
}

func TestListCiliumExternalWorkloads(t *testing.T) {
	assert := assert.New(t)
	client := &fakeClient{}

	externalWorkloads, err := client.ListCiliumExternalWorkloads(context.Background(), metav1.ListOptions{})
	assert.NoError(err)
	assert.GreaterOrEqual(len(externalWorkloads.Items), 0)
}

type execRequest struct {
	namespace string
	pod       string
	container string
	command   string
}

type execResult struct {
	stderr []byte
	stdout []byte
	err    error
}

type fakeClient struct {
	nodeList *corev1.NodeList
	execs    map[execRequest]execResult
}

func (c *fakeClient) ListCiliumBGPPeeringPolicies(ctx context.Context, opts metav1.ListOptions) (*ciliumv2alpha1.CiliumBGPPeeringPolicyList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumLoadBalancerIPPools(ctx context.Context, opts metav1.ListOptions) (*ciliumv2alpha1.CiliumLoadBalancerIPPoolList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumNodeConfigs(ctx context.Context, namespace string, opts metav1.ListOptions) (*ciliumv2alpha1.CiliumNodeConfigList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumClusterwideEnvoyConfigs(ctx context.Context, opts metav1.ListOptions) (*ciliumv2.CiliumClusterwideEnvoyConfigList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumEnvoyConfigs(ctx context.Context, namespace string, options metav1.ListOptions) (*ciliumv2.CiliumEnvoyConfigList, error) {
	panic("implement me")
}

func (c *fakeClient) ListIngresses(ctx context.Context, o metav1.ListOptions) (*networkingv1.IngressList, error) {
	panic("implement me")
}

func (c *fakeClient) CopyFromPod(ctx context.Context, namespace, pod, container, fromFile, destFile string, retryLimit int) error {
	panic("implement me")
}

func (c *fakeClient) AutodetectFlavor(ctx context.Context) k8s.Flavor {
	panic("implement me")
}

func (c *fakeClient) GetPod(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Pod, error) {
	panic("implement me")
}

func (c *fakeClient) CreatePod(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	panic("implement me")
}

func (c *fakeClient) DeletePod(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error {
	panic("implement me")
}

func (c *fakeClient) expectExec(namespace, pod, container string, command []string, expectedStdout []byte, expectedStderr []byte, expectedErr error) {
	r := execRequest{namespace, pod, container, strings.Join(command, " ")}
	c.execs[r] = execResult{
		stdout: expectedStdout,
		stderr: expectedStderr,
		err:    expectedErr,
	}
}

func (c *fakeClient) ExecInPod(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, error) {
	stdout, _, err := c.ExecInPodWithStderr(ctx, namespace, pod, container, command)
	return stdout, err
}

func (c *fakeClient) ExecInPodWithStderr(ctx context.Context, namespace, pod, container string, command []string) (bytes.Buffer, bytes.Buffer, error) {
	r := execRequest{namespace, pod, container, strings.Join(command, " ")}
	out, ok := c.execs[r]
	if !ok {
		panic(fmt.Sprintf("unexpected exec: %v", r))
	}
	return *bytes.NewBuffer(out.stdout), *bytes.NewBuffer(out.stderr), out.err
}

func (c *fakeClient) GetCiliumVersion(ctx context.Context, p *corev1.Pod) (*semver.Version, error) {
	panic("implement me")
}

func (c *fakeClient) GetConfigMap(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.ConfigMap, error) {
	panic("implement me")
}

func (c *fakeClient) GetDaemonSet(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.DaemonSet, error) {
	return nil, nil
}

func (c *fakeClient) GetDeployment(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*appsv1.Deployment, error) {
	return nil, nil
}

func (c *fakeClient) GetLogs(ctx context.Context, namespace, name, container string, sinceTime time.Time, limitBytes int64, previous bool) (string, error) {
	panic("implement me")
}

func (c *fakeClient) GetPodsTable(ctx context.Context) (*metav1.Table, error) {
	panic("implement me")
}

func (c *fakeClient) GetSecret(ctx context.Context, namespace, name string, opts metav1.GetOptions) (*corev1.Secret, error) {
	panic("implement me")
}

func (c *fakeClient) GetVersion(ctx context.Context) (string, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumClusterwideNetworkPolicies(ctx context.Context, opts metav1.ListOptions) (*ciliumv2.CiliumClusterwideNetworkPolicyList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumIdentities(ctx context.Context) (*ciliumv2.CiliumIdentityList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumEgressGatewayPolicies(ctx context.Context, opts metav1.ListOptions) (*ciliumv2.CiliumEgressGatewayPolicyList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumEndpoints(ctx context.Context, namespace string, options metav1.ListOptions) (*ciliumv2.CiliumEndpointList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumEndpointSlices(ctx context.Context, options metav1.ListOptions) (*ciliumv2alpha1.CiliumEndpointSliceList, error) {
	ciliumEndpointSliceList := ciliumv2alpha1.CiliumEndpointSliceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		ListMeta: metav1.ListMeta{},
		Items: []ciliumv2alpha1.CiliumEndpointSlice{{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CiliumEndpointSlice",
				APIVersion: "v2alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "testEndpointSlice1",
			},
			Endpoints: []ciliumv2alpha1.CoreCiliumEndpoint{{
				Name:       "EndpointSlice1",
				IdentityID: 1,
				Networking: &ciliumv2.EndpointNetworking{
					Addressing: ciliumv2.AddressPairList{{
						IPV4: "10.0.0.1",
					},
						{
							IPV4: "10.0.0.2",
						},
					},
				},
				Encryption: ciliumv2.EncryptionSpec{},
				NamedPorts: models.NamedPorts{},
			},
			},
		},
		},
	}
	return &ciliumEndpointSliceList, nil
}

func (c *fakeClient) ListCiliumExternalWorkloads(ctx context.Context, options metav1.ListOptions) (*ciliumv2.CiliumExternalWorkloadList, error) {
	ciliumExternalWorkloadList := ciliumv2.CiliumExternalWorkloadList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		ListMeta: metav1.ListMeta{},
		Items: []ciliumv2.CiliumExternalWorkload{{
			TypeMeta: metav1.TypeMeta{
				Kind:       "CiliumEndpointSlice",
				APIVersion: "v2alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "testEndpointSlice1",
			},
			Spec: ciliumv2.CiliumExternalWorkloadSpec{
				IPv4AllocCIDR: "10.100.0.0/24",
				IPv6AllocCIDR: "FD00::/64",
			},
		},
		},
	}
	return &ciliumExternalWorkloadList, nil
}

func (c *fakeClient) ListCiliumLocalRedirectPolicies(ctx context.Context, namespace string, options metav1.ListOptions) (*ciliumv2.CiliumLocalRedirectPolicyList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumNetworkPolicies(ctx context.Context, namespace string, opts metav1.ListOptions) (*ciliumv2.CiliumNetworkPolicyList, error) {
	panic("implement me")
}

func (c *fakeClient) ListCiliumNodes(ctx context.Context) (*ciliumv2.CiliumNodeList, error) {
	panic("implement me")
}

func (c *fakeClient) ListDaemonSet(ctx context.Context, namespace string, o metav1.ListOptions) (*appsv1.DaemonSetList, error) {
	panic("implement me")
}

func (c *fakeClient) ListEvents(ctx context.Context, o metav1.ListOptions) (*corev1.EventList, error) {
	panic("implement me")
}

func (c *fakeClient) ListNamespaces(ctx context.Context, o metav1.ListOptions) (*corev1.NamespaceList, error) {
	panic("implement me")
}

func (c *fakeClient) ListEndpoints(ctx context.Context, o metav1.ListOptions) (*corev1.EndpointsList, error) {
	panic("implement me")
}

func (c *fakeClient) ListNetworkPolicies(ctx context.Context, o metav1.ListOptions) (*networkingv1.NetworkPolicyList, error) {
	panic("implement me")
}

func (c *fakeClient) ListNodes(ctx context.Context, options metav1.ListOptions) (*corev1.NodeList, error) {
	return c.nodeList, nil
}

func (c *fakeClient) ListPods(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.PodList, error) {
	panic("implement me")
}

func (c *fakeClient) ListServices(ctx context.Context, namespace string, options metav1.ListOptions) (*corev1.ServiceList, error) {
	panic("implement me")
}

func (c *fakeClient) ListUnstructured(ctx context.Context, gvr schema.GroupVersionResource, namespace *string, o metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	panic("implement me")
}

func (c *fakeClient) ListTetragonTracingPolicies(ctx context.Context, options metav1.ListOptions) (*tetragonv1alpha1.TracingPolicyList, error) {
	tetragonTracingPolicy := tetragonv1alpha1.TracingPolicyList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		ListMeta: metav1.ListMeta{},
		Items: []tetragonv1alpha1.TracingPolicy{{
			TypeMeta: metav1.TypeMeta{
				Kind:       "TracingPolicy",
				APIVersion: "v1alpha",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "testPolicy1",
			},
			Spec: tetragonv1alpha1.TracingPolicySpec{
				KProbes:     []tetragonv1alpha1.KProbeSpec{},
				Tracepoints: []tetragonv1alpha1.TracepointSpec{},
				Loader:      true,
			},
		},
		},
	}

	return &tetragonTracingPolicy, nil
}

func (c *fakeClient) CreateEphemeralContainer(ctx context.Context, pod *corev1.Pod, container *corev1.EphemeralContainer) (*corev1.Pod, error) {
	panic("implement me")
}

func (c *fakeClient) GetNamespace(_ context.Context, ns string, _ metav1.GetOptions) (*corev1.Namespace, error) {
	if ns == "kube-system" {
		return &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		}, nil
	}
	return nil, &errors.StatusError{
		ErrStatus: metav1.Status{
			Code: http.StatusNotFound,
		},
	}
}

func (c *fakeClient) ListOpenshiftAPIRequestCounts(ctx context.Context, o metav1.ListOptions) (*apiserverv1.APIRequestCountList, error) {
	return &apiserverv1.APIRequestCountList{}, nil
}

func (c *fakeClient) ListOpenshiftDeploymentConfigs(ctx context.Context, namespace string, o metav1.ListOptions) (*openshiftAppsv1.DeploymentConfigList, error) {
	return &openshiftAppsv1.DeploymentConfigList{}, nil
}

func (c *fakeClient) ListOpenshiftAPIServers(ctx context.Context, o metav1.ListOptions) (*configv1.APIServerList, error) {
	return &configv1.APIServerList{}, nil
}

func (c *fakeClient) ListOpenshiftBuilds(ctx context.Context, o metav1.ListOptions) (*configv1.BuildList, error) {
	return &configv1.BuildList{}, nil
}

func (c *fakeClient) ListOpenshiftClusterOperators(ctx context.Context, o metav1.ListOptions) (*configv1.ClusterOperatorList, error) {
	return &configv1.ClusterOperatorList{}, nil
}

func (c *fakeClient) ListOpenshiftClusterVersions(ctx context.Context, o metav1.ListOptions) (*configv1.ClusterVersionList, error) {
	return &configv1.ClusterVersionList{}, nil
}

func (c *fakeClient) ListOpenshiftConfigConsoles(ctx context.Context, o metav1.ListOptions) (*configv1.ConsoleList, error) {
	return &configv1.ConsoleList{}, nil
}

func (c *fakeClient) ListOpenshiftConfigDNSes(ctx context.Context, o metav1.ListOptions) (*configv1.DNSList, error) {
	return &configv1.DNSList{}, nil
}

func (c *fakeClient) ListOpenshiftFeatureGates(ctx context.Context, o metav1.ListOptions) (*configv1.FeatureGateList, error) {
	return &configv1.FeatureGateList{}, nil
}

func (c *fakeClient) ListOpenshiftConfigImages(ctx context.Context, o metav1.ListOptions) (*configv1.ImageList, error) {
	return &configv1.ImageList{}, nil
}

func (c *fakeClient) ListOpenshiftInfrastructures(ctx context.Context, o metav1.ListOptions) (*configv1.InfrastructureList, error) {
	return &configv1.InfrastructureList{}, nil
}

func (c *fakeClient) ListOpenshiftIngresses(ctx context.Context, o metav1.ListOptions) (*configv1.IngressList, error) {
	return &configv1.IngressList{}, nil
}

func (c *fakeClient) ListOpenshiftConfigNetworks(ctx context.Context, o metav1.ListOptions) (*configv1.NetworkList, error) {
	return &configv1.NetworkList{}, nil
}

func (c *fakeClient) ListOpenshiftNodes(ctx context.Context, o metav1.ListOptions) (*configv1.NodeList, error) {
	return &configv1.NodeList{}, nil
}

func (c *fakeClient) ListOpenshiftOperatorHubs(ctx context.Context, o metav1.ListOptions) (*configv1.OperatorHubList, error) {
	return &configv1.OperatorHubList{}, nil
}

func (c *fakeClient) ListOpenshiftConfigProjects(ctx context.Context, o metav1.ListOptions) (*configv1.ProjectList, error) {
	return &configv1.ProjectList{}, nil
}

func (c *fakeClient) ListOpenshiftProxies(ctx context.Context, o metav1.ListOptions) (*configv1.ProxyList, error) {
	return &configv1.ProxyList{}, nil
}

func (c *fakeClient) ListOpenshiftSchedulers(ctx context.Context, o metav1.ListOptions) (*configv1.SchedulerList, error) {
	return &configv1.SchedulerList{}, nil
}

func (c *fakeClient) ListOpenshiftCloudPrivateIPConfigs(ctx context.Context, o metav1.ListOptions) (*cloudnetworkv1.CloudPrivateIPConfigList, error) {
	return &cloudnetworkv1.CloudPrivateIPConfigList{}, nil
}

func (c *fakeClient) ListOpenshiftImages(ctx context.Context, o metav1.ListOptions) (*imagev1.ImageList, error) {
	return &imagev1.ImageList{}, nil
}

func (c *fakeClient) ListOpenshiftControlPlaneMachineSets(ctx context.Context, namespace string, o metav1.ListOptions) (*machinev1.ControlPlaneMachineSetList, error) {
	return &machinev1.ControlPlaneMachineSetList{}, nil
}

func (c *fakeClient) ListOpenshiftClusterNetworks(ctx context.Context, o metav1.ListOptions) (*networkv1.ClusterNetworkList, error) {
	return &networkv1.ClusterNetworkList{}, nil
}

func (c *fakeClient) ListOpenshiftEgressNetworkPolicies(ctx context.Context, namespace string, o metav1.ListOptions) (*networkv1.EgressNetworkPolicyList, error) {
	return &networkv1.EgressNetworkPolicyList{}, nil
}

func (c *fakeClient) ListOpenshiftHostSubnets(ctx context.Context, o metav1.ListOptions) (*networkv1.HostSubnetList, error) {
	return &networkv1.HostSubnetList{}, nil
}

func (c *fakeClient) ListOpenshiftNetNamespaces(ctx context.Context, o metav1.ListOptions) (*networkv1.NetNamespaceList, error) {
	return &networkv1.NetNamespaceList{}, nil
}

func (c *fakeClient) ListOpenshiftCSISnapshotControllers(ctx context.Context, o metav1.ListOptions) (*operatorv1.CSISnapshotControllerList, error) {
	return &operatorv1.CSISnapshotControllerList{}, nil
}

func (c *fakeClient) ListOpenshiftCloudCredentials(ctx context.Context, o metav1.ListOptions) (*operatorv1.CloudCredentialList, error) {
	return &operatorv1.CloudCredentialList{}, nil
}

func (c *fakeClient) ListOpenshiftClusterCSIDrivers(ctx context.Context, o metav1.ListOptions) (*operatorv1.ClusterCSIDriverList, error) {
	return &operatorv1.ClusterCSIDriverList{}, nil
}

func (c *fakeClient) ListOpenshiftConfigs(ctx context.Context, o metav1.ListOptions) (*operatorv1.ConfigList, error) {
	return &operatorv1.ConfigList{}, nil
}

func (c *fakeClient) ListOpenshiftConsoles(ctx context.Context, o metav1.ListOptions) (*operatorv1.ConsoleList, error) {
	return &operatorv1.ConsoleList{}, nil
}

func (c *fakeClient) ListOpenshiftDNSes(ctx context.Context, o metav1.ListOptions) (*operatorv1.DNSList, error) {
	return &operatorv1.DNSList{}, nil
}

func (c *fakeClient) ListOpenshiftEtcds(ctx context.Context, o metav1.ListOptions) (*operatorv1.EtcdList, error) {
	return &operatorv1.EtcdList{}, nil
}

func (c *fakeClient) ListOpenshiftIngressControllers(ctx context.Context, namespace string, o metav1.ListOptions) (*operatorv1.IngressControllerList, error) {
	return &operatorv1.IngressControllerList{}, nil
}

func (c *fakeClient) ListOpenshiftInsightsOperators(ctx context.Context, o metav1.ListOptions) (*operatorv1.InsightsOperatorList, error) {
	return &operatorv1.InsightsOperatorList{}, nil
}

func (c *fakeClient) ListOpenshiftKubeAPIServers(ctx context.Context, o metav1.ListOptions) (*operatorv1.KubeAPIServerList, error) {
	return &operatorv1.KubeAPIServerList{}, nil
}

func (c *fakeClient) ListOpenshiftKubeControllerManagers(ctx context.Context, o metav1.ListOptions) (*operatorv1.KubeControllerManagerList, error) {
	return &operatorv1.KubeControllerManagerList{}, nil
}

func (c *fakeClient) ListOpenshiftKubeSchedulers(ctx context.Context, o metav1.ListOptions) (*operatorv1.KubeSchedulerList, error) {
	return &operatorv1.KubeSchedulerList{}, nil
}

func (c *fakeClient) ListOpenshiftKubeStorageVersionMigrators(ctx context.Context, o metav1.ListOptions) (*operatorv1.KubeStorageVersionMigratorList, error) {
	return &operatorv1.KubeStorageVersionMigratorList{}, nil
}

func (c *fakeClient) ListOpenshiftNetworks(ctx context.Context, o metav1.ListOptions) (*operatorv1.NetworkList, error) {
	return &operatorv1.NetworkList{}, nil
}

func (c *fakeClient) ListOpenshiftOperatorAPIServers(ctx context.Context, o metav1.ListOptions) (*operatorv1.OpenShiftAPIServerList, error) {
	return &operatorv1.OpenShiftAPIServerList{}, nil
}

func (c *fakeClient) ListOpenshiftOperatorControllerManagers(ctx context.Context, o metav1.ListOptions) (*operatorv1.OpenShiftControllerManagerList, error) {
	return &operatorv1.OpenShiftControllerManagerList{}, nil
}

func (c *fakeClient) ListOpenshiftServiceCAs(ctx context.Context, o metav1.ListOptions) (*operatorv1.ServiceCAList, error) {
	return &operatorv1.ServiceCAList{}, nil
}

func (c *fakeClient) ListOpenshiftServiceCatalogAPIServers(ctx context.Context, o metav1.ListOptions) (*operatorv1.ServiceCatalogAPIServerList, error) {
	return &operatorv1.ServiceCatalogAPIServerList{}, nil
}

func (c *fakeClient) ListOpenshiftServiceCatalogControllerManagers(ctx context.Context, o metav1.ListOptions) (*operatorv1.ServiceCatalogControllerManagerList, error) {
	return &operatorv1.ServiceCatalogControllerManagerList{}, nil
}

func (c *fakeClient) ListOpenshiftStorages(ctx context.Context, o metav1.ListOptions) (*operatorv1.StorageList, error) {
	return &operatorv1.StorageList{}, nil
}

func (c *fakeClient) ListOpenshiftProjects(ctx context.Context, o metav1.ListOptions) (*projectv1.ProjectList, error) {
	return &projectv1.ProjectList{}, nil
}

func (c *fakeClient) ListOpenshiftAppliedClusterResourceQuotas(ctx context.Context, namespace string, o metav1.ListOptions) (*quotav1.AppliedClusterResourceQuotaList, error) {
	return &quotav1.AppliedClusterResourceQuotaList{}, nil
}

func (c *fakeClient) ListOpenshiftClusterResourceQuotas(ctx context.Context, o metav1.ListOptions) (*quotav1.ClusterResourceQuotaList, error) {
	return &quotav1.ClusterResourceQuotaList{}, nil
}

func (c *fakeClient) ListOpenshiftRoutes(ctx context.Context, namespace string, o metav1.ListOptions) (*routev1.RouteList, error) {
	return &routev1.RouteList{}, nil
}

func (c *fakeClient) ListOpenshiftSampleConfigs(ctx context.Context, o metav1.ListOptions) (*samplesv1.ConfigList, error) {
	return &samplesv1.ConfigList{}, nil
}
