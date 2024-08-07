// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package check

import appsv1 "k8s.io/api/apps/v1"

const (
	// SocatServerPort is the port on which the socat server listens.
	netshootSocatServerDaemonsetName  = "netshoot-socat-server"
	netshootSocatClientDeploymentName = "netshoot-socat-client"
)

func NewSocatServerDaemonSet(params Parameters) *appsv1.DaemonSet {
	ds := newDaemonSet(daemonSetParameters{
		Name:    netshootSocatServerDaemonsetName,
		Kind:    netshootSocatServerDaemonsetName,
		Image:   params.NetshootImage,
		Command: []string{"/bin/bash", "-c", "sleep 10000000"},
	})
	return ds
}

func NewSocatClientDeployment(params Parameters) *appsv1.Deployment {
	dep := newDeployment(deploymentParameters{
		Name:     netshootSocatClientDeploymentName,
		Kind:     netshootSocatClientDeploymentName,
		Image:    params.NetshootImage,
		Replicas: 1,
		Command:  []string{"/bin/bash", "-c", "sleep 10000000"},
	})
	return dep
}
