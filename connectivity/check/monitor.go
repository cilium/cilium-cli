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

package check

import (
	"context"
	"fmt"
	"time"

	"github.com/cilium/cilium-cli/defaults"
	"github.com/cilium/cilium-cli/internal/utils"
)

// Monitor keeps state about one Cilium Monitor instance.
type Monitor struct {
	name    string
	greeted bool
	done    chan struct{}
	stdin   *utils.CtrlCReader
	stdout  *utils.SyncBuffer
	err     error
	*Action // unnamed for logging convenience
}

// NewMonitor starts a new monitor command in background.
func (a *Action) NewMonitor(pod Pod, opts ...string) *Monitor {
	m := &Monitor{
		Action: a,
		name:   pod.Name(),
		done:   make(chan struct{}, 1),
		stdin:  utils.NewCtrlCReader(),
		stdout: utils.NewSyncBuffer(),
	}

	cmd := append([]string{"cilium", "monitor"}, opts...)

	go func() {
		defer close(m.done)
		m.err = pod.K8sClient.ExecInPodWithTTY(context.TODO(), pod.Pod.Namespace, pod.Pod.Name,
			"cilium-agent", cmd, m.stdin, m.stdout)
	}()

	return m
}

// Wait waits until the monitor has started executing or until the 5
// second timeout is exceeded.
func (m *Monitor) Wait() {
	// Return immediately if greeting has already been received
	if m.greeted {
		return
	}

	// Process initial monitor output
	done := m.stdout.ReadUntilLine([]byte("Press Ctrl-C to quit"))

	// Wait until Cilium monitor starts
	ticker := time.NewTicker(defaults.MonitorStartGracePeriod)
	defer ticker.Stop()
	for {
		select {
		case err := <-done:
			if err != nil {
				m.Debugf("Cilium monitor for %s did not start properly: %s", m.name, err)
				m.Stop()
			} else {
				m.Debugf("Cilium monitor for %s successfully started", m.name)
				m.greeted = true
			}
			return
		case <-m.done:
			m.Warnf("Cilium monitor for %s exited before starting.", m.name)
			return
		case <-ticker.C:
			m.Warnf("Cilium monitor for %s failed to start within %s.", m.name, defaults.MonitorStartGracePeriod)
			m.Stop()
			return
		}
	}
}

// Stop tells the monitor to stop execution.
func (m *Monitor) Stop() {
	m.stdin.Close()
}

// Output returns the monitor output after first waiting that the
// monitor has stopped.
func (m *Monitor) Output() (string, error) {
	// Wait for the monitor to have stopped.
	select {
	case <-m.done:
		return m.stdout.String(), m.err
	case <-time.After(5 * time.Second):
		return m.stdout.String(), fmt.Errorf("Timed out waiting for monitor to close, monitor is left running on the Cilium pod")
	}
}

// Name returns the name of the Cilium pod this monitor is running on.
func (m *Monitor) Name() string {
	return m.name
}
