# TestLib

TestLib is a library that improves on Go's internal unit testing framework.

## Installation

```bash
go get github.com/liquidgecka/testlib
```

## Usage

Within a Test function you can setup and use a TestLib instance. This will
provide very short and easy to follow functions for doing things within the
unit test. Ideally you can pass the TestLib object to helper functions as
well.

```go
package example

import (
    "os"
    "testing"
    "github.com/liquidgecka/testlib"
)

func TestSomething(t *testing.T) {
    T := testlib.NewT(t)
    defer T.Finish()

    // Create a temporary file that will be cleaned up when the test finishes.
    // If there is an error making the temporary file then the test is
    // automatically terminated via a call to Fatal.
    tmpFile := T.WriteTempFile("contents")

    // Use the os.ReadFile call to read the contents of the file. Since
    // the file should exist we can use ExpectSuccess to assert that the
    // call succeeded without error. If err is not nil then the test will be
    // terminated with a call to Fatal.
    contents, err := os.ReadFile(tmpFile)
    T.ExpectSuccess(err)

    // Next we can check that the read results match the contents we wrote
    // earlier. If not then this call will automatically call Fatal. Note that
    // if the objects are NOT the same then this will also output a difference
    // between them to make it easier to see what happened.
    T.Equal(contents, []byte("contents"))
}
```

### Generic Helpers

Go does not support type parameters on methods, so `testlib` provides
package-level generic functions alongside the existing `*T` methods. Each takes
`T` as its first argument:

```go
// Compile-time type checking and untyped constant inference for equality.
testlib.Equal(T, int64(42), 42)
testlib.NotEqual(T, int64(10), 1)

// Generic panic expectation.
testlib.ExpectPanic(T, func() { panic("boom") }, "boom")

// Accepts either string or []byte (including defined types with those underlying types).
tmpFile := testlib.WriteTempFile(T, []byte("contents"))
```

## Testing
[![CI](https://github.com/liquidgecka/testlib/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/liquidgecka/testlib/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/liquidgecka/testlib.svg)](https://pkg.go.dev/github.com/liquidgecka/testlib)

## Contribution

I gladly accept PR's and will work on issues if filed.

## License (Apache 2)

Copyright 2014 Brady Catherman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
