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

package cmd

import (
	"context"

	"github.com/cilium/cilium-cli/exec"

	"github.com/spf13/cobra"
)

const (
	defaultCommand = "/bin/bash"
)

func newCmdExec() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "exec (NODE_NAME|CILIUM_POD_NAME|NAMESPACE/POD_NAME) [-- COMMAND [args...]]",
		DisableFlagsInUseLine: true,
		Short:                 "Connectivity troubleshooting",
		Long:                  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle the case when there's a single argument (i.e. the target).
			if len(args) == 1 && cmd.ArgsLenAtDash() == -1 {
				return execInTarget(args[0], defaultCommand)
			}
			// Handle the case where there are at least two arguments (i.e. the target, '--', and the command(s) to run inside the Cilium agent).
			if len(args) >= 2 && cmd.ArgsLenAtDash() == 1 {
				return execInTarget(args[0], args[cmd.ArgsLenAtDash():]...)
			}
			// Everything else should be an invalid.
			return cmd.Usage()
		},
	}
	return cmd
}

func execInTarget(target string, command ...string) error {
	i, err := exec.NewCiliumExecImplementation(k8sClient, exec.Parameters{
		Target:  target,
		Command: command,
	})
	if err != nil {
		return err
	}
	return i.ExecInTarget(context.Background())
}
