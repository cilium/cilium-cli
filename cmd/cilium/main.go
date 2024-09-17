// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"fmt"
	"os"

	"github.com/cilium/cilium/cilium-cli/cli"
	_ "github.com/cilium/cilium/cilium-cli/logging" // necessary to disable unwanted cfssl log messages
)

func main() {
	if err := cli.NewDefaultCiliumCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
