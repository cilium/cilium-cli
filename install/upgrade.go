// Copyright 2021 Authors of Cilium
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

package install

import (
	"context"
	"fmt"
	"strings"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/internal/utils"
	"github.com/cilium/cilium-cli/status"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	configNameEnableHubble = "enable-hubble"
)

func (k *K8sInstaller) Upgrade(ctx context.Context) error {
	daemonSet, err := k.client.GetDaemonSet(ctx, k.params.Namespace, defaults.AgentDaemonSetName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to retrieve DaemonSet of cilium-agent: %s", err)
	}

	deployment, err := k.client.GetDeployment(ctx, k.params.Namespace, defaults.OperatorDeploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to retrieve Deployment of cilium-operator: %s", err)
	}

	cm, err := k.client.GetConfigMap(ctx, k.params.Namespace, defaults.ConfigMapName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("unable to retrieve ConfigMap %q: %w", defaults.ConfigMapName, err)
	}

	if cm.Data == nil {
		return fmt.Errorf("ConfigMap %q does not contain any configuration", defaults.ConfigMapName)
	}

	var patched int

	if deployment.Spec.Template.Spec.Containers[0].Image == k.fqOperatorImage() {
		k.Log("✅ Cilium-operator is already up to date")
	} else {
		k.Log("🚀 Upgrading cilium-operator to version %s...", k.fqOperatorImage())
		patch := []byte(`{"spec":{"template":{"spec":{"containers":[{"name": "cilium-operator", "image":"` + k.fqOperatorImage() + `"}]}}}}`)

		_, err = k.client.PatchDeployment(ctx, k.params.Namespace, defaults.OperatorDeploymentName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("unable to patch Deployment %s with patch %q: %w", defaults.OperatorDeploymentName, patch, err)
		}

		patched++
	}

	agentImage := k.fqAgentImage()
	var containerPatches []string
	for _, c := range daemonSet.Spec.Template.Spec.Containers {
		if c.Image != agentImage {
			containerPatches = append(containerPatches, `{"name":"`+c.Name+`", "image":"`+agentImage+`"}`)
		}
	}
	var initContainerPatches []string
	for _, c := range daemonSet.Spec.Template.Spec.InitContainers {
		if c.Image != agentImage {
			initContainerPatches = append(initContainerPatches, `{"name":"`+c.Name+`", "image":"`+agentImage+`"}`)
		}
	}

	if len(containerPatches) == 0 && len(initContainerPatches) == 0 {
		k.Log("✅ Cilium is already up to date")
	} else {
		k.Log("🚀 Upgrading cilium to version %s...", k.fqAgentImage())

		patch := []byte(`{"spec":{"template":{"spec":{"containers":[` + strings.Join(containerPatches, ",") + `], "initContainers":[` + strings.Join(initContainerPatches, ",") + `]}}}}`)
		_, err = k.client.PatchDaemonSet(ctx, k.params.Namespace, defaults.AgentDaemonSetName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("unable to patch DaemonSet %s with patch %q: %w", defaults.AgentDaemonSetName, patch, err)
		}

		patched++
	}

	hubbleDeployment, err := k.client.GetDeployment(ctx, k.params.Namespace, defaults.RelayDeploymentName, metav1.GetOptions{})
	enableHubble, ok := cm.Data[configNameEnableHubble]
	if err != nil {
		k.Log("❌  Unable to retrieve Deployment of hubble-relay, unable to upgrade")
	} else if !ok || strings.ToLower(enableHubble) != "true" {
		k.Log("❌  Hubble is not enabled in ConfigMap, %q is not set to true,unable to upgrade", configNameEnableHubble)
	} else if hubbleDeployment.Spec.Template.Spec.Containers[0].Image == k.RelayImage() {
		k.Log("✅ Hubble-relay is already up to date")
	} else {
		k.Log("🚀 Upgrading hubble-relay to version %s...", k.RelayImage())
		patch := []byte(`{"spec":{"template":{"spec":{"containers":[{"name": "hubble-relay", "image":"` + k.RelayImage() + `"}]}}}}`)

		_, err = k.client.PatchDeployment(ctx, k.params.Namespace, defaults.RelayDeploymentName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			return fmt.Errorf("unable to patch Deployment %s with patch %q: %w", defaults.RelayDeploymentName, patch, err)
		}

		patched++
	}

	if patched == 1 && k.params.Wait {
		k.Log("⌛ Waiting for Hubble-relay to be upgraded...")
		collector, err := status.NewK8sStatusCollector(ctx, k.client, status.K8sStatusParameters{
			Namespace:       k.params.Namespace,
			Wait:            true,
			WaitDuration:    k.params.WaitDuration,
			WarningFreePods: []string{defaults.RelayDeploymentName},
		})
		if err != nil {
			return err
		}

		s, err := collector.Status(ctx)
		if err != nil {
			if s != nil {
				fmt.Println(s.Format())
			}
			return err
		}
	}

	if patched == 2 && k.params.Wait {
		k.Log("⌛ Waiting for Cilium to be upgraded...")
		collector, err := status.NewK8sStatusCollector(ctx, k.client, status.K8sStatusParameters{
			Namespace:       k.params.Namespace,
			Wait:            true,
			WaitDuration:    k.params.WaitDuration,
			WarningFreePods: []string{defaults.AgentDaemonSetName, defaults.OperatorDeploymentName},
		})
		if err != nil {
			return err
		}

		s, err := collector.Status(ctx)
		if err != nil {
			if s != nil {
				fmt.Println(s.Format())
			}
			return err
		}
	}

	if patched >= 3 && k.params.Wait {
		k.Log("⌛ Waiting for Cilium & Hubble-relay to be upgraded...")
		collector, err := status.NewK8sStatusCollector(ctx, k.client, status.K8sStatusParameters{
			Namespace:       k.params.Namespace,
			Wait:            true,
			WaitDuration:    k.params.WaitDuration,
			WarningFreePods: []string{defaults.AgentDaemonSetName, defaults.OperatorDeploymentName, defaults.RelayDeploymentName},
		})
		if err != nil {
			return err
		}

		s, err := collector.Status(ctx)
		if err != nil {
			if s != nil {
				fmt.Println(s.Format())
			}
			return err
		}
	}

	return nil
}

func (k *K8sInstaller) RelayImage() string {
	return utils.BuildImagePath(k.params.RelayImage, defaults.RelayImage, k.params.Version, defaults.Version)
}
