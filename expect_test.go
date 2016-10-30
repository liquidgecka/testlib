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

func TestT_ExpectErrorf(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckPass(t, func() { T.ExpectErrorf(fmt.Errorf("EXPECTED"), "foo %d", 2) })
	m.CheckFail(t, func() { T.ExpectErrorf(nil, "foo %d", 3) })
	if msg == "" {
		t.Fatalf("No error message was reported.")
	} else if !strings.HasPrefix(msg, "foo 3: ") {
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

func TestT_ExpectSuccessf(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckPass(t, func() { T.ExpectSuccessf(nil, "foo %d", 4) })
	m.CheckFail(t, func() { T.ExpectSuccessf(fmt.Errorf("ERROR"), "foo %d", 5) })
	if msg == "" {
		t.Fatalf("No error message was reported.")
	} else if !strings.HasPrefix(msg, "foo 5: ") {
		t.Fatalf("The prefix was not prepended to the message: '''%s'''", msg)
	}
}

func TestT_ExpectErrorMessage(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() {
		T.ExpectErrorMessage(nil, "test")
	})
	m.CheckFail(t, func() {
		T.ExpectErrorMessage(fmt.Errorf("XXX"), "test")
	})
	m.CheckFail(t, func() {
		T.ExpectErrorMessage(fmt.Errorf("ERROR"), "test", "prefix")
	})
	if msg == "" {
		t.Fatalf("No error message was reported.")
	} else if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("The prefix was not prepended to the message: '''%s'''", msg)
	}

	m.CheckPass(t, func() {
		T.ExpectErrorMessage(fmt.Errorf("XXX"), "XXX")
	})
}

func TestT_ExpectErrorMessagef(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() {
		T.ExpectErrorMessagef(nil, "test", "foo %d", 6)
	})
	m.CheckFail(t, func() {
		T.ExpectErrorMessagef(fmt.Errorf("XXX"), "test", "foo %d", 7)
	})
	m.CheckFail(t, func() {
		T.ExpectErrorMessagef(fmt.Errorf("ERROR"), "test", "foo %d", 8)
	})
	if msg == "" {
		t.Fatalf("No error message was reported.")
	} else if !strings.HasPrefix(msg, "foo 8: ") {
		t.Fatalf("The prefix was not prepended to the message: '''%s'''", msg)
	}

	m.CheckPass(t, func() {
		T.ExpectErrorMessage(fmt.Errorf("XXX"), "XXX")
	})
}
