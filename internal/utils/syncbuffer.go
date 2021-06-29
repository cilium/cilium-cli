// Copyright 2020-2021 Authors of Cilium
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

package utils

import (
	"bytes"
	"io"
	"sync"
)

// SyncBuffer is a concurrency safe buffered Reader/Writer.  This is
// needed as we are reading the initial lines concurrently with the
// monitor execution.
type SyncBuffer struct {
	sync.Mutex
	readReady *sync.Cond
	buffer    bytes.Buffer
	closed    bool
}

// NewSyncBuffer returns a new SyncBuffer instance
func NewSyncBuffer() *SyncBuffer {
	sb := &SyncBuffer{}
	sb.readReady = sync.NewCond(&sb.Mutex)
	return sb
}

// Read implements io.Reader.
func (b *SyncBuffer) Read(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()
	for !b.closed && b.buffer.Len() == 0 {
		// Blocks but releases Mutex for the duration of the Wait
		b.readReady.Wait()
	}
	l := b.buffer.Len()
	if l == 0 && b.closed {
		return 0, io.EOF
	}
	if l > len(p) {
		l = len(p)
	}
	// Quaranteed to not block since we limit the buffer to the available length
	return b.buffer.Read(p[:l])
}

// ReadBytes wraps bytes.Buffer.ReadBytes.
func (b *SyncBuffer) ReadBytes(delim byte) (line []byte, err error) {
	b.Lock()
	defer b.Unlock()
	// Wait without holding the mutex until buffer contains delim
	for !b.closed && !bytes.Contains(b.buffer.Bytes(), []byte{delim}) {
		// Blocks but releases Mutex for the duration of the Wait
		b.readReady.Wait()

		l := b.buffer.Len()
		if l == 0 && b.closed {
			return []byte{}, io.EOF
		}
	}
	// Quaranteed to not block since we know buffer contains 'delim'
	return b.buffer.ReadBytes(delim)
}

// Write implements io.Writer.
func (b *SyncBuffer) Write(p []byte) (n int, err error) {
	// Writer can expand the buffer so it should never block, so holding the lock is safe.
	b.Lock()
	defer b.Unlock()
	defer b.readReady.Signal()
	return b.buffer.Write(p)
}

// Close implements io.Closer.
func (b *SyncBuffer) Close() error {
	b.Lock()
	b.closed = true
	b.Unlock()
	return nil
}

// String implements fmt.Stringer.
func (b *SyncBuffer) String() string {
	b.Lock()
	defer b.Unlock()
	return b.buffer.String()
}

// ReadUntilLine consumes lines from the buffer line at at a time, and
// stops only after receiving a line starting with the 'greeting'.
// Must read data byte at a time to not consume any data after
func (b *SyncBuffer) ReadUntilLine(greeting []byte) (done chan error) {
	done = make(chan error, 1)
	go func() {
		for {
			// Read next line
			line, err := b.ReadBytes('\n')
			if bytes.HasPrefix(line, greeting) {
				close(done)
				return
			}
			if err != nil {
				done <- err
				close(done)
				return
			}
		}
	}()

	return done
}
