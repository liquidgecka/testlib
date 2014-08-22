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
	"runtime"
	"sync"
	"testing"
)

type mockT struct {
	// This is set to true if Fatal or Fatalf have been called.
	failed bool

	// This is set to true if Skip or Skipf have been called.
	skipped bool

	// These functions will be called if set to non nil.
	funcErrorf  func(format string, args ...interface{})
	funcError   func(args ...interface{})
	funcFailed  func() bool
	funcFail    func()
	funcFailNow func()
	funcFatalf  func(format string, args ...interface{})
	funcFatal   func(args ...interface{})
	funcLogf    func(format string, args ...interface{})
	funcLog     func(args ...interface{})
	funcSkipf   func(format string, args ...interface{})
	funcSkip    func(args ...interface{})
	funcSkipNow func()
	funcSkipped func() bool
}

func (m *mockT) runTest(t *testing.T, fails bool, f func()) {
	m.failed = false
	m.skipped = false

	// We run the test in another goroutine since Fatal terminates
	// running goroutines.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
	wg.Wait()
	_, file, line, _ := runtime.Caller(2)
	if fails && !m.failed {
		t.Fatalf("\n%s:%d\nTest didn't fail but should have.", file, line)
	} else if !fails && m.failed {
		t.Fatalf("\n%s:%d\nTest failed but shouldn't have.", file, line)
	}
}

func (m *mockT) CheckPass(t *testing.T, f func()) {
	m.runTest(t, false, f)
}

func (m *mockT) CheckFail(t *testing.T, f func()) {
	m.runTest(t, true, f)
}

// testing.TB compatibility functions.

func (m *mockT) Error(args ...interface{}) {
	m.failed = true
	if m.funcError != nil {
		m.funcError(args...)
	}
}

func (m *mockT) Errorf(format string, args ...interface{}) {
	m.failed = true
	if m.funcErrorf != nil {
		m.funcErrorf(format, args...)
	}
}

func (m *mockT) Failed() bool {
	if m.funcFailed != nil {
		return m.funcFailed()
	}
	return m.failed
}

func (m *mockT) Fail() {
	m.failed = true
}

func (m *mockT) FailNow() {
	m.failed = true
	runtime.Goexit()
}

func (m *mockT) Fatal(args ...interface{}) {
	m.failed = true
	if m.funcFatal != nil {
		m.funcFatal(args...)
	}
	runtime.Goexit()
}

func (m *mockT) Fatalf(format string, args ...interface{}) {
	m.failed = true
	if m.funcFatalf != nil {
		m.funcFatalf(format, args...)
	}
	runtime.Goexit()
}

func (m *mockT) Log(args ...interface{}) {
	if m.funcLog != nil {
		m.funcLog(args...)
	}
}

func (m *mockT) Logf(format string, args ...interface{}) {
	if m.funcLogf != nil {
		m.funcLogf(format, args...)
	}
}

func (m *mockT) Skipf(format string, args ...interface{}) {
	m.skipped = true
	if m.funcSkipf != nil {
		m.funcSkipf(format, args...)
	}
	runtime.Goexit()
}

func (m *mockT) Skip(args ...interface{}) {
	m.skipped = true
	if m.funcSkip != nil {
		m.funcSkip(args...)
	}
	runtime.Goexit()
}

func (m *mockT) SkipNow() {
	m.skipped = true
	if m.funcSkipNow != nil {
		m.funcSkipNow()
	}
	runtime.Goexit()
}

func (m *mockT) Skipped() bool {
	if m.funcSkipped != nil {
		return m.funcSkipped()
	}
	return m.skipped
}

func (m *mockT) private() {}
