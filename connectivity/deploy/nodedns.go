// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package deploy

import (
	"context"
	_ "embed"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/cilium/cilium-cli/connectivity/builder/manifests/template"
	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/k8s"
)

var (
	//go:embed manifests/node-local-dns-service.yaml
	nodeLocalDNSService string

	//go:embed manifests/node-local-dns-cm.yaml
	nodeLocalDNSConfigMap string

	//go:embed manifests/node-local-dns-ds.yaml
	nodeLocalDNSDaemonSet string
)

func NodeDNS(ctx context.Context, t *check.Test, ct *check.ConnectivityTest) error {
	var (
		ns      = "kube-system"
		name    = "node-local-dns"
		svcName = "kube-dns-upstream"
		clients = ct.Clients()
	)

	client := clients[0]

	_, err := client.GetServiceAccount(ctx, ns, name, metav1.GetOptions{})
	if err != nil {
		_, err = client.CreateServiceAccount(ctx, ns, k8s.NewServiceAccount(name), metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to create service account %s: %w", name, err)
		}
	}

	t.WithFinalizer(func(_ context.Context) error {
		// Use a detached context to make sure this call is not affected by
		// context cancellation. This deletion needs to happen event when the
		// user interrupted the program.
		if err := client.DeleteServiceAccount(context.TODO(), ns, name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("unable to delete service account %s: %w", name, err)
		}
		return nil
	})

	_, err = client.GetService(ctx, ns, svcName, metav1.GetOptions{})
	if err != nil {
		sv := &corev1.Service{}
		err = yaml.Unmarshal([]byte(nodeLocalDNSService), &sv)
		if err != nil {
			return err
		}

		_, err = client.CreateService(ctx, ns, sv, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to create service %s: %w", name, err)
		}
	}

	t.WithFinalizer(func(_ context.Context) error {
		if err := client.DeleteService(context.TODO(), ns, svcName, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("unable to delete service %s: %w", svcName, err)
		}
		return nil
	})

	_, err = client.GetConfigMap(ctx, ns, name, metav1.GetOptions{})
	if err != nil {
		cm := &corev1.ConfigMap{}
		err = yaml.Unmarshal([]byte(nodeLocalDNSConfigMap), &cm)
		if err != nil {
			return err
		}

		_, err = client.CreateConfigMap(ctx, ns, cm, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to create configmap %s: %s", name, err)
		}
	}

	t.WithFinalizer(func(_ context.Context) error {
		if err := client.DeleteConfigMap(context.TODO(), ns, name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("unable to delete configMap %s: %w", name, err)
		}
		return nil
	})

	_, err = client.GetDaemonSet(ctx, ns, name, metav1.GetOptions{})
	if err != nil {
		kubeDNSService, err := client.GetService(ctx, ns, "kube-dns", metav1.GetOptions{})
		if err != nil {
			return err
		}

		dsYaml, err := template.Render(nodeLocalDNSDaemonSet, map[string]interface{}{
			"PILLAR_DNS_SERVER": kubeDNSService.Spec.ClusterIP,
		})
		if err != nil {
			return err
		}

		ds := &appsv1.DaemonSet{}
		err = yaml.Unmarshal([]byte(dsYaml), &ds)
		if err != nil {
			return err
		}
		_, err = client.CreateDaemonSet(ctx, ns, ds, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to create daemonset %s: %w", name, err)
		}
	}

	t.WithFinalizer(func(_ context.Context) error {
		if err := client.DeleteDaemonSet(context.TODO(), ns, name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("unable to delete DaemonSet %s: %w", name, err)
		}
		return nil
	})

	if err := check.WaitForDaemonSet(ctx, ct, client, ns, name); err != nil {
		return err
	}

	return nil
}
