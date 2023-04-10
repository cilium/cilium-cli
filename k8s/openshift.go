package k8s

import (
	oseAPI "github.com/openshift/api"

	//openshiftBuild "github.com/openshift/api/build"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"

	apiserverv1 "github.com/openshift/client-go/apiserver/clientset/versioned/typed/apiserver/v1"
	appsv1 "github.com/openshift/client-go/apps/clientset/versioned/typed/apps/v1"
	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	cloudnetworkv1 "github.com/openshift/client-go/cloudnetwork/clientset/versioned/typed/cloudnetwork/v1"
	configv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	imagev1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	machinev1 "github.com/openshift/client-go/machine/clientset/versioned/typed/machine/v1"
	networkv1 "github.com/openshift/client-go/network/clientset/versioned/typed/network/v1"
	operatorv1 "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	quotav1 "github.com/openshift/client-go/quota/clientset/versioned/typed/quota/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	samplesv1 "github.com/openshift/client-go/samples/clientset/versioned/typed/samples/v1"
)

// Not including
// authorization - concerned about security implications
// console - not necessary for sysdump
// helm - not necessary for sysdump
// imageregistry - not necessary for sysdump
// insights - not necessary for sysdump
// monitoring - not necessary for sysdump
// oauth - concerned about security implications
// operatorcontrolplane - doesn't do anything
// servicecertsigner - v1alpha - not sure if needed for sysdump
// sharedResource - v1alpha1 - not sure if needed
// template - not necessary for sysdump
// user - concerned about security implications
type OpenshiftClient struct {
	ApiserverClient    apiserverv1.ApiserverV1Interface
	AppsClient         appsv1.AppsV1Interface
	BuildClient        buildv1.BuildV1Interface
	ConfigClient       configv1.ConfigV1Interface
	CloudnetworkClient cloudnetworkv1.CloudV1Interface
	ImageClient        imagev1.ImageV1Interface
	MachineClient      machinev1.MachineV1Interface
	NetworkClient      networkv1.NetworkV1Interface
	OperatorClient     operatorv1.OperatorV1Interface
	ProjectClient      projectv1.ProjectV1Interface
	QuotaClient        quotav1.QuotaV1Interface
	RouteClient        routev1.RouteV1Interface
	SamplesClient      samplesv1.SamplesV1Interface
}

func NewOpenshiftClient(scheme *runtime.Scheme, config *rest.Config) (*OpenshiftClient, error) {

	_ = oseAPI.Install(scheme)

	apiserverClientset, err := apiserverv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	appsClientset, err := appsv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	buildClientset, err := buildv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	configClientset, err := configv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	cloudnetworkClientset, err := cloudnetworkv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	imageClientset, err := imagev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	machineClientset, err := machinev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	networkClientset, err := networkv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	operatorClientset, err := operatorv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	projectClientset, err := projectv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	quotaClientset, err := quotav1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	routeClientset, err := routev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	samplesClientset, err := samplesv1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &OpenshiftClient{
		ApiserverClient:    apiserverClientset,
		AppsClient:         appsClientset,
		BuildClient:        buildClientset,
		CloudnetworkClient: cloudnetworkClientset,
		ConfigClient:       configClientset,
		ImageClient:        imageClientset,
		MachineClient:      machineClientset,
		NetworkClient:      networkClientset,
		OperatorClient:     operatorClientset,
		ProjectClient:      projectClientset,
		QuotaClient:        quotaClientset,
		RouteClient:        routeClientset,
		SamplesClient:      samplesClientset,
	}, nil
}
