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
	"reflect"
	"strings"
)

// This file contains a super utility for checking the equality of structures
// or values. It is designed to compare anything to basically anything.

// Compares two values to ensure that they are equal to each other. This will
// deep inspect both values to ensure that the full structure tree is equal.
// It also walks through pointers ensuring that everything is equal.
func (t *TestLib) Equal(have, want interface{}, desc ...string) {
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}

	// Check to see if either value is nil and then verify that the are
	// either both nil, or fail if one is nil.
	haveNil := t.isNil(have)
	wantNil := t.isNil(want)
	if haveNil && wantNil {
		return
	} else if haveNil && !wantNil {
		t.Fatalf("%sExpected non nil, got nil.", prefix)
	} else if !haveNil && wantNil {
		t.Fatalf("%sExpected nil, got nil.", prefix)
	}

	// Next we need to get the value of both objects so we can compare them.
	haveValue := reflect.ValueOf(have)
	wantValue := reflect.ValueOf(want)
	visited := make(map[uintptr]*visitedNode)
	reason := t.deepEqual("", haveValue, wantValue, visited)
	if len(reason) > 0 {
		t.Fatalf("%sNot Equal\n%s", prefix, strings.Join(reason, "\n"))
	}
}

// Like Equal() except that it asserts that the two values are not equal
// to each other.
func (t *TestLib) NotEqual(have, want interface{}, desc ...string) {
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}

	// Check to see if either value is nil and then verify that the are
	// either both nil, or fail if one is nil.
	haveNil := t.isNil(have)
	wantNil := t.isNil(want)
	if haveNil && wantNil {
		t.Fatalf("%sEquality not expected, have=nil", prefix)
	} else if haveNil || wantNil {
		return
	}

	// Next we need to get the value of both objects so we can compare them.
	haveValue := reflect.ValueOf(have)
	wantValue := reflect.ValueOf(want)
	visited := make(map[uintptr]*visitedNode)
	reason := t.deepEqual("", haveValue, wantValue, visited)
	if len(reason) == 0 {
		t.Fatalf("%sValues are not expected to be equal: %#v", prefix, have)
	}
}

// Tracks access to specific pointers so we do not recurse.
type visitedNode struct {
	a1   uintptr
	a2   uintptr
	typ  reflect.Type
	next *visitedNode
}

// Returns true if the underlying object is nil.
func (t *TestLib) isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Slice:
	default:
		return false
	}
	return v.IsNil()
}

