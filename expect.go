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
)

// This file contains functions to assert specific expectations.

// This call will check that the given error object is non nil and if it is
// not it will automatically Fatalf the test with a message.
func (t *T) ExpectError(err error, desc ...string) {
	if err != nil {
		return
	}
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}
	t.Fatalf("%sError not returned when one was expected.", prefix)
}

// ExpectErrorf checks if the given error object is non-nil and if it is
// not then it will Fatalf the test with a message formed by *f formatting.
func (t *T) ExpectErrorf(err error, spec string, args ...interface{}) {
	if err != nil {
		return
	}
	prefix := fmt.Sprintf(spec, args...) + ": "
	t.Fatalf("%sError not returned when one was expected.", prefix)
}

// Checks to see that the given error object is nil. This is handy for
// performing quick checks on calls that are expected to work.
func (t *T) ExpectSuccess(err error, desc ...string) {
	if err == nil {
		return
	}
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}
	t.Fatalf("%sUnexpected error encountered: %#v (%s)",
		prefix, err, err.Error())
}

// ExpectSuccessf checks that the given error object is nil. If non-nil
// then report as a test failure.
func (t *T) ExpectSuccessf(err error, spec string, args ...interface{}) {
	if err == nil {
		return
	}
	prefix := fmt.Sprintf(spec, args...) + ": "
	t.Fatalf("%sUnexpected error encountered: %#v (%s)",
		prefix, err, err.Error())
}

// Fails if the error message does not contain the given string.
func (t *T) ExpectErrorMessage(err error, msg string, desc ...string) {
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}
	if err == nil {
		t.Fatalf("%sExpected error was not returned.", prefix)
	} else if !strings.Contains(err.Error(), msg) {
		t.Fatalf("%sError message didn't contain the expected message:\n"+
			"Error message=%s\nExpected string=%s", prefix, err.Error(), msg)
	}
}
