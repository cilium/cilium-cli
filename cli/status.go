// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/status"
)

func newCmdStatus() *cobra.Command {
	var params = status.K8sStatusParameters{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Display status",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, _ []string) {
			pod, _ := cmd.Flags().GetString("pod")
			output, _ := cmd.Flags().GetString("output")
			if pod != "" && output != "text" {
				fatalf("specifying a --pod requires output option -o text")
			}
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			params.Namespace = namespace

			collector, err := status.NewK8sStatusCollector(k8sClient, params)
			if err != nil {
				return err
			}

			s, err := collector.Status(context.Background())
			if err != nil {
				// Report the most recent status even if an error occurred.
				fmt.Fprint(os.Stderr, s.Format())
				fatalf("Unable to determine status:  %s", err)
			}
			switch params.Output {
			case status.OutputJSON:
				{
					jsonStatus, err := json.MarshalIndent(s, "", " ")
					if err != nil {
						// Report the most recent status even if an error occurred.
						fmt.Fprint(os.Stderr, s.Format())
						fatalf("Unable to marshal status to JSON:  %s", err)
					}
					fmt.Println(string(jsonStatus))
				}
			case status.OutputText:
				if params.Pod == "" {
					for key, value := range s.CiliumStatusInText {
						// Randomly pick one of the agent to print out the status when the pods is empty
						fmt.Printf("cilium status --verbose from %s\n", key)
						fmt.Println(value)
						break
					}

				} else {
					if _, value := s.CiliumStatusInText[params.Pod]; value {
						fmt.Printf("cilium status --verbose from %s\n", params.Pod)
						fmt.Println(s.CiliumStatusInText[params.Pod])
					} else {
						fatalf("the pod name for cilium status -o text doesn't exist")
					}
				}

			default:
				fmt.Print(s.Format())
			}

			if err == nil && len(s.CollectionErrors) > 0 {
				errs := make([]string, 0, len(s.CollectionErrors))
				for _, e := range s.CollectionErrors {
					errs = append(errs, e.Error())
				}
				err = fmt.Errorf("status check failed: [%s]", strings.Join(errs, ", "))
			}
			return err
		},
	}
	cmd.Flags().BoolVar(&params.Wait, "wait", false, "Wait for status to report success (no errors and warnings)")
	cmd.Flags().DurationVar(&params.WaitDuration, "wait-duration", defaults.StatusWaitDuration, "Maximum time to wait for status")
	cmd.Flags().BoolVar(&params.IgnoreWarnings, "ignore-warnings", false, "Ignore warnings when waiting for status to report success")
	cmd.Flags().IntVar(&params.WorkerCount,
		"worker-count", status.DefaultWorkerCount,
		"The number of workers to use")
	cmd.Flags().StringVarP(&params.Output, "output", "o", status.OutputSummary, "Output format. One of: json, text, summary")
	cmd.Flags().StringVarP(&params.Pod, "pod", "p", "", "Pod name when use cilium status -o text")
	cmd.Flags().BoolVar(&params.Interactive, "interactive", true, "Refresh the status summary output after each retry when --wait flag is specified")

	return cmd
}
