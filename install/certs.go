// SPDX-License-Identifier: Apache-2.0
// Copyright 2020 Authors of Cilium

package install

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/k8s"
)

func (k *K8sUninstaller) uninstallCerts(ctx context.Context) (err error) {
	if err2 := k.client.DeleteSecret(ctx, k.params.Namespace, defaults.CASecretName, metav1.DeleteOptions{}); err2 != nil {
		err2 = fmt.Errorf("unable to delete CA secret %s/%s: %w", k.params.Namespace, defaults.CASecretName, err2)
		if err == nil {
			err = err2
		}
	}

	return err
}

func (k *K8sInstaller) installCerts(ctx context.Context) error {
	if k.params.InheritCA != "" {
		caCluster, err := k8s.NewClient(k.params.InheritCA, "")
		if err != nil {
			return fmt.Errorf("unable to create Kubernetes client to derive CA from: %w", err)
		}

		s, err := caCluster.GetSecret(ctx, k.params.Namespace, defaults.CASecretName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("secret %s not found to derive CA from: %w", defaults.CASecretName, err)
		}

		newSecret := k8s.NewSecret(defaults.CASecretName, k.params.Namespace, s.Data)
		_, err = k.client.CreateSecret(ctx, k.params.Namespace, newSecret, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to create secret to store CA: %w", err)
		}
		k.pushRollbackStep(func(ctx context.Context) {
			if err := k.client.DeleteSecret(ctx, k.params.Namespace, defaults.CASecretName, metav1.DeleteOptions{}); err != nil {
				k.Log("Cannot delete %s Secret: %s", defaults.CASecretName, err)
			}
		})
	}

	caSecret, created, err := k.certManager.GetOrCreateCASecret(ctx, defaults.CASecretName, true)
	if err != nil {
		k.Log("❌ Unable to get or create the Cilium CA Secret: %s", err)
		return err
	}

	if caSecret != nil {
		err = k.certManager.LoadCAFromK8s(ctx, caSecret)
		if err != nil {
			k.pushRollbackStep(func(ctx context.Context) {
				if err := k.client.DeleteSecret(ctx, k.params.Namespace, caSecret.Name, metav1.DeleteOptions{}); err != nil {
					k.Log("Cannot delete %s Secret: %s", caSecret.Name, err)
				}
			})
			return err
		}
		if created {
			k.Log("🔑 Created CA in secret %s", caSecret.Name)
		} else {
			k.Log("🔑 Found CA in secret %s", caSecret.Name)
		}
	}

	return nil
}
