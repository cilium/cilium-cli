// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cli

import (
	"os"
	"time"

	"github.com/cilium/cilium-cli/multicast"
	"github.com/spf13/cobra"
)

func newCmdMulticast() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multicast",
		Short: "Manage multicast groups",
		Long:  ``,
	}
	cmd.AddCommand(
		newCmdMulticastList(),
		newCmdMulticastAdd(),
		newCmdMulticastDel(),
	)
	return cmd
}

func newCmdMulticastList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Shows the information about multicast groups",
		Long:  ``,
	}

	cmd.AddCommand(
		newCmdMulticastListGroup(),
		newCmdMulticastListSubscriber(),
	)
	return cmd

}

func newCmdMulticastListGroup() *cobra.Command {
	var params = multicast.Parameters{
		Writer: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "group",
		Short: "Shows list of multicast groups in every node",
		RunE: func(_ *cobra.Command, _ []string) error {
			params.CiliumNamespace = namespace
			mc := multicast.NewMulticast(k8sClient, params)
			err := mc.ListGroups()
			if err != nil {
				fatalf("Unable to list multicast groups: %s", err)
			}
			return nil
		},
	}
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", 1*time.Minute, "Maximum time to wait for result, default 1 minute")
	return cmd

}

func newCmdMulticastListSubscriber() *cobra.Command {
	var params = multicast.Parameters{
		Writer: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "subscriber",
		Short: "Shows list of subscribers belonging to the specified multicast group",
		RunE: func(_ *cobra.Command, _ []string) error {
			params.CiliumNamespace = namespace
			mc := multicast.NewMulticast(k8sClient, params)
			err := mc.ListSubscribers()
			if err != nil {
				fatalf("Unable to list subscribers of the multicast group: %s", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&params.MulticastGroupIP, "group-ip", "g", "", "Multicast group IP address")
	cmd.Flags().BoolVar(&params.All, "all", false, "Show all subscribers")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", 1*time.Minute, "Maximum time to wait for result, default 1 minute")
	return cmd

}

func newCmdMulticastAdd() *cobra.Command {
	var params = multicast.Parameters{
		Writer: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add all nodes to the specified multicast group as subscribers in every cilium-agent",
		RunE: func(_ *cobra.Command, _ []string) error {
			params.CiliumNamespace = namespace
			mc := multicast.NewMulticast(k8sClient, params)
			err := mc.AddAllNodes()
			if err != nil {
				fatalf("Unable to add all nodes: %s", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&params.MulticastGroupIP, "group-ip", "g", "", "Multicast group IP address")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", 1*time.Minute, "Maximum time to wait for result, default 1 minute")
	return cmd
}

func newCmdMulticastDel() *cobra.Command {
	var params = multicast.Parameters{
		Writer: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "del",
		Short: "Delete the specified multicast group in every cilium-agent",
		RunE: func(_ *cobra.Command, _ []string) error {
			params.CiliumNamespace = namespace
			mc := multicast.NewMulticast(k8sClient, params)
			err := mc.DelAllNodes()
			if err != nil {
				fatalf("Unable to delete all nodes: %s", err)
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&params.MulticastGroupIP, "group-ip", "g", "", "Multicast group IP address")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", 1*time.Minute, "Maximum time to wait for result, default 1 minute")
	return cmd
}
