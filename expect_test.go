// Copyright 2014 Brady Catherman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testlib

import (
	"fmt"
	"strings"
	"testing"
)

func TestT_ExpectError(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckPass(t, func() { T.ExpectError(fmt.Errorf("EXPECTED")) })
	m.CheckFail(t, func() { T.ExpectError(nil, "prefix") })
	if msg == "" {
		t.Fatalf("No error message was reported.")
	} else if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("The prefix was not prepended to the message: '''%s'''", msg)
	}
}

func TestT_ExpectSuccess(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckPass(t, func() { T.ExpectSuccess(nil) })
	m.CheckFail(t, func() { T.ExpectSuccess(fmt.Errorf("ERROR"), "prefix") })
	if msg == "" {
		t.Fatalf("No error message was reported.")
	} else if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("The prefix was not prepended to the message: '''%s'''", msg)
	}
}
