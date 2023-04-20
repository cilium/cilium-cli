// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/cilium/cilium-cli/connectivity/check"
	"github.com/cilium/cilium-cli/internal/cli/cmd"
	"github.com/cilium/cilium-cli/sysdump"
)

func NewDefaultCiliumCommand() *cobra.Command {
	return NewCiliumCommand(&NopHooks{})
}

func NewCiliumCommand(hooks Hooks) *cobra.Command {
	return cmd.NewCiliumCommand(hooks)
}

type (
	Hooks                 = cmd.Hooks
	ConnectivityTestHooks = cmd.ConnectivityTestHooks
	SysdumpHooks          = cmd.SysdumpHooks
)

type NopHooks struct{}

var _ Hooks = &NopHooks{}

func (*NopHooks) AddSysdumpFlags(flags *pflag.FlagSet)                  {}
func (*NopHooks) AddSysdumpTasks(*sysdump.Collector) error              { return nil }
func (*NopHooks) AddConnectivityTestFlags(flags *pflag.FlagSet)         {}
func (*NopHooks) AddConnectivityTests(ct *check.ConnectivityTest) error { return nil }
