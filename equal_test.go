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
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"unicode"
)

type testEqualCustomStruct struct {
	Field1 string
	Field2 string
}

type testEqualCustomNode struct {
	value string
	next  *testEqualCustomNode
}

func TestT_EqualAndNotEqual(t *testing.T) {
	t.Parallel()

	// Describes a given object. This is used to make it easier to read
	// the output of a test failure.
	describe := func(i interface{}) string {
		return fmt.Sprintf("%T(%#v)", i, i)
	}

	// This is a list of valid runes that can be in a string. We use this
	// when generating random  strings later.
	validRunes := make([]rune, 0, 10000)
	for _, r16 := range unicode.Latin.R16 {
		for i := r16.Lo; i < r16.Hi; i = i + r16.Stride {
			validRunes = append(validRunes, rune(i))
		}
	}
	for _, r32 := range unicode.Latin.R32 {
		for i := r32.Lo; i < r32.Hi; i = i + r32.Stride {
			validRunes = append(validRunes, rune(i))
		}
	}

	// Runs a test against sets of like objects. All 'same' objects will
	// be evaluated to ensure that T.Equal() passes, and T.NotEqual() fails.
	// All same objects will be run against all diff objects to ensure that
	// T.Equal() fails and T.NotEqual passes.
	runTest := func(same, diff []interface{}) {
		m, T := testSetup()
		for _, s1 := range same {
			// Ensure equality for all of the "same" objects.
			for _, s2 := range same {
				description := fmt.Sprintf("Equal() failed for:\n%s\n%s",
					describe(s1), describe(s2))
				m.CheckPass(t, func() { T.Equal(s1, s2) }, description)
				description = fmt.Sprintf("Equal() failed for:\n%s\n%s",
					describe(s2), describe(s1))
				m.CheckPass(t, func() { T.Equal(s2, s1) }, description)
				description = fmt.Sprintf("NotEqual() passed for:\n%s\n%s",
					describe(s1), describe(s2))
				m.CheckFail(t, func() { T.NotEqual(s1, s2) }, description)
				description = fmt.Sprintf("NotEqual() passed for:\n%s\n%s",
					describe(s2), describe(s1))
				m.CheckFail(t, func() { T.NotEqual(s2, s1) }, description)
			}

			// Ensure non-equality for all of the "diff" objects.
			for _, d := range diff {
				description := fmt.Sprintf("Equal() passed for:\n%s\n%s",
					describe(s1), describe(d))
				m.CheckFail(t, func() { T.Equal(s1, d) }, description)
				description = fmt.Sprintf("Equal() passed for:\n%s\n%s",
					describe(d), describe(s1))
				m.CheckFail(t, func() { T.Equal(d, s1) }, description)
				description = fmt.Sprintf("NotEqual() failed for:\n%s\n%s",
					describe(s1), describe(d))
				m.CheckPass(t, func() { T.NotEqual(s1, d) }, description)
				description = fmt.Sprintf("NotEqual() failed for:\n%s\n%s",
					describe(d), describe(s1))
				m.CheckPass(t, func() { T.NotEqual(d, s1) }, description)
			}
		}
	}

	// To start we run the test against basic types. The equality of these
	// should be really obvious. We iterate through each so we get a good
	// sampling of each type.
	for loop := 0; loop < 100; loop++ {
		// Integers
		ir := rand.Int63()
		runTest(
			[]interface{}{int(ir), int(ir)},
			[]interface{}{int(ir + 1), uint(ir)})
		runTest(
			[]interface{}{int8(ir), int8(ir)},
			[]interface{}{int8(ir + 1), uint8(ir)})
		runTest(
			[]interface{}{int16(ir), int16(ir)},
			[]interface{}{int16(ir + 1), uint16(ir)})
		runTest(
			[]interface{}{int32(ir), int32(ir)},
			[]interface{}{int32(ir + 1), uint32(ir)})
		runTest(
			[]interface{}{int64(ir), int64(ir)},
			[]interface{}{int64(ir + 1), uint64(ir)})
		runTest(
			[]interface{}{uint(ir), uint(ir)},
			[]interface{}{uint(ir + 1), int(ir)})
		runTest(
			[]interface{}{uint8(ir), uint8(ir)},
			[]interface{}{uint8(ir + 1), int8(ir)})
		runTest(
			[]interface{}{uint16(ir), uint16(ir)},
			[]interface{}{uint16(ir + 1), int16(ir)})
		runTest(
			[]interface{}{uint32(ir), uint32(ir)},
			[]interface{}{uint32(ir + 1), int32(ir)})
		runTest(
			[]interface{}{uint64(ir), uint64(ir)},
			[]interface{}{uint64(ir + 1), int64(ir)})
		runTest(
			[]interface{}{uintptr(ir), uintptr(ir)},
			[]interface{}{uintptr(ir + 1), uint64(ir)})

		// Floating Point.
		fr := rand.NormFloat64()
		runTest(
			[]interface{}{float32(fr), float32(fr)},
			[]interface{}{float32(fr + 1), float64(fr)})
		runTest(
			[]interface{}{float64(fr), float64(fr)},
			[]interface{}{float64(fr + 1), float32(fr)})

		// Strings
		rarray := make([]rune, 100)
		for i := range rarray {
			rarray[i] = validRunes[rand.Intn(len(validRunes))]
		}
		s1 := string(rarray)
		s2 := string(rarray)
		d1 := string(rarray) + "1"
		d2 := string(rarray[0:50])
		if rarray[10] == 'A' {
			rarray[10] = 'B'
		} else {
			rarray[10] = 'A'
		}
		d3 := string(rarray)
		runTest(
			[]interface{}{s1, s2},
			[]interface{}{"", d1, d2, d3})
	}

	// Boolean (no need to test with 100 values.)
	runTest(
		[]interface{}{true},
		[]interface{}{false})

	// Basic nil value check.
	var nilOsFilePtr *os.File
	var nilStringPtr *string
	var nilInterface interface{} = nilOsFilePtr
	var str string
	runTest(
		[]interface{}{nil, nilInterface, nilOsFilePtr, nilStringPtr},
		[]interface{}{new(os.File), &str})

	// Array check.
	sIntArray1 := [5]int{0, 1, 2, 3, 4}
	sIntArray2 := [5]int{0, 1, 2, 3, 4}
	dIntArray1 := [4]int{0, 1, 2, 3}
	dIntArray2 := [6]int{0, 1, 2, 3, 4, 5}
	dIntArray3 := [5]int{0, 1, 2, 3, 5}
	runTest(
		[]interface{}{sIntArray1, sIntArray2},
		[]interface{}{dIntArray1, dIntArray2, dIntArray3})

	// Slice check.
	sIntSlice1 := []int{0, 1, 2, 3, 4}
	sIntSlice2 := []int{0, 1, 2, 3, 4}
	dIntSlice1 := []int{0, 1, 2, 3}
	dIntSlice2 := []int{0, 1, 2, 3, 4, 5}
	dIntSlice3 := []int{0, 1, 2, 3, 5}
	runTest(
		[]interface{}{sIntSlice1, sIntSlice2},
		[]interface{}{dIntSlice1, dIntSlice2, dIntSlice3})

	// Map check
	sIntMap1 := map[string]int{"a": 1, "b": 2, "c": 3}
	sIntMap2 := map[string]int{"a": 1, "b": 2, "c": 3}
	dIntMap1 := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	dIntMap2 := map[string]int{"a": 1, "b": 2}
	dIntMap3 := map[string]int{"a": 1, "b": 2, "c": 4}
	runTest(
		[]interface{}{sIntMap1, sIntMap2},
		[]interface{}{dIntMap1, dIntMap2, dIntMap3})

	// Pointer check.
	sString1 := "a"
	sString2 := "a"
	dString1 := "b"
	runTest(
		[]interface{}{&sString1, &sString2},
		[]interface{}{&dString1})

	// Interface check
	var sInterface1 interface{} = sString1
	var sInterface2 interface{} = sString2
	var dInterface1 interface{} = dString1
	var dInterface2 interface{}
	runTest(
		[]interface{}{&sInterface1, &sInterface2},
		[]interface{}{&dInterface1, &dInterface2})

	// Functions
	// Note that the only check that can be done here is a pure pointer
	// check so we have limited options.
	sFunc1 := func() {}
	sFunc2 := sFunc1
	var dFunc1 func()
	runTest(
		[]interface{}{sFunc1, sFunc2},
		[]interface{}{dFunc1})

	// Structure types.
	sCust1 := testEqualCustomStruct{Field1: "a", Field2: "b"}
	sCust2 := testEqualCustomStruct{Field1: "a", Field2: "b"}
	dCust1 := testEqualCustomStruct{Field1: "a", Field2: "c"}
	runTest(
		[]interface{}{sCust1, sCust2},
		[]interface{}{dCust1})

	// Structures in a slice.
	sCustSlice1 := []testEqualCustomStruct{sCust1, sCust2}
	sCustSlice2 := []testEqualCustomStruct{sCust1, sCust2}
	dCustSlice1 := []testEqualCustomStruct{sCust1}
	dCustSlice2 := []testEqualCustomStruct{sCust2}
	dCustSlice3 := []testEqualCustomStruct{sCust1, dCust1}
	runTest(
		[]interface{}{sCustSlice1, sCustSlice2},
		[]interface{}{dCustSlice1, dCustSlice2, dCustSlice3})

	// Invalid structures.
	sValid1 := reflect.Value{}
	sValid2 := reflect.Value{}
	dValid1 := reflect.ValueOf(TestT_EqualAndNotEqual)
	runTest(
		[]interface{}{sValid1, sValid2},
		[]interface{}{dValid1})

	// Circular linked list.
	sCircle1 := &testEqualCustomNode{value: "a"}
	sCircle1.next = sCircle1
	sCircle2 := &testEqualCustomNode{value: "a"}
	sCircle2.next = sCircle2
	dCircle1 := &testEqualCustomNode{value: "a"}
	dCircle1.next = &testEqualCustomNode{value: "b", next: dCircle1}
	runTest(
		[]interface{}{sCircle1, sCircle2},
		[]interface{}{dCircle1})

	// Channel
	sChan1 := make(chan bool, 10)
	sChan2 := make(chan bool, 10)
	dChan1 := make(chan bool, 1000)
	dChan2 := make(chan bool, 100)
	runTest(
		[]interface{}{sChan1, sChan2},
		[]interface{}{dChan1, dChan2})

	// Lastly we verify that the error messages get prefixed.
	m, T := testSetup()
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() { T.Equal(1, 2, "prefix") })
	if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("Prefix was not prepended to the error: %s", msg)
	}
	msg = ""
	m.CheckFail(t, func() { T.NotEqual(1, 1, "prefix") })
	if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("Prefix was not prepended to the error: %s", msg)
	}
}

