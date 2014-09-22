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

package testlib_test

import (
	"os"
	"testing"

	"github.com/liquidgecka/testlib"
)

// Note that the function name needs to not start with an underscore. That
// is a relic of how golang's "playable examples" work.
func _TestBasicExample(t *testing.T) {
	// Basic set of the testlib wrapper.
	T := testlib.NewT(t)
	defer T.Finish()

	// We can create a temporary file fairly easily. This file will be cleaned
	// up automatically after this test runs.
	fd := T.TempFile()

	// We can check equality of two objects, including full structures, lists,
	// maps, etc. If the objects are not equal then the test will be
	// stopped via a call to Fatalf()
	T.Equal(1, 1)
	T.Equal([]string{"a", "b", "c"}, []string{"a", "b", "c"})
	T.Equal(T, T)

	// We can also assert inequality.
	T.NotEqual(10, 1)

	// It is also possible to do a very simple error check that will Fatalf the
	// test if an error is returned.
	_, err := os.Open(fd.Name())
	T.ExpectSuccess(err)

	// Likewise we can assert that something SHOULD return an error, and lack of
	// an error is actually a problem.
	_, err = os.Open("/a_file_that_doesnot_exist")
	T.ExpectError(err)

	// And if we want to we can just flat out Fatalf with a message.
	T.Fatalf("The test has failed.")
}

// Simple documentation showing the basic usage of this library with some really
// easy to follow examples.
func Example() {
	// Ignore this function.. It is required to make this file a "playable
	// example" in godoc.
}
