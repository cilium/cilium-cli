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

package hubble

import (
	"context"

	"github.com/cilium/cilium-cli/defaults"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var hubbleCLIReplicas = int32(1)

func (k *K8sHubble) generateHubbleCLIDeployment() *appsv1.Deployment {
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   defaults.HubbleCLIDeploymentName,
			Labels: defaults.HubbleCLIDeploymentLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &hubbleCLIReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: defaults.HubbleCLIDeploymentLabels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   defaults.HubbleCLIDeploymentName,
					Labels: defaults.HubbleCLIDeploymentLabels,
				},
				Spec: corev1.PodSpec{
					RestartPolicy:      corev1.RestartPolicyAlways,
					ServiceAccountName: defaults.RelayServiceAccountName,
					Containers: []corev1.Container{
						{
							Name: "hubble-cli",
							// NOTE: just run forever
							Command: []string{"tail"},
							Args: []string{
								"-f",
								"/dev/null",
							},
							Image: "quay.io/cilium/hubble:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "HUBBLE_SERVER",
									Value: "$(HUBBLE_RELAY_SERVICE_HOST):$(HUBBLE_RELAY_SERVICE_PORT)",
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "hubble-sock-dir",
									MountPath: "/var/run/cilium",
									ReadOnly:  true,
								},
								{
									Name:      "tls",
									MountPath: "/var/lib/hubble-relay/tls",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "hubble-sock-dir",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/run/cilium",
									Type: &hostPathDirectoryOrCreate,
								},
							},
						},
						{
							Name: "tls",
							VolumeSource: corev1.VolumeSource{
								Projected: &corev1.ProjectedVolumeSource{
									Sources: []corev1.VolumeProjection{
										{
											Secret: &corev1.SecretProjection{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: defaults.RelayClientSecretName,
												},
												Items: []corev1.KeyToPath{
													{
														Key:  corev1.TLSCertKey,
														Path: "client.crt",
													},
													{
														Key:  corev1.TLSPrivateKeyKey,
														Path: "client.key",
													},
												},
											},
										},
										{
											Secret: &corev1.SecretProjection{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: defaults.CASecretName,
												},
												Items: []corev1.KeyToPath{
													{
														Key:  defaults.CASecretCertName,
														Path: "hubble-server-ca.crt",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return d
}

func (k *K8sHubble) disableHubbleCLI(ctx context.Context) error {
	k.Log("🔥 Deleting Hubble CLI...")
	k.client.DeleteDeployment(ctx, k.params.Namespace, defaults.HubbleCLIDeploymentName, metav1.DeleteOptions{})
	return nil
}

func (k *K8sHubble) enableHubbleCLI(ctx context.Context) error {
	_, err := k.client.GetDeployment(ctx, k.params.Namespace, defaults.HubbleCLIDeploymentName, metav1.GetOptions{})
	if err == nil {
		k.Log("✅ Hubble CLI is already deployed")
		return nil
	}

	// We need relay to be deployed since we're using its ServiceAccountName
	// and Secret.
	_, err = k.client.GetDeployment(ctx, k.params.Namespace, defaults.RelayDeploymentName, metav1.GetOptions{})
	if err != nil {
		k.Log("❌ Relay is not deployed")
		return err
	}

	k.Log("✨ Deploying Hubble CLI...")
	_, err = k.client.CreateDeployment(ctx, k.params.Namespace, k.generateHubbleCLIDeployment(), metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