/*

func TestT_Equal(t *testing.T) {
	t.Parallel()
	m := &mockT{}
	T := NewT(m)

	ct1 := testEqualCustomStruct{Field1: "s", Field2: "t"}
	ct2 := testEqualCustomStruct{Field1: "s", Field2: "t"}
	ct3 := testEqualCustomStruct{Field1: "s", Field2: "o"}
	ctSlice1 := []testEqualCustomStruct{ct1, ct2}
	ctSlice2 := []testEqualCustomStruct{ct1, ct2}
	ctSlice3 := []testEqualCustomStruct{ct1, ct3}
	f1 := func() {}
	var f2 func()
	var i1 interface{} = ct1
	var i2 interface{} = ct1
	var i3 interface{}
	mArray1 := [3]*mockT{m, m, m}
	mArray2 := [3]*mockT{m, m, m}
	mArray3 := [3]*mockT{m}
	mSlice1 := []*mockT{m, m, m}
	mSlice2 := []*mockT{m, m, m}
	mSlice3 := []*mockT{m, m}
	myStr1 := testEqualCustomType("one")
	myStr2 := testEqualCustomType("one")
	var nilPtr *mockT
	strMap1 := map[string]int{"A": 1, "B": 2, "C": 3}
	strMap2 := map[string]int{"C": 3, "B": 2, "A": 1}
	strMap3 := map[string]int{"C": 3, "B": 2, "A": -1}
	strMap4 := map[string]int{"C": 3, "B": 2, "D": 3}
	strSlice1 := []string{"A", "B", "C"}
	strSlice2 := []string{"A", "B", "C"}
	strSlice3 := []string{"X", "B", "C"}

	// These tests should all succeed.



	m.CheckPass(t, func() { T.Equal(ct1, ct2) })
	m.CheckPass(t, func() { T.Equal(ctSlice1, ctSlice2) })
	m.CheckPass(t, func() { T.Equal(&i1, &i2) })
	m.CheckPass(t, func() { T.Equal(f1, f1) })
	m.CheckPass(t, func() { T.Equal(m, m) })
	m.CheckPass(t, func() { T.Equal(mArray1, mArray2) })
	m.CheckPass(t, func() { T.Equal(mSlice1, mSlice2) })
	m.CheckPass(t, func() { T.Equal(myStr1, myStr2) })
	m.CheckPass(t, func() { T.Equal(nilPtr, nil) })
	m.CheckPass(t, func() { T.Equal(nil, nilPtr) })
	m.CheckPass(t, func() { T.Equal(strMap1, strMap2) })
	m.CheckPass(t, func() { T.Equal(strSlice1, strSlice2) })

	// Expected failure conditions.
	m.CheckFail(t, func() { T.Equal(&mockT{}, "A") })
	m.CheckFail(t, func() { T.Equal(&mockT{}, nil) })
	m.CheckFail(t, func() { T.Equal(nil, &mockT{}) })
	m.CheckFail(t, func() { T.Equal([]*mockT{nil}, []*mockT{&mockT{}}) })
	m.CheckFail(t, func() { T.Equal([]*mockT{&mockT{}}, []*mockT{nil}) })
	m.CheckFail(t, func() { T.Equal(false, true) })
	m.CheckFail(t, func() { T.Equal(int(2), int(1)) })
	m.CheckFail(t, func() { T.Equal(int8(2), int8(1)) })
	m.CheckFail(t, func() { T.Equal(int16(2), int16(1)) })
	m.CheckFail(t, func() { T.Equal(int32(2), int32(1)) })
	m.CheckFail(t, func() { T.Equal(int64(2), int64(1)) })
	m.CheckFail(t, func() { T.Equal(uint(2), uint(1)) })
	m.CheckFail(t, func() { T.Equal(uint8(2), uint8(1)) })
	m.CheckFail(t, func() { T.Equal(uint16(2), uint16(1)) })
	m.CheckFail(t, func() { T.Equal(uint32(2), uint32(1)) })
	m.CheckFail(t, func() { T.Equal(uint64(2), uint64(1)) })
	m.CheckFail(t, func() { T.Equal(uintptr(2), uintptr(1)) })
	m.CheckFail(t, func() { T.Equal(float32(2), float32(1)) })
	m.CheckFail(t, func() { T.Equal(float64(2), float64(1)) })
	m.CheckFail(t, func() { T.Equal("2", "1") })
	m.CheckFail(t, func() { T.Equal("22", "1") })
	m.CheckFail(t, func() { T.Equal(ct1, ct3) })
	m.CheckFail(t, func() { T.Equal(ctSlice1, ctSlice3) })
	m.CheckFail(t, func() { T.Equal(&i1, &i3) })
	m.CheckFail(t, func() { T.Equal(f1, f2) })
	m.CheckFail(t, func() { T.Equal(mArray1, mArray3) })
	m.CheckFail(t, func() { T.Equal(mSlice1, mSlice3) })
	m.CheckFail(t, func() { T.Equal(strMap1, strMap3) })
	m.CheckFail(t, func() { T.Equal(strMap1, strMap4) })
	m.CheckFail(t, func() { T.Equal(strSlice1, strSlice3) })

	// Check the returned error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() { T.Equal(1, 2, "prefix") })
    if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("Error message did not contain the prefix: '''%s'''", msg)
	}
}

func TestT_NotEqual(t *testing.T) {
	t.Parallel()
	m := &mockT{}
	T := NewT(m)

	ct1 := testEqualCustomStruct{Field1: "s", Field2: "t"}
	ct2 := testEqualCustomStruct{Field1: "s", Field2: "t"}
	ct3 := testEqualCustomStruct{Field1: "s", Field2: "o"}
	ctSlice1 := []testEqualCustomStruct{ct1, ct2}
	ctSlice2 := []testEqualCustomStruct{ct1, ct2}
	ctSlice3 := []testEqualCustomStruct{ct1, ct3}
	f1 := func() {}
	var f2 func()
	var i1 interface{} = ct1
	var i2 interface{} = ct2
	var i3 interface{}
	mArray1 := [3]*mockT{m, m, m}
	mArray2 := [3]*mockT{m, m, m}
	mArray3 := [3]*mockT{m}
	mSlice1 := []*mockT{m, m, m}
	mSlice2 := []*mockT{m, m, m}
	mSlice3 := []*mockT{m, m}
	var nilPtr *mockT
	strMap1 := map[string]int{"A": 1, "B": 2, "C": 3}
	strMap2 := map[string]int{"C": 3, "B": 2, "A": 1}
	strMap3 := map[string]int{"C": 3, "B": 2, "A": -1}
	strMap4 := map[string]int{"C": 3, "B": 2, "D": 3}
	strSlice1 := []string{"A", "B", "C"}
	strSlice2 := []string{"A", "B", "C"}
	strSlice3 := []string{"X", "B", "C"}

	// CheckFailure cases.
	m.CheckFail(t, func() { T.NotEqual(nil, nil) })
	m.CheckFail(t, func() { T.NotEqual(true, true) })
	m.CheckFail(t, func() { T.NotEqual(int(1), int(1)) })
	m.CheckFail(t, func() { T.NotEqual(int8(1), int8(1)) })
	m.CheckFail(t, func() { T.NotEqual(int16(1), int16(1)) })
	m.CheckFail(t, func() { T.NotEqual(int32(1), int32(1)) })
	m.CheckFail(t, func() { T.NotEqual(int64(1), int64(1)) })
	m.CheckFail(t, func() { T.NotEqual(uint(1), uint(1)) })
	m.CheckFail(t, func() { T.NotEqual(uint8(1), uint8(1)) })
	m.CheckFail(t, func() { T.NotEqual(uint16(1), uint16(1)) })
	m.CheckFail(t, func() { T.NotEqual(uint32(1), uint32(1)) })
	m.CheckFail(t, func() { T.NotEqual(uint64(1), uint64(1)) })
	m.CheckFail(t, func() { T.NotEqual(uintptr(1), uintptr(1)) })
	m.CheckFail(t, func() { T.NotEqual(float32(1), float32(1)) })
	m.CheckFail(t, func() { T.NotEqual(float64(1), float64(1)) })
	m.CheckFail(t, func() { T.NotEqual("1", "1") })
	m.CheckFail(t, func() { T.NotEqual(ct1, ct2) })
	m.CheckFail(t, func() { T.NotEqual(ctSlice1, ctSlice2) })
	m.CheckFail(t, func() { T.NotEqual(&i1, &i2) })
	m.CheckFail(t, func() { T.NotEqual(f1, f1) })
	m.CheckFail(t, func() { T.NotEqual(m, m) })
	m.CheckFail(t, func() { T.NotEqual(mArray1, mArray2) })
	m.CheckFail(t, func() { T.NotEqual(mSlice1, mSlice2) })
	m.CheckFail(t, func() { T.NotEqual(nilPtr, nil) })
	m.CheckFail(t, func() { T.NotEqual(nil, nilPtr) })
	m.CheckFail(t, func() { T.NotEqual(strMap1, strMap2) })
	m.CheckFail(t, func() { T.NotEqual(strSlice1, strSlice2) })

	// Non failure cases.
	m.CheckPass(t, func() { T.NotEqual(&mockT{}, "A") })
	m.CheckPass(t, func() { T.NotEqual(&mockT{}, nil) })
	m.CheckPass(t, func() { T.NotEqual(nil, &mockT{}) })
	m.CheckPass(t, func() { T.NotEqual([]*mockT{nil}, []*mockT{&mockT{}}) })
	m.CheckPass(t, func() { T.NotEqual([]*mockT{&mockT{}}, []*mockT{nil}) })
	m.CheckPass(t, func() { T.NotEqual(false, true) })
	m.CheckPass(t, func() { T.NotEqual(int(2), int(1)) })
	m.CheckPass(t, func() { T.NotEqual(int8(2), int8(1)) })
	m.CheckPass(t, func() { T.NotEqual(int16(2), int16(1)) })
	m.CheckPass(t, func() { T.NotEqual(int32(2), int32(1)) })
	m.CheckPass(t, func() { T.NotEqual(int64(2), int64(1)) })
	m.CheckPass(t, func() { T.NotEqual(uint(2), uint(1)) })
	m.CheckPass(t, func() { T.NotEqual(uint8(2), uint8(1)) })
	m.CheckPass(t, func() { T.NotEqual(uint16(2), uint16(1)) })
	m.CheckPass(t, func() { T.NotEqual(uint32(2), uint32(1)) })
	m.CheckPass(t, func() { T.NotEqual(uint64(2), uint64(1)) })
	m.CheckPass(t, func() { T.NotEqual(uintptr(2), uintptr(1)) })
	m.CheckPass(t, func() { T.NotEqual(float32(2), float32(1)) })
	m.CheckPass(t, func() { T.NotEqual(float64(2), float64(1)) })
	m.CheckPass(t, func() { T.NotEqual("2", "1") })
	m.CheckPass(t, func() { T.NotEqual("22", "1") })
	m.CheckPass(t, func() { T.NotEqual(ct1, ct3) })
	m.CheckPass(t, func() { T.NotEqual(ctSlice1, ctSlice3) })
	m.CheckPass(t, func() { T.NotEqual(&i1, &i3) })
	m.CheckPass(t, func() { T.NotEqual(f1, f2) })
	m.CheckPass(t, func() { T.NotEqual(mArray1, mArray3) })
	m.CheckPass(t, func() { T.NotEqual(mSlice1, mSlice3) })
	m.CheckPass(t, func() { T.NotEqual(strMap1, strMap3) })
	m.CheckPass(t, func() { T.NotEqual(strMap1, strMap4) })
	m.CheckPass(t, func() { T.NotEqual(strSlice1, strSlice3) })

	// Check the returned error message.
	msg := ""
	m.funcFatal = func(args ...interface{}) {
		msg = fmt.Sprint(args...)
	}
	m.CheckFail(t, func() { T.NotEqual(1, 1, "prefix") })
    if !strings.HasPrefix(msg, "prefix: ") {
		t.Fatalf("Error message did not contain the prefix: '''%s'''", msg)
	}
}

*/
