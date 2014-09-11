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
	"io/ioutil"
	"os"
)

var fmtFprintf func(io.Writer, string, ...interface{}) (int, error) = fmt.Fprintf
var ioutilTempDir func(string, string) (string, error) = ioutil.TempDir
var ioutilTempFile func(string, string) (*os.File, error) = ioutil.TempFile
var osChmod func(string, os.FileMode) error = os.Chmod
var osExit func(int) = os.Exit
var osRemoveAll func(string) error = os.RemoveAll
var osRemove func(string) error = os.Remove
var osTempDir func() string = os.TempDir
