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

// +build !privileged_tests

package utils

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type UtilsSuite struct{}

var _ = check.Suite(&UtilsSuite{})

func (b *UtilsSuite) TestExec(c *check.C) {
	failIfCalled := func(f string, a ...interface{}) {
		c.Error("log function should not be called")
	}
	_, err := Exec(failIfCalled, "true")
	c.Assert(err, check.IsNil)

	logCalled := 0
	countLog := func(f string, a ...interface{}) { logCalled++ }
	_, err = Exec(countLog, "false")
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(logCalled, check.Equals, 1)

	logCalled = 0
	_, err = Exec(countLog, "sh", "-c", "'echo foo; exit 1'")
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(logCalled, check.Equals, 2)
}
