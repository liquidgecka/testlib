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
	"path"
	"runtime"
	"strings"
)

// This is a mirror of testing.TB except that it does not include the private
// function. This allows us to use a mocked testing library, which for some
// reason the Go devs appear to think is a really bad thing however ends up
// being necessary to test libraries designed to work with testing.. =/
type testingTB interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

// This is the tracking object used to store various variables and bits needed
// for tracking the test. Each test should create this object and then destroy
// it when the test finishes by using a defer() call.
type TestLib struct {
	// This stores a reference to the testing.T or testing.B object used to
	// create the initial test.
	t testingTB

	// This is the name of the test (function name) that created this TestLib
	// instance.
	name string

	// This is a list of functions that need to be run when the test finishes,
	// regardless of how the test finishes. This allows us to setup cleanup
	// functionality without imposing more than a single defer on the
	// calling test function.
	finalizers []func()
}

// This should be called when the test is started. It will initialize a
// TestLib instance for the specific test.
func NewTestLib(t testingTB) *TestLib {
	return &TestLib{t: t}
}

// This function should be immediately added as a defer after initializing
// the TestLib structure. This will clean up after the test.
func (t *TestLib) Finish() {
	for i := len(t.finalizers) - 1; i >= 0; i-- {
		t.finalizers[i]()
	}

	// FIXME: Logging?
}

// This adds a function that will be called once the test completes. The
// functions are called when the test finishes in reverse order from how
// they were added.
func (t *TestLib) AddFinalizer(f func()) {
	t.finalizers = append(t.finalizers, f)
}

// This call will make a stack trace message for the Fatal/Fatalf and
// Error/Errorf function calls. This will insert "msg" at the top of the
// stack and return a string.
func (t *TestLib) makeStack(msg string) string {
	lines := make([]string, 0, 100)
	lines = append(lines, msg)

	// We want to eliminate any part of the stack trace that includes the
	// testlib module since that is just noise. In order to do that we
	// get the directory of this function from the runtime module and compare
	// all lines against that directory.
	_, thisfile, _, _ := runtime.Caller(0)
	thisdir := path.Dir(thisfile)

	// Now walk through generating a stack trace.
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		} else if path.Dir(file) == thisdir {
			continue
		}
		trace := fmt.Sprintf("%s:%d", file, line)
		lines = append(lines, trace)
	}

	return strings.Join(lines, "\n")
}

// Wraps the testing.T.Error function call in order to provide full stack
// traces around the error being reported rather than just the calling line.
func (t *TestLib) Error(args ...interface{}) {
	t.t.Error(t.makeStack(fmt.Sprint(args...)))
}

// Like Error() but for formatted strings.
func (t *TestLib) Errorf(format string, args ...interface{}) {
	t.t.Error(t.makeStack(fmt.Sprintf(format, args...)))
}

// This is a wrapper around Fatal in order to provide much nicer output.
// Specifically this will output a full stack trace rather than just the
// failing line. This is optimal because it makes debugging when loops
// or helper functions are used far easier.
func (t *TestLib) Fatal(args ...interface{}) {
	// TODO: Add pre-failure helper functions.
	t.t.Fatal(t.makeStack(fmt.Sprint(args...)))
}

// Like Fatal() except for formatted strings.
func (t *TestLib) Fatalf(format string, args ...interface{}) {
	// TODO: Add pre-failure helper functions.
	t.t.Fatal(t.makeStack(fmt.Sprintf(format, args...)))
}

// Gets the function name of the running test. This is useful since there is
// no other programatic way of finding out which test is running.
func (t *TestLib) Name() string {
	// If we already calculated this then just return the cached value.
	if t.name != "" {
		return t.name
	}

	// Next we need to walk through the call stack checking the name of
	// each function that is running. Each function follows the form:
	// "module/file.function" so we need to split on . and take the
	// last element. We then need to find the first function named "Test*"
	// in the list.
	for i := 0; true; i++ {
		if pc, _, _, ok := runtime.Caller(i); !ok {
			break
		} else if f := runtime.FuncForPC(pc); f == nil {
			continue
		} else {
			name := f.Name()
			index := strings.LastIndex(name, ".")
			if index < 1 {
				continue
			}
			name = name[index:]
			if strings.HasPrefix(name, "Test") {
				t.name = name
			} else if strings.HasPrefix(name, "Benchmark") {
				t.name = name
			}
		}
	}

	return t.name
}
