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
	"time"
)

func TestT_TryUntilPasses(t *testing.T) {
	t.Parallel()
	m, T := testSetup()
	m.CheckPass(t, func() {
		T.TryUntil(func() bool { return true }, time.Second)
	})
}

func TestT_TryUntilfPasses(t *testing.T) {
	t.Parallel()
	m, T := testSetup()
	m.CheckPass(t, func() {
		T.TryUntilf(func() bool { return true }, time.Second, "foo %d", 1)
	})
}

func TestT_TryUntilFails(t *testing.T) {
	t.Parallel()
	m, T := testSetup()
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() {
		T.TryUntil(func() bool { return false }, time.Second/100, "prefix")
	})
	if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("Error message did not contain the prefix: '''%s'''", msg)
	}
}

func TestT_TryUntilfFails(t *testing.T) {
	t.Parallel()
	m, T := testSetup()
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() {
		T.TryUntilf(func() bool { return false }, time.Second/100, "foo %d", 2)
	})
	if !strings.HasPrefix(msg, "foo 2: ") {
		t.Fatalf("Error message did not contain the prefix: '''%s'''", msg)
	}
}

func TestT_TryUntilYield(t *testing.T) {
	t.Parallel()
	m, T := testSetup()
	l := sync.Mutex{}
	unlocked := false
	go func() {
		time.Sleep(time.Second / 100)
		l.Lock()
		unlocked = true
		l.Unlock()
	}()
	getUnlocked := func() bool {
		l.Lock()
		defer l.Unlock()
		time.Sleep(time.Millisecond)
		return unlocked
	}
	m.CheckPass(t, func() {
		T.TryUntil(getUnlocked, time.Second)
	})
}
