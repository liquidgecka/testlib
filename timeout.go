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
	"runtime"
	"strings"
	"time"
)

// This file contains timeout functions to help with unit testing.

// This function will try to run the function 'f' until 'timeout' duration
// has elapsed. If 'f' returns true this will return, otherwise if 'f' returns
// false for the whole period then this will automatically call Fatal to
// terminate the test.
func (t *T) TryUntil(
	f func() bool, timeout time.Duration, desc ...string,
) {
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}

	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		if f() {
			return
		}
		// Yield the processor so that other goroutines have a chance to work.
		// This is necessary since the function may not actually sleep at
		// all.
		runtime.Gosched()
	}

	t.Fatalf("%sTimeout after %s", prefix, timeout)
}

// TryUntilf is the same as TryUntil but uses Printf formatting for the
// description message.
func (t *T) TryUntilf(
	f func() bool, timeout time.Duration, spec string, args ...interface{},
) {
	prefix := fmt.Sprintf(spec, args...) + ": "

	end := time.Now().Add(timeout)
	for time.Now().Before(end) {
		if f() {
			return
		}
		// Yield the processor so that other goroutines have a chance to work.
		// This is necessary since the function may not actually sleep at
		// all.
		runtime.Gosched()
	}

	t.Fatalf("%sTimeout after %s", prefix, timeout)
}
