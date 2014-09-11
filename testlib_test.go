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
	"sync"
	"testing"
)

func TestT_Finalizers(t *testing.T) {
	t.Parallel()
	T := NewT(new(mockT))

	// Some values we will use to keep track of order.
	fin1 := -1
	fin2 := -2
	fin3 := -3
	order := 0

	// Add finalizers.
	T.AddFinalizer(func() {
		fin1 = order
		order += 1
	})
	T.AddFinalizer(func() {
		fin2 = order
		order += 1
	})
	T.AddFinalizer(func() {
		fin3 = order
		order += 1
	})

	// Call the Finish function.
	T.Finish()

	if fin1 != 2 || fin2 != 1 || fin3 != 0 {
		t.Fatalf("Order of execution error! %d, %d, %d", fin1, fin2, fin3)
	}
}

func TestT_Error(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcError = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	T.Error("xxx")
	if !strings.Contains(msg, "xxx") {
		t.Fatalf("The error was not passed through.")
	} else if !m.failed {
		t.Fatalf("The test was not marked as having failed.")
	}
}

func TestT_Errorf(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcError = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	T.Errorf("xxx %s", "yyy")
	if !strings.Contains(msg, "xxx yyy") {
		t.Fatalf("The error was not passed through.")
	} else if !m.failed {
		t.Fatalf("The test was not marked as having failed.")
	}
}

func TestT_FailNow(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Ensure that the check fails.
	m.CheckFail(t, func() {
		T.FailNow()
	})
	if m.failed != true {
		t.Fatalf("The test was not marked as having failed.")
	}
}

func TestT_Failed(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	if T.Failed() {
		t.Fatalf("Failed() returned true when it shouldn't have.")
	}
	m.failed = true
	if !T.Failed() {
		t.Fatalf("Failed() returned false when it shouldn't have.")
	}
}

func TestT_Fatal(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() { T.Fatal("xxx") })
	if !strings.Contains(msg, "xxx") {
		t.Fatalf("The error was not passed through.")
	} else if !m.failed {
		t.Fatalf("The test was not marked as having failed.")
	}
}

func TestT_Fatalf(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() { T.Fatalf("xxx %s", "yyy") })
	if !strings.Contains(msg, "xxx yyy") {
		t.Fatalf("The error was not passed through.")
	} else if !m.failed {
		t.Fatalf("The test was not marked as having failed.")
	}
}

func TestT_Log(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcLog = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	T.Log("xxx")
	if !strings.Contains(msg, "xxx") {
		t.Fatalf("The message was not passed through.")
	}
}

func TestT_Logf(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcLogf = func(f string, args ...interface{}) {
		msg = fmt.Sprintf(f, args...)
	}
	T.Logf("xxx %s", "yyy")
	if !strings.Contains(msg, "xxx yyy") {
		t.Fatalf("The message was not passed through.")
	}
}

// This is a function helper for the TestT_Name test.
func BenchmarktestT_Name(T *T) string {
	return T.Name()
}

func TestT_Name(t *testing.T) {
	t.Parallel()
	T := NewT(t)

	// Simple test.
	name := T.Name()
	if name != "TestT_Name" {
		t.Fatalf("Name() returned the wrong name: %s", name)
	} else if name != T.Name() {
		t.Fatalf("Name() returned differing results.")
	}

	// Test with Benchmark as a prefix.
	T = NewT(t)
	wg := sync.WaitGroup{}
	wg.Add(1)
	name = ""
	go func() {
		name = BenchmarktestT_Name(T)
		wg.Done()
	}()
	wg.Wait()
	if name != "BenchmarktestT_Name" {
		t.Fatalf("Name() returned the wrong name: %s", name)
	}
}

func TestT_SkipNow(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	m.CheckSkips(t, func() { T.SkipNow() })
	if !m.skipped {
		t.Fatalf("The test was not marked as having skipped.")
	}
}

func TestT_Skip(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcSkip = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckSkips(t, func() { T.Skip("xxx") })
	if !strings.Contains(msg, "xxx") {
		t.Fatalf("The skip message was not passed through.")
	} else if !m.skipped {
		t.Fatalf("The test was not marked as having skipped.")
	}
}

func TestT_Skipf(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	// Capture the error.
	msg := ""
	m.funcSkip = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckSkips(t, func() { T.Skipf("xxx %s", "yyy") })
	if !strings.Contains(msg, "xxx yyy") {
		t.Fatalf("The skip message was not passed through.")
	} else if !m.skipped {
		t.Fatalf("The test was not marked as having skipped.")
	}
}

func TestT_Skipped(t *testing.T) {
	t.Parallel()
	m, T := testSetup()

	if T.Skipped() {
		t.Fatalf("Skipped() returned true when it shouldn't have.")
	}
	m.skipped = true
	if !T.Skipped() {
		t.Fatalf("Skipped() returned false when it shouldn't have.")
	}
}
