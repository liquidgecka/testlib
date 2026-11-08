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
func (t *T) Equal(have, want interface{}, desc ...string) {
	t.Helper()
	t.EqualWithIgnores(have, want, nil, desc...)
}

// Equalf is the same as Equal but uses Printf style formatting to construct
// the description message.
func (t *T) Equalf(have, want interface{}, spec string, args ...interface{}) {
	t.Helper()
	t.EqualWithIgnoresf(have, want, nil, spec, args...)
}

// Like Equal, except the third argument is a list of paths that should not
// be considered. This can be used to mask out expected differences in objects.
//
// The ignores list contains strings which match the output format of Equal.
func (t *T) EqualWithIgnores(
	have, want interface{}, ignores []string, desc ...string,
) {
	t.Helper()
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}
	t.equalWithIgnoresPrefix_(have, want, ignores, prefix)
}

// EqualWithIgnoresf is the same as EqualWithIgnores but uses Printf
// formatting for the description.
func (t *T) EqualWithIgnoresf(
	have, want interface{}, ignores []string, spec string, args ...interface{},
) {
	t.Helper()
	prefix := fmt.Sprintf(spec, args...) + ": "
	t.equalWithIgnoresPrefix_(have, want, ignores, prefix)
}

func (t *T) equalWithIgnoresPrefix_(
	have, want interface{}, ignores []string, prefix string,
) {
	t.Helper()
	// Check to see if either value is nil and then verify that the are
	// either both nil, or fail if one is nil.
	haveNil := t.isNil(have)
	wantNil := t.isNil(want)
	if haveNil && wantNil {
		return
	} else if haveNil && !wantNil {
		t.Fatalf("%sExpected non nil, got nil.", prefix)
	} else if !haveNil && wantNil {
		t.Fatalf("%sExpected nil, got non nil.", prefix)
	}

	// Next we need to get the value of both objects so we can compare them.
	haveValue := reflect.ValueOf(have)
	wantValue := reflect.ValueOf(want)
	visited := make(map[uintptr]*visitedNode)
	reason := t.deepEqual("", haveValue, wantValue, ignores, visited)
	if len(reason) > 0 {
		t.Fatalf("%sNot Equal\n%s", prefix, strings.Join(reason, "\n"))
	}
}

// Like Equal() except that it asserts that the two values are not equal
// to each other.
func (t *T) NotEqual(have, unwanted interface{}, desc ...string) {
	t.Helper()
	prefix := ""
	if len(desc) > 0 {
		prefix = strings.Join(desc, " ") + ": "
	}
	t.notEqualPrefix_(have, unwanted, prefix)
}

// NotEqualf is NotEqual using Printf style format strings.
func (t *T) NotEqualf(have, unwanted interface{}, spec string, args ...interface{}) {
	t.Helper()
	prefix := fmt.Sprintf(spec, args...) + ": "
	t.notEqualPrefix_(have, unwanted, prefix)
}

func (t *T) notEqualPrefix_(have, unwanted interface{}, prefix string) {
	t.Helper()
	// Check to see if either value is nil and then verify that the are
	// either both nil, or fail if one is nil.
	haveNil := t.isNil(have)
	unwantedNil := t.isNil(unwanted)
	if haveNil && unwantedNil {
		t.Fatalf("%sEquality not expected, have=nil", prefix)
	} else if haveNil || unwantedNil {
		return
	}

	// Next we need to get the value of both objects so we can compare them.
	haveValue := reflect.ValueOf(have)
	unwantedValue := reflect.ValueOf(unwanted)
	visited := make(map[uintptr]*visitedNode)
	reason := t.deepEqual("", haveValue, unwantedValue, nil, visited)
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
func (t *T) isNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}

// Returns a string representation of a Value that is sanitized based on what
// we are allowed to see. Private fields can not be exposed via a call to
// Interface() and trying will cause a panic, so we must use String()
// in that case.
func stringValue(v reflect.Value) string {
	if !v.IsValid() {
		return "<invalid reflect.Value>"
	}
	if v.CanInterface() {
		return fmt.Sprintf("%#v", v.Interface())
	}
	return v.String()
}

