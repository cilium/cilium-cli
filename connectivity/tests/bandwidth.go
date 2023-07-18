// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package tests

import (
	"context"
	"fmt"

	"github.com/cilium/cilium-cli/connectivity/check"
)

// BandWidth Manager
func BandWidthManager(n string) check.Scenario {
	return &bandWidthManager{
		name: n,
	}
}

// bandWidthManager implements a Scenario.
type bandWidthManager struct {
	name string
}

func (b *bandWidthManager) Name() string {
	tn := "bandwidth-manager"
	if b.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, b.name)
}

func (b *bandWidthManager) Run(ctx context.Context, t *check.Test) {
	for _, c := range t.Context().PerfClientPods() {
		c := c
		for _, server := range t.Context().PerfServerPod() {
			action := t.NewAction(b, "bandwidth", &c, server, check.IPFamilyV4)
			action.Run(func(a *check.Action) {})
		}
	}
}
