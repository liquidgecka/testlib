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
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

//
// RootTempDir Tests
//

// In order to check the RootTempDir() function completely we need to actually
// use subprocesses. This is because this call only works when the process
// that calls it terminates. As such we will use the same trick it is using
// to ensure that we can start children successfully.

// Each element of this map represents a way that the test can error. In all
// of these cases the temporary directory must still be cleaned up at the end
// of the test.
var rootTempDirExitModes map[string]func() = map[string]func(){
	"panic":       func() { panic(fmt.Errorf("Expected panic.")) },
	"osexit0":     func() { os.Exit(0) },
	"osexit1":     func() { os.Exit(1) },
	"syscallexit": testSyscallExit,
}

// The interceptor we will use here.
const testInterceptorArg2 = testInterceptorArg + "_2"

// This call checks to see if this process is the child started with the
// testInterceptorArg2 token. If so we will prevent main() from starting.
func init() {
	if len(os.Args) != 3 {
		return
	} else if os.Args[1] != testInterceptorArg2 {
		return
	}

	// If we get here we are the child of the testing process. This means
	// that we need to run the actual RootTempDir() call and report the
	// outputs.

	// Setup some mock test objects.
	t := &mockT{}
	T := NewT(t)
	defer T.Finish()

	// Next we replace the error functions in t with ones that will actually
	// print the error so we can capture it and see that something has
	// gone wrong.
	t.funcErrorf = func(f string, args ...interface{}) {
		fmt.Printf(f, args...)
	}
	t.funcError = func(args ...interface{}) { fmt.Print(args...) }
	t.funcFatalf = t.funcErrorf
	t.funcFatal = t.funcError

	// And now we make the call, reporting the output back to the user.
	dir := T.RootTempDir()
	fmt.Printf("%s:%s\n", testInterceptorArg2, dir)

	// Okay, now we need to figure out which test pattern we need to actually
	// run. All the tests should end the same way, but each is a unique
	// termination mode.
	if fn, ok := rootTempDirExitModes[os.Args[2]]; ok {
		fn()
		os.Exit(1)
	} else {
		panic("Unknown test passed as an argument.")
	}
}

func TestRootTempDir(t *testing.T) {
	// This function runs a test with a given string name.
	runTest := func(name string) {
		cmd := exec.Command(os.Args[0], testInterceptorArg2, name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				t.Fatalf("%s: Error executing the command: %s", name, err)
			}
		}

		// Convert the output to a list of lines.
		lines := strings.Split(string(output), "\n")

		// And now ensure that the process actually executed correctly by
		// ensuring that the first line contains the interceptor arg,
		// then a colon, then the temp directory and a return.
		if !strings.HasPrefix(lines[0], testInterceptorArg2+":") {
			t.Fatalf("%s: Child didn't execute properly.\nOutput: %s",
				name, strings.Join(lines, "\n"))
		}

		// Get just the directory part:
		dir := strings.SplitN(lines[0], ":", 2)[1]

		// The test ran successfully, now we need to ensure that the root
		// directory is removed at some point. We need to try several times
		// since the grand-child process may still be in the process of
		// removing the directory.
		end := time.Now().Add(time.Second * 5)
		for {
			if time.Now().After(end) {
				t.Fatalf("%s: Timed out waiting for the %s to be removed.",
					name, dir)
			} else if _, err := os.Stat(dir); err == nil {
				time.Sleep(time.Second / 100)
			} else if !os.IsNotExist(err) {
				t.Fatalf("%s: Error stating file: %s", name, err)
			} else {
				break
			}
		}

		// Success!
		return
	}

	// Run each test in the map.
	for name, _ := range rootTempDirExitModes {
		runTest(name)
	}

	// For complete coverage we call RootTempDir() just once here.
	NewT(t).RootTempDir()
}

func TestRoomTempDirInit(t *testing.T) {
	// Ensure that the defaults get set again once this test finishes.
	defer func() {
		fmtFprintf = fmt.Fprintf
		osExit = os.Exit
		osRemoveAll = os.RemoveAll
		osTempDir = os.TempDir
	}()

	// Replace some system calls with stable testing calls.
	exited := 0
	osExit = func(i int) { exited = i }
	osTempDir = func() string { return "PREFIX" }
	osRemoveAll = func(dir string) error {
		if dir == "PREFIX/WORK" {
			return nil
		}
		return fmt.Errorf("Expected error")
	}
	fmtFprintf = func(w io.Writer, s string, args ...interface{}) (int, error) {
		return 0, nil
	}
	r := strings.NewReader("")

	// Test the arg count rule.
	exited = -1
	initRootTempDir([]string{}, r)
	if exited != -1 {
		t.Fatalf("initRootTempDir should not have exited.")
	}

	// Test the arg token rule.
	initRootTempDir([]string{"1", "2", "3"}, r)
	if exited != -1 {
		t.Fatalf("initRootTempDir should not have exited.")
	}

	// Test that a bad prefix causes the process to exit with code 1.
	initRootTempDir([]string{"argv0", testInterceptorArg, "BAD_PREFIX"}, r)
	if exited != 1 {
		t.Fatalf("initRootTempDir should have exited with code 1.")
	}

	// Test that an error while reading causes a exit code of 2.
	exited = -1
	pr, pw := io.Pipe()
	pw.CloseWithError(fmt.Errorf("expected"))
	initRootTempDir([]string{"argv0", testInterceptorArg, "PREFIX"}, pr)
	if exited != 2 {
		t.Fatalf("initRootTempDir should have exited with code 2.")
	}

	// Check that an error in the removeall stage returns code 3.
	exited = -1
	pw.CloseWithError(fmt.Errorf("expected"))
	initRootTempDir([]string{"argv0", testInterceptorArg, "PREFIX"}, r)
	if exited != 3 {
		t.Fatalf("initRootTempDir should have exited with code 3.")
	}

	// And lastly check that all the right stuff workd.
	exited = -1
	pw.CloseWithError(fmt.Errorf("expected"))
	initRootTempDir([]string{"argv0", testInterceptorArg, "PREFIX/WORK"}, r)
	if exited != 0 {
		t.Fatalf("initRootTempDir should have exited with code 0.")
	}

}
