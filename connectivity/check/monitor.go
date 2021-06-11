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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

// CtrlCReader implements a simple Reader/Closer that returns Ctrl-C and EOF
// on Read() after it has been closed, and nothing before it.
type CtrlCReader struct {
	once sync.Once
	done chan struct{}
}

// NewCtrlCReader returns a new CtrlCReader instance
func NewCtrlCReader() *CtrlCReader {
	return &CtrlCReader{
		done: make(chan struct{}),
	}
}

// Read implements io.Reader.
func (cc *CtrlCReader) Read(p []byte) (n int, err error) {
	select {
	case <-cc.done:
		if len(p) > 0 {
			p[0] = byte(3) // Ctrl-C
			return 1, io.EOF
		}
	default:
	}
	return 0, nil
}

// Close implements io.Closer. Note that we do not return an error on
// second close, not do we wait for the close to have any effect.
func (cc *CtrlCReader) Close() error {
	cc.once.Do(func() {
		close(cc.done)
	})
	return nil
}

// Buffer is a concurrency safe buffered Reader/Writer.  This is
// needed as we are reading the initial lines concurrently with the
// monitor execution.
type Buffer struct {
	sync.Mutex
	b bytes.Buffer
}

// Read implements io.Reader.
func (b *Buffer) Read(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	return b.b.Read(p)
}

// Write implements io.Writer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	return b.b.Write(p)
}

// String implements fmt.Stringer.
func (b *Buffer) String() string {
	b.Lock()
	defer b.Unlock()
	return b.b.String()
}

// Monitor keeps state about one Cilium Monitor instance.
type Monitor struct {
	name    string
	greeted bool
	done    chan struct{}
	stdin   *CtrlCReader
	stdout  *Buffer
	err     error
	*Action // unnamed for logging convenience
}

// NewMonitor starts a new monitor command in background.
func (a *Action) NewMonitor(pod Pod, opts ...string) *Monitor {
	m := &Monitor{
		Action: a,
		name:   pod.Name(),
		done:   make(chan struct{}, 1),
		stdin:  NewCtrlCReader(),
		stdout: &Buffer{},
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

	// Wait until Cilium monitor responds with "Press Ctrl-C to quit".
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.done:
			m.Warnf("Cilium monitor exited before starting.")
			return
		case <-ticker.C:
			m.Warnf("Cilium monitor failed to start within 5 seconds.")
			return
		default:
			reader := bufio.NewReader(m.stdout)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					break
				}
				if err == nil && string(line) == "Press Ctrl-C to quit" {
					// Greeting received, do not wait again in future
					m.greeted = true
					m.Debugf("Cilium monitor for %s successfully started", m.name)
					return
				}
			}
			time.Sleep(50 * time.Millisecond)
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
