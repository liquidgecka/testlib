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
	"io/ioutil"
	"os"
	"testing"
)

func TestTempDirMode(t *testing.T) {
	// Test 1: ioutil.TempDir() failure.
	m, T := testSetup()
	m.CheckFail(t, func() {
		ioutilTempDir = func(a, b string) (string, error) {
			return "", fmt.Errorf("Expected")
		}
		defer func() { ioutilTempDir = ioutil.TempDir }()
		T.TempDirMode(os.FileMode(0644))
	})

	// Test 2: ioutil.TempDir() returns "".
	m, T = testSetup()
	m.CheckFail(t, func() {
		ioutilTempDir = func(a, b string) (string, error) {
			return "", nil
		}
		defer func() { ioutilTempDir = ioutil.TempDir }()
		T.TempDirMode(os.FileMode(0644))
	})

	// Test 3: os.Chmod fails.
	m, T = testSetup()
	m.CheckFail(t, func() {
		osChmod = func(a string, m os.FileMode) error {
			return fmt.Errorf("Expected error")
		}
		defer func() { osChmod = os.Chmod }()
		T.TempDirMode(os.FileMode(0644))
	})

	// Test 4: Success.
	m, T = testSetup()
	var dir string
	m.CheckPass(t, func() {
		dir = T.TempDirMode(os.FileMode(0715))
	})
	expectedMode := os.FileMode(0715) + os.ModeDir
	if dir == "" {
		t.Fatalf("Returned directory can not be empty.")
	} else if stat, err := os.Stat(dir); err != nil {
		t.Fatalf("Error stating the returned directory: %s", err)
	} else if stat.Mode() != expectedMode {
		t.Fatalf("Invalid mode on the created directory: %s", stat.Mode())
	}

	// Ensure that the Finalizer function runs.
	ran := false
	m.CheckPass(t, func() {
		osRemoveAll = func(s string) error {
			ran = true
			return os.RemoveAll(s)
		}
		defer func() { osRemoveAll = os.RemoveAll }()
		T.Finish()
	})
	if ran != true {
		t.Fatalf("The Finalized failed to run.")
	} else if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("The file %s shouldn't exist.", dir)
	}
}

func TestT_TempDir(t *testing.T) {
	m, T := testSetup()
	var dir string
	m.CheckPass(t, func() {
		dir = T.TempDir()
	})
	expectedMode := os.FileMode(0755) + os.ModeDir
	if dir == "" {
		t.Fatalf("Returned directory can not be empty.")
	} else if stat, err := os.Stat(dir); err != nil {
		t.Fatalf("Error stating the returned directory: %s", err)
	} else if stat.Mode() != expectedMode {
		t.Fatalf("Invalid mode on the created directory: %s", stat.Mode())
	}
	T.Finish()
}

func TestTempFileMode(t *testing.T) {
	// Test 1: ioutil.TempFile() failure.
	m, T := testSetup()
	m.CheckFail(t, func() {
		ioutilTempFile = func(a, b string) (*os.File, error) {
			return nil, fmt.Errorf("Expected")
		}
		defer func() { ioutilTempFile = ioutil.TempFile }()
		T.TempFileMode(os.FileMode(0644))
	})

	// Test 2: ioutil.TempFile() returns nil.
	m, T = testSetup()
	m.CheckFail(t, func() {
		ioutilTempFile = func(a, b string) (*os.File, error) {
			return nil, nil
		}
		defer func() { ioutilTempFile = ioutil.TempFile }()
		T.TempFileMode(os.FileMode(0644))
	})

	// Test 3: os.Chmod fails.
	m, T = testSetup()
	m.CheckFail(t, func() {
		osChmod = func(a string, m os.FileMode) error {
			return fmt.Errorf("Expected error")
		}
		defer func() { osChmod = os.Chmod }()
		T.TempFileMode(os.FileMode(0644))
	})

	// Test 4: Success.
	m, T = testSetup()
	var file *os.File
	m.CheckPass(t, func() {
		file = T.TempFileMode(os.FileMode(0614))
	})
	expectedMode := os.FileMode(0614)
	if file == nil {
		t.Fatalf("Returned file can not be nil.")
	} else if stat, err := os.Stat(file.Name()); err != nil {
		t.Fatalf("Error stating the returned file: %s", err)
	} else if stat.Mode() != expectedMode {
		t.Fatalf("Invalid mode on the created file: %s", stat.Mode())
	}

	// Ensure that the Finalizer function runs.
	ran := false
	m.CheckPass(t, func() {
		osRemove = func(s string) error {
			ran = true
			return os.Remove(s)
		}
		defer func() { osRemove = os.Remove }()
		T.Finish()
	})
	if ran != true {
		t.Fatalf("The Finalized failed to run.")
	} else if _, err := os.Stat(file.Name()); !os.IsNotExist(err) {
		t.Fatalf("The file %s shouldn't exist.", file.Name())
	}
}

func TestT_TempFile(t *testing.T) {
	m, T := testSetup()
	var file *os.File
	m.CheckPass(t, func() {
		file = T.TempFile()
	})
	expectedMode := os.FileMode(0644)
	if file == nil {
		t.Fatalf("Returned file can not be nil.")
	} else if stat, err := os.Stat(file.Name()); err != nil {
		t.Fatalf("Error stating the returned file: %s", err)
	} else if stat.Mode() != expectedMode {
		t.Fatalf("Invalid mode on the created file: %s", stat.Mode())
	}
	T.Finish()
}

func TestT_WhiteTempFileMode(t *testing.T) {
	m, T := testSetup()
	var file string
	m.CheckPass(t, func() {
		file = T.WriteTempFileMode("contents", os.FileMode(0614))
	})
	if file == "" {
		t.Fatalf("Returned file can not be empty.")
	} else if stat, err := os.Stat(file); err != nil {
		t.Fatalf("Error statting returned file %s: %s", file, err)
	} else if stat.Mode() != os.FileMode(0614) {
		t.Fatalf("Returned file '%s' has the wrong mode: %x", file, stat.Mode())
	} else if contents, err := ioutil.ReadFile(file); err != nil {
		t.Fatalf("Error reading %s: %s", file, err)
	} else if string(contents) != "contents" {
		t.Fatalf("File contained the wrong contents")
	}
	T.Finish()
}

func TestT_WhiteTempFile(t *testing.T) {
	m, T := testSetup()
	var file string
	m.CheckPass(t, func() {
		file = T.WriteTempFile("contents")
	})
	if file == "" {
		t.Fatalf("Returned file can not be empty.")
	} else if stat, err := os.Stat(file); err != nil {
		t.Fatalf("Error statting returned file %s: %s", file, err)
	} else if stat.Mode() != os.FileMode(0644) {
		t.Fatalf("Returned file '%s' has the wrong mode: %x", file, stat.Mode())
	} else if contents, err := ioutil.ReadFile(file); err != nil {
		t.Fatalf("Error reading %s: %s", file, err)
	} else if string(contents) != "contents" {
		t.Fatalf("File contained the wrong contents")
	}
	T.Finish()
}