// Deep comparison. This is based on golang's reflect.DeepEqual functionality.
func (t *T) deepEqual(
	desc string, have, want reflect.Value, ignores []string,
	visited map[uintptr]*visitedNode,
) (diffs []string) {
	for _, ignore := range ignores {
		if desc == ignore {
			return nil
		}
	}
	if !want.IsValid() && !have.IsValid() {
		return nil
	} else if !want.IsValid() && have.IsValid() {
		prefix := ""
		if desc != "" {
			prefix = desc + ": "
		}
		return []string{
			fmt.Sprintf("%swanted an invalid or nil object, have valid object.", prefix),
		}
	} else if want.IsValid() && !have.IsValid() {
		prefix := ""
		if desc != "" {
			prefix = desc + ": "
		}
		return []string{
			fmt.Sprintf("%swanted a valid, non nil object, have invalid or nil object.", prefix),
		}
	} else if want.Type() != have.Type() {
		prefix := ""
		if desc != "" {
			prefix = desc + ": "
		}
		return []string{fmt.Sprintf(
			"%sNot the same type have: '%s', want: '%s'",
			prefix, have.Type(), want.Type())}
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
		prefix := ""
		if desc != "" {
			prefix = desc + ": "
		}
		if want.IsNil() && !have.IsNil() {
			diffs = append(diffs, fmt.Sprintf("%snot equal.", prefix))
			diffs = append(diffs, fmt.Sprintf("  have: %s", stringValue(have)))
			diffs = append(diffs, "  want: nil")
			return true
		} else if !want.IsNil() && have.IsNil() {
			diffs = append(diffs, fmt.Sprintf("%snot equal.", prefix))
			diffs = append(diffs, "  have: nil")
			diffs = append(diffs, fmt.Sprintf("  want: %s", stringValue(want)))
			return true
		}
		return false
	}

	// Checks to see that the lengths of both objects are equal.
	checkLen := func() bool {
		if want.Len() != have.Len() {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			diffs = append(diffs, fmt.Sprintf(
				"%s(len(have): %d, len(want): %d)",
				prefix, have.Len(), want.Len()))
			diffs = append(diffs, fmt.Sprintf("  have: %s", stringValue(have)))
			diffs = append(diffs, fmt.Sprintf("  want: %s", stringValue(want)))
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
					have.Index(i), want.Index(i), ignores, visited)
				diffs = append(diffs, newdiffs...)
			}
		}

	case reflect.Chan:
		// Channels are complex to compare so we rely on the existing type
		// checks to assert correctness, and then we add an additional
		// capacity check to assert buffer size.
		hcap := have.Cap()
		wcap := want.Cap()
		if hcap != wcap {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			diffs = append(diffs, fmt.Sprintf(
				"%sCapacities differ:\n  have: %d\n  want: %d",
				prefix, hcap, wcap))
			return diffs
		}

	case reflect.Func:
		if have.Pointer() == want.Pointer() {
			return nil
		}
		if !checkNil() {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			diffs = append(diffs, fmt.Sprintf("%snot equal.", prefix))
			diffs = append(diffs, fmt.Sprintf("  have: %s", stringValue(have)))
			diffs = append(diffs, fmt.Sprintf("  want: %s", stringValue(want)))
		}

	case reflect.Interface:
		if have.IsNil() && want.IsNil() {
			return nil
		}
		if !checkNil() {
			newdiffs := t.deepEqual(
				desc, have.Elem(), want.Elem(), ignores, visited)
			diffs = append(diffs, newdiffs...)
		}

	case reflect.Map:
		if have.Pointer() == want.Pointer() && have.Pointer() != 0 {
			return nil
		}
		if !checkNil() {
			// Check that the keys are present in both maps.
			for _, k := range want.MapKeys() {
				var keyStr string
				if k.Kind() == reflect.String {
					keyStr = fmt.Sprintf("[%q]", k.String())
				} else {
					keyStr = fmt.Sprintf("[%v]", k.Interface())
				}
				keyDesc := keyStr
				if desc != "" {
					keyDesc = desc + keyStr
				}

				if !have.MapIndex(k).IsValid() {
					prefix := ""
					if desc != "" {
						prefix = desc + ": "
					}
					diffs = append(diffs, fmt.Sprintf(
						"%sExpected key %s is missing.", prefix, keyStr))
					diffs = append(diffs, "  have: not present")
					diffs = append(diffs, fmt.Sprintf("  want: %s",
						stringValue(want.MapIndex(k))))
					continue
				}
				newdiffs := t.deepEqual(
					keyDesc,
					have.MapIndex(k), want.MapIndex(k), ignores, visited)
				diffs = append(diffs, newdiffs...)
			}
			for _, k := range have.MapKeys() {
				if !want.MapIndex(k).IsValid() {
					var keyStr string
					if k.Kind() == reflect.String {
						keyStr = fmt.Sprintf("[%q]", k.String())
					} else {
						keyStr = fmt.Sprintf("[%v]", k.Interface())
					}
					prefix := ""
					if desc != "" {
						prefix = desc + ": "
					}
					diffs = append(diffs, fmt.Sprintf(
						"%sUnexpected key %s.", prefix, keyStr))
					diffs = append(diffs,
						fmt.Sprintf("  have: %s", stringValue(have.MapIndex(k))))
					diffs = append(diffs, "  want: not present")
				}
			}
		}

	case reflect.Ptr:
		if have.Pointer() == want.Pointer() {
			return nil
		}
		if !checkNil() {
			newdiffs := t.deepEqual(
				desc, have.Elem(), want.Elem(), ignores, visited)
			diffs = append(diffs, newdiffs...)
		}

	case reflect.Slice:
		if !checkNil() && !checkLen() {
			for i := 0; i < want.Len(); i++ {
				newdiffs := t.deepEqual(
					fmt.Sprintf("%s[%d]", desc, i),
					have.Index(i), want.Index(i), ignores, visited)
				diffs = append(diffs, newdiffs...)
			}
		}

	case reflect.String:
		// We know the underlying type is a string so calling String()
		// will return the underlying value. Trying to call Interface()
		// and assert to a string will panic.
		hstr := have.String()
		wstr := want.String()
		if hstr == wstr {
			// Cheap equality test passed, no need to continue.
			break
		}
		hrunes := []rune(hstr)
		wrunes := []rune(wstr)
		if len(hrunes) != len(wrunes) {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%slen(have) %d != len(want) %d.",
					prefix, len(hrunes), len(wrunes)),
				fmt.Sprintf("  have: %#v", hstr),
				fmt.Sprintf("  want: %#v", wstr),
			}
		}
		for i, r := range hrunes {
			if r != wrunes[i] {
				prefix := ""
				if desc != "" {
					prefix = desc + ": "
				}
				return []string{
					fmt.Sprintf("%sdifference at rune %d.", prefix, i),
					fmt.Sprintf("  have: %#v", hstr),
					fmt.Sprintf("  want: %#v", wstr),
				}
			}
		}

	case reflect.Struct:
		for i, n := 0, want.NumField(); i < n; i++ {
			name := want.Type().Field(i).Name
			// Make sure that we don't print a strange error if the
			// first object given to us is a struct.
			if desc == "" {
				newdiffs := t.deepEqual(
					name, have.Field(i), want.Field(i), ignores, visited)
				diffs = append(diffs, newdiffs...)
			} else {
				newdiffs := t.deepEqual(
					fmt.Sprintf("%s.%s", desc, name),
					have.Field(i), want.Field(i), ignores, visited)
				diffs = append(diffs, newdiffs...)
			}
		}

	case reflect.Uintptr:
		// Uintptr's work like UnsafePointers. We can't evaluate them or
		// do much with them so we have to cast them into a number and
		// compare them that way.
		havePtr := have.Uint()
		wantPtr := want.Uint()
		if havePtr != wantPtr {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal.", prefix),
				fmt.Sprintf("  have: %#v", havePtr),
				fmt.Sprintf("  want: %#v", wantPtr),
			}
		}

	case reflect.UnsafePointer:
		// Unsafe pointers can cause us problems as they fall ill of the
		// Interface() restrictions. As such we have to special case them
		// and cast them as integers.
		havePtr := have.Pointer()
		wantPtr := want.Pointer()
		if havePtr != wantPtr {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal.", prefix),
				fmt.Sprintf("  have: %#v", havePtr),
				fmt.Sprintf("  want: %#v", wantPtr),
			}
		}

	case reflect.Bool:
		haveBool := have.Bool()
		wantBool := want.Bool()
		if haveBool != wantBool {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal.", prefix),
				fmt.Sprintf("  have: bool(%t)", haveBool),
				fmt.Sprintf("  want: bool(%t)", wantBool),
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Basic integer types.
		haveInt := have.Int()
		wantInt := want.Int()
		if haveInt != wantInt {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal", prefix),
				fmt.Sprintf("  have: %s(%d)", have.Type(), haveInt),
				fmt.Sprintf("  want: %s(%d)", want.Type(), wantInt),
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Basic unsigned integer types.
		haveUint := have.Uint()
		wantUint := want.Uint()
		if haveUint != wantUint {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal", prefix),
				fmt.Sprintf("  have: %s(%d)", have.Type(), haveUint),
				fmt.Sprintf("  want: %s(%d)", want.Type(), wantUint),
			}
		}

	case reflect.Float32, reflect.Float64:
		// Float types.
		haveFloat := have.Float()
		wantFloat := want.Float()
		if haveFloat != wantFloat {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal", prefix),
				fmt.Sprintf("  have: %s(%f)", have.Type(), haveFloat),
				fmt.Sprintf("  want: %s(%f)", want.Type(), wantFloat),
			}
		}

	case reflect.Complex64, reflect.Complex128:
		// Complex types.
		haveComplex := have.Complex()
		wantComplex := want.Complex()
		if haveComplex != wantComplex {
			prefix := ""
			if desc != "" {
				prefix = desc + ": "
			}
			return []string{
				fmt.Sprintf("%snot equal", prefix),
				fmt.Sprintf("  have: %s(%v)", have.Type(), haveComplex),
				fmt.Sprintf("  want: %s(%v)", want.Type(), wantComplex),
			}
		}
	}

	return diffs
}