// Deep comparison. This is based on golang 1.2's reflect.Equal functionality.
func (t *TestLib) deepEqual(
	desc string, have, want reflect.Value, visited map[uintptr]*visitedNode,
) (diffs []string) {
	if !want.IsValid() && !have.IsValid() {
		return nil
	} else if !want.IsValid() && have.IsValid() {
		// This is rare, not sure how to document this better.
		return []string{
			fmt.Sprintf("%s: have invalid object.", desc),
		}
	} else if want.IsValid() && !have.IsValid() {
		// This is rare, not sure how to document this better.
		return []string{
			fmt.Sprintf("%s: wanted a valid object.", desc),
		}
	} else if want.Type() != have.Type() {
		return []string{fmt.Sprintf(
			"%s: Not the same type, have: '%s', want: '%s'",
			desc, have.Type(), want.Type())}
	}

	if want.CanAddr() && have.CanAddr() {
		addr1 := want.UnsafeAddr()
		addr2 := have.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are identical ...
		if addr1 == addr2 {
			return []string{}
		}

		// ... or already seen
		h := 17*addr1 + addr2
		seen := visited[h]
		typ := want.Type()
		for p := seen; p != nil; p = p.next {
			if p.a1 == addr1 && p.a2 == addr2 && p.typ == typ {
				return []string{}
			}
		}

		// Remember for later.
		visited[h] = &visitedNode{addr1, addr2, typ, seen}
	}

	// Checks to see if one value is nil, while the other is not.
	checkNil := func() bool {
		if want.IsNil() && !have.IsNil() {
			diffs = append(diffs, fmt.Sprintf("%s: not equal.", desc))
			diffs = append(diffs, fmt.Sprintf("have %#v", have))
			diffs = append(diffs, "want: nil")
			return true
		} else if !want.IsNil() && have.IsNil() {
			diffs = append(diffs, fmt.Sprintf("%s: not equal.", desc))
			diffs = append(diffs, "have: nil")
			diffs = append(diffs, fmt.Sprintf("want %#v", want))
			return true
		}
		return false
	}

	// Checks to see that the lengths of both objects are equal.
	checkLen := func() bool {
		if want.Len() != have.Len() {
			diffs = append(diffs, fmt.Sprintf(
				"%s: (len(have): %d, len(want): %d)",
				desc, have.Len(), want.Len()))
			diffs = append(diffs, fmt.Sprintf("have: %#v", have))
			diffs = append(diffs, fmt.Sprintf("want: %#v", want))
			return true
		}
		return false
	}

	switch want.Kind() {
	case reflect.Array:
		if !checkLen() {
			for i := 0; i < want.Len(); i++ {
				newdiffs := t.deepEqual(
					fmt.Sprintf("%s[%d]", desc, i),
					want.Index(i), have.Index(i), visited)
				diffs = append(diffs, newdiffs...)
			}
		}

	case reflect.Slice:
		if !checkNil() && !checkLen() {
			for i := 0; i < want.Len(); i++ {
				newdiffs := t.deepEqual(
					fmt.Sprintf("%s[%d]", desc, i),
					want.Index(i), have.Index(i), visited)
				diffs = append(diffs, newdiffs...)
			}
		}

	case reflect.Interface:
		if !checkNil() {
			newdiffs := t.deepEqual(
				desc, want.Elem(), have.Elem(), visited)
			diffs = append(diffs, newdiffs...)
		}

	case reflect.Ptr:
		newdiffs := t.deepEqual(
			desc, want.Elem(), have.Elem(), visited)
		diffs = append(diffs, newdiffs...)

	case reflect.Struct:
		for i, n := 0, want.NumField(); i < n; i++ {
			name := want.Type().Field(i).Name
			// Make sure that we don't print a strange error if the
			// first object given to us is a struct.
			if desc == "" {
				newdiffs := t.deepEqual(
					name, want.Field(i), have.Field(i), visited)
				diffs = append(diffs, newdiffs...)
			} else {
				newdiffs := t.deepEqual(
					fmt.Sprintf("%s.%s", desc, name),
					want.Field(i), have.Field(i), visited)
				diffs = append(diffs, newdiffs...)
			}
		}

	case reflect.Map:
		if !checkNil() {
			// Check that the keys are present in both maps.
			for _, k := range want.MapKeys() {
				if !have.MapIndex(k).IsValid() {
					// Add the error.
					diffs = append(diffs, fmt.Sprintf(
						"%sExpected key [%q] is missing.", desc, k))
					diffs = append(diffs, "have: not present")
					diffs = append(diffs, fmt.Sprintf("want: %#v",
						want.MapIndex(k)))
					continue
				}
				newdiffs := t.deepEqual(
					fmt.Sprintf("%s[%q] ", desc, k),
					want.MapIndex(k), have.MapIndex(k), visited)
				diffs = append(diffs, newdiffs...)
			}
			for _, k := range have.MapKeys() {
				if !want.MapIndex(k).IsValid() {
					// Add the error.
					diffs = append(diffs, fmt.Sprintf(
						"%sUnexpected key [%q].", desc, k))
					diffs = append(diffs,
						fmt.Sprintf("have: %#v", have.MapIndex(k)))
					diffs = append(diffs, "want: not present")
				}
			}
		}

	case reflect.Func:
		// Can't do better than this:
		checkNil()

	case reflect.String:
		// We know the underlying type is a string so calling String()
		// will return the underlying value. Trying to call Interface()
		// and assert to a string will panic.
		hstr := have.String()
		wstr := want.String()
		if len(hstr) != len(wstr) {
			return []string{
				fmt.Sprintf("%s: len(have) %d != len(want) %d.",
					desc, len(hstr), len(wstr)),
				fmt.Sprintf("have: %#v", hstr),
				fmt.Sprintf("want: %#v", wstr),
			}
		}
		for i := range hstr {
			if hstr[i] != wstr[i] {
				return []string{
					fmt.Sprintf("%s: difference at index %d.", desc, i),
					fmt.Sprintf("have: %#v", hstr),
					fmt.Sprintf("want: %#v", wstr),
				}
			}
		}

	default:
		// Specific low level types:
		switch want.Interface().(type) {
		case bool:
			hbool := have.Interface().(bool)
			wbool := want.Interface().(bool)
			if hbool != wbool {
				return []string{fmt.Sprintf(
					"%s: have %t, want %t", desc, hbool, wbool)}
			}
		case int:
			hint := have.Interface().(int)
			wint := want.Interface().(int)
			if hint != wint {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, hint, wint)}
			}
		case int8:
			hint8 := have.Interface().(int8)
			wint8 := want.Interface().(int8)
			if hint8 != wint8 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, hint8, wint8)}
			}
		case int16:
			hint16 := have.Interface().(int16)
			wint16 := want.Interface().(int16)
			if hint16 != wint16 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, hint16, wint16)}
			}
		case int32:
			hint32 := have.Interface().(int32)
			wint32 := want.Interface().(int32)
			if hint32 != wint32 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, hint32, wint32)}
			}
		case int64:
			hint64 := have.Interface().(int64)
			wint64 := want.Interface().(int64)
			if hint64 != wint64 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, hint64, wint64)}
			}
		case uint:
			huint := have.Interface().(uint)
			wuint := want.Interface().(uint)
			if huint != wuint {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, huint, wuint)}
			}
		case uint8:
			huint8 := have.Interface().(uint8)
			wuint8 := want.Interface().(uint8)
			if huint8 != wuint8 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, huint8, wuint8)}
			}
		case uint16:
			huint16 := have.Interface().(uint16)
			wuint16 := want.Interface().(uint16)
			if huint16 != wuint16 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, huint16, wuint16)}
			}
		case uint32:
			huint32 := have.Interface().(uint32)
			wuint32 := want.Interface().(uint32)
			if huint32 != wuint32 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, huint32, wuint32)}
			}
		case uint64:
			huint64 := have.Interface().(uint64)
			wuint64 := want.Interface().(uint64)
			if huint64 != wuint64 {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, huint64, wuint64)}
			}
		case uintptr:
			huintptr := have.Interface().(uintptr)
			wuintptr := want.Interface().(uintptr)
			if huintptr != wuintptr {
				return []string{fmt.Sprintf(
					"%s: have %d, want %d", desc, huintptr, wuintptr)}
			}
		case float32:
			hf32 := have.Interface().(float32)
			wf32 := want.Interface().(float32)
			if hf32 != wf32 {
				return []string{fmt.Sprintf(
					"%s: have %f, want %f", desc, hf32, wf32)}
			}
		case float64:
			hf64 := have.Interface().(float64)
			wf64 := want.Interface().(float64)
			if hf64 != wf64 {
				return []string{fmt.Sprintf(
					"%s: have %f, want %f", desc, hf64, wf64)}
			}
		default:
			// Normal equality suffices
			if !reflect.DeepEqual(want.Interface(), have.Interface()) {
				return []string{
					fmt.Sprintf("%s: not equal.", desc),
					fmt.Sprintf("have: %#v", have),
					fmt.Sprintf("want: %#v", want),
				}
			}
		}
	}

	// This shouldn't ever be reached.
	return diffs
}
