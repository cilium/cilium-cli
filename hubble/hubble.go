// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package hubble

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/internal/certs"
	"github.com/cilium/cilium-cli/internal/helm"
	"github.com/cilium/cilium-cli/k8s"
	"helm.sh/helm/v3/pkg/cli/values"
)

type K8sHubble struct {
	params      Parameters
	certManager *certs.CertManager
}

type Parameters struct {
	Namespace        string
	Relay            bool
	RelayImage       string
	RelayVersion     string
	RelayServiceType string
	PortForward      int
	CreateCA         bool
	UI               bool
	UIImage          string
	UIBackendImage   string
	UIVersion        string
	UIPortForward    int
	Writer           io.Writer
	Context          string // Only for 'kubectl' pass-through commands
	Wait             bool
	WaitDuration     time.Duration

	// K8sVersion is the Kubernetes version that will be used to generate the
	// kubernetes manifests. If the auto-detection fails, this flag can be used
	// as a workaround.
	K8sVersion string
	// HelmChartDirectory points to the location of a helm chart directory.
	// Useful to test from upstream where a helm release is not available yet.
	HelmChartDirectory string

	// HelmOpts are all the options the user used to pass into the Cilium cli
	// template.
	HelmOpts values.Options

	// HelmGenValuesFile points to the file that will store the generated helm
	// options.
	HelmGenValuesFile string

	// HelmValuesSecretName is the name of the secret where helm values will be
	// stored.
	HelmValuesSecretName string

	// RedactHelmCertKeys does not print helm certificate keys into the terminal.
	RedactHelmCertKeys bool

	// UIOpenBrowser will automatically open browser if true
	UIOpenBrowser bool
}

func (p *Parameters) Log(format string, a ...interface{}) {
	fmt.Fprintf(p.Writer, format+"\n", a...)
}

func (k *K8sHubble) Log(format string, a ...interface{}) {
	if k.params.RedactHelmCertKeys {
		formattedString := fmt.Sprintf(format+"\n", a...)
		for _, certKey := range []string{
			certs.EncodeCertBytes(k.certManager.CAKeyBytes()),
		} {
			if certKey != "" {
				formattedString = strings.ReplaceAll(formattedString, certKey, "[--- REDACTED WHEN PRINTING TO TERMINAL (USE --redact-helm-certificate-keys=false TO PRINT) ---]")
			}
		}
		fmt.Fprint(k.params.Writer, formattedString)
		return
	}
	fmt.Fprintf(k.params.Writer, format+"\n", a...)
}

func EnableWithHelm(ctx context.Context, k8sClient *k8s.Client, params Parameters) error {
	options := values.Options{
		Values: []string{
			fmt.Sprintf("hubble.relay.enabled=%t", params.Relay),
			fmt.Sprintf("hubble.ui.enabled=%t", params.UI),
		},
	}
	vals, err := helm.MergeVals(options, nil, nil, nil)
	if err != nil {
		return err
	}
	upgradeParams := helm.UpgradeParameters{
		Namespace:   params.Namespace,
		Name:        defaults.HelmReleaseName,
		Values:      vals,
		ResetValues: false,
		ReuseValues: true,
	}
	_, err = helm.Upgrade(ctx, k8sClient.HelmActionConfig, upgradeParams)
	return err
}

func DisableWithHelm(ctx context.Context, k8sClient *k8s.Client, params Parameters) error {
	options := values.Options{
		Values: []string{"hubble.relay.enabled=false", "hubble.ui.enabled=false"},
	}
	vals, err := helm.MergeVals(options, nil, nil, nil)
	if err != nil {
		return err
	}
	upgradeParams := helm.UpgradeParameters{
		Namespace:   params.Namespace,
		Name:        defaults.HelmReleaseName,
		Values:      vals,
		ResetValues: false,
		ReuseValues: true,
	}
	_, err = helm.Upgrade(ctx, k8sClient.HelmActionConfig, upgradeParams)
	return err
}
