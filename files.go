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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// This file contains functions for dealing with files.

// Calling this function will create a process specific temporary directory
// that will be cleaned up when the process terminates. Data in this directory
// will survive the entire lifetime of the testing process.
//
// Files in this directory are cleaned up by a child process that is forked
// from the running process so that nothing can stop them from being cleaned.
func (t *T) RootTempDir() string {
	testLibRootDirOnce.Do(func() {
		var err error
		var reader *os.File
		mode := os.FileMode(0777)
		testLibRootDir, err = ioutilTempDir("", "golang-testlib")
		t.NotEqual(testLibRootDir, "")
		t.ExpectSuccess(err)
		t.ExpectSuccess(osChmod(testLibRootDir, mode))
		reader, testLibRootDirStdin, err = os.Pipe()
		t.ExpectSuccess(err)
		cmd := exec.Command(os.Args[0], testInterceptorArg,
			testLibRootDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = reader
		t.ExpectSuccess(cmd.Start())
		t.ExpectSuccess(reader.Close())
	})
	return testLibRootDir

}

// Creates a temporary directory for this specific test which will be cleaned
// once the test has finished executing. This calls RootTempDir() to create the
// base directory.
func (t *T) TempDirMode(mode os.FileMode) string {
	f, err := ioutilTempDir(t.RootTempDir(), t.Name())
	t.ExpectSuccess(err)
	t.NotEqual(f, "")
	t.ExpectSuccess(osChmod(f, mode))
	t.AddFinalizer(func() {
		osRemoveAll(f)
	})
	return f
}

// Like TempDirMode except this sets the default mode to 0755.
func (t *T) TempDir() string {
	return t.TempDirMode(os.FileMode(0755))
}

// Creates a temporary file in a temporary directory with a specific mode
// set on it. This will return the file descriptor of the open file.
func (t *T) TempFileMode(mode os.FileMode) *os.File {
	f, err := ioutilTempFile(t.RootTempDir(), t.Name())
	t.ExpectSuccess(err)
	t.NotEqual(f, nil)
	t.ExpectSuccess(osChmod(f.Name(), mode))
	name := f.Name()
	t.AddFinalizer(func() {
		osRemove(name)
	})
	return f
}

// Like TempFileMode except that it uses a default mode of 0644.
func (t *T) TempFile() *os.File {
	return t.TempFileMode(os.FileMode(0644))
}

// Makes a temporary file with the given string as contents. This returns
// the name of the created file.
func (t *T) WriteTempFileMode(contents string, mode os.FileMode) string {
	f := t.TempFileMode(mode)
	name := f.Name()
	_, err := io.WriteString(f, contents)
	t.ExpectSuccess(err)
	t.ExpectSuccess(f.Close())
	return name
}

// Like WriteTempFileMode except this uses the default temp file mode.
func (t *T) WriteTempFile(contents string) string {
	return t.WriteTempFileMode(contents, 0644)
}

// -------------------------------
// Temporary Dir Cleanup Internals
// -------------------------------

// If the process is started with this string as its first argument and
// a directory as its second argument then the startup flow will be
// intercepted to allow the process to clean up after the parent.
const testInterceptorArg = "wledfhs9d8fs9id"

// This function is used to intercept the process startup and check to see if
// if its a clean up process. Args will be os.Args, and reader will be
// os.Stdin.
func initRootTempDir(args []string, reader io.Reader) {
	if len(args) != 3 {
		return
	} else if args[1] != testInterceptorArg {
		return
	}

	// Only remove files if it is in the operating systems temporary directory
	// structure. This is a safety trap to prevent us from accidentally
	// removing files critical to the system.
	if !strings.HasPrefix(args[2], osTempDir()) {
		fmtFprintf(os.Stderr, "Refusing to clean a non temporary directory: "+
			"%s since it is not under %s\n", args[2], osTempDir())
		osExit(1)
		return
	}

	// The parent process holds our stdin open until it dies, once that happens
	// we need to remove the directory.
	if _, err := ioutil.ReadAll(reader); err != nil {
		fmtFprintf(
			os.Stderr, "Error cleaning up directory %s: %s\n",
			args[2], err)
		osExit(2)
	} else if err := osRemoveAll(args[2]); err != nil {
		fmtFprintf(
			os.Stderr, "Error cleaning up directory %s: %s\n",
			args[2], err)
		osExit(3)
	} else {
		osExit(0)
	}
}

// On startup call the initRootTempDir() function.
func init() {
	initRootTempDir(os.Args, os.Stdin)
}

// The private global variables that stores the root directories location
// so it is preserved between tests.
var (
	testLibRootDir      string
	testLibRootDirOnce  sync.Once
	testLibRootDirStdin io.Writer
)
