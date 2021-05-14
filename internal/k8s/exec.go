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

package k8s

import (
	"context"
	"fmt"
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/util/interrupt"
	"k8s.io/kubectl/pkg/util/term"
)

type ExecParameters struct {
	Namespace       string
	Pod             string
	Container       string
	Command         []string
	Interactive     bool
	In              io.Reader
	Out             io.Writer
	Err             io.Writer
	InterruptParent *interrupt.Handler
}

func (c *Client) execInPod(ctx context.Context, p *ExecParameters) error {
	t, err := setupTTY(p)
	if err != nil {
		return err
	}
	var sizeQueue remotecommand.TerminalSizeQueue
	if t.Raw {
		sizeQueue = t.MonitorSize(t.GetSize())
	}

	fn := func() error {
		req := c.Clientset.CoreV1().RESTClient().Post().Resource("pods").Name(p.Pod).Namespace(p.Namespace).SubResource("exec")

		scheme := runtime.NewScheme()
		if err := corev1.AddToScheme(scheme); err != nil {
			return fmt.Errorf("error adding to scheme: %w", err)
		}

		parameterCodec := runtime.NewParameterCodec(scheme)

		req.VersionedParams(&corev1.PodExecOptions{
			Command:   p.Command,
			Container: p.Container,
			Stdin:     p.Interactive,
			Stdout:    true,
			Stderr:    true,
			TTY:       p.Interactive,
		}, parameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(c.Config, "POST", req.URL())
		if err != nil {
			return fmt.Errorf("error while creating executor: %w", err)
		}
		return exec.Stream(remotecommand.StreamOptions{
			Stdin:             p.In,
			Stdout:            p.Out,
			Stderr:            p.Err,
			Tty:               p.Interactive,
			TerminalSizeQueue: sizeQueue,
		})
	}
	if err := t.Safe(fn); err != nil {
		return err
	}
	return nil
}

func setupTTY(p *ExecParameters) (*term.TTY, error) {
	t := &term.TTY{
		Parent: p.InterruptParent,
		Out:    p.Out,
	}
	if !p.Interactive {
		return t, nil
	}
	t.Raw = true
	t.Out = os.Stdout
	t.In = os.Stdin
	if !t.IsTerminalIn() {
		return nil, fmt.Errorf("unable to create a TTY as input is not a terminal")
	}
	p.In = t.In
	p.Out = t.Out
	p.Err = nil // Unset p.Err because both stdout and stderr go over p.Out when tty is true.
	return t, nil
}
