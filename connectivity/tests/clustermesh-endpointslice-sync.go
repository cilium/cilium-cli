// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package tests

import (
	"context"
	"fmt"

	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/utils/features"
)

func ClusterMeshEndpointSliceSync() check.Scenario {
	return &clusterMeshEndpointSliceSync{}
}

type clusterMeshEndpointSliceSync struct{}

func (s *clusterMeshEndpointSliceSync) Name() string {
	return "clustermesh-endpointslice-sync"
}

func (s *clusterMeshEndpointSliceSync) Run(ctx context.Context, t *check.Test) {
	ct := t.Context()
	client := ct.RandomClientPod()

	service, ok := ct.EchoServicesAll()[check.EchoOtherNodeDeploymentHeadlessServiceName]
	if !ok {
		t.Fatal("Cannot get echo-other-node-headless service")
	}

	t.ForEachIPFamily(func(ipFam features.IPFamily) {
		t.NewAction(s, fmt.Sprintf("dig-%s", ipFam), client, service, ipFam).Run(func(a *check.Action) {
			a.ExecInPod(ctx, ct.DigCommandService(service, ipFam))
		})
	})
}
