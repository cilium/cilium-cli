// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	ciliumiov2 "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned/typed/cilium.io/v2"
	gentype "k8s.io/client-go/gentype"
)

// fakeCiliumClusterwideEnvoyConfigs implements CiliumClusterwideEnvoyConfigInterface
type fakeCiliumClusterwideEnvoyConfigs struct {
	*gentype.FakeClientWithList[*v2.CiliumClusterwideEnvoyConfig, *v2.CiliumClusterwideEnvoyConfigList]
	Fake *FakeCiliumV2
}

func newFakeCiliumClusterwideEnvoyConfigs(fake *FakeCiliumV2) ciliumiov2.CiliumClusterwideEnvoyConfigInterface {
	return &fakeCiliumClusterwideEnvoyConfigs{
		gentype.NewFakeClientWithList[*v2.CiliumClusterwideEnvoyConfig, *v2.CiliumClusterwideEnvoyConfigList](
			fake.Fake,
			"",
			v2.SchemeGroupVersion.WithResource("ciliumclusterwideenvoyconfigs"),
			v2.SchemeGroupVersion.WithKind("CiliumClusterwideEnvoyConfig"),
			func() *v2.CiliumClusterwideEnvoyConfig { return &v2.CiliumClusterwideEnvoyConfig{} },
			func() *v2.CiliumClusterwideEnvoyConfigList { return &v2.CiliumClusterwideEnvoyConfigList{} },
			func(dst, src *v2.CiliumClusterwideEnvoyConfigList) { dst.ListMeta = src.ListMeta },
			func(list *v2.CiliumClusterwideEnvoyConfigList) []*v2.CiliumClusterwideEnvoyConfig {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v2.CiliumClusterwideEnvoyConfigList, items []*v2.CiliumClusterwideEnvoyConfig) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
