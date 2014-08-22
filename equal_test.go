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
	"testing"
)

type testEqualCutomType string

func TestTestLibEqual(t *testing.T) {
	m := &mockT{}
	T := NewTestLib(m)

	var nilPtr *mockT
	strSlice1 := []string{"A", "B", "C"}
	strSlice2 := []string{"A", "B", "C"}
	strSlice3 := []string{"X", "B", "C"}
	strMap1 := map[string]int{"A": 1, "B": 2, "C": 3}
	strMap2 := map[string]int{"C": 3, "B": 2, "A": 1}
	strMap3 := map[string]int{"C": 3, "B": 2, "A": -1}
	myStr1 := testEqualCutomType("one")
	myStr2 := testEqualCutomType("one")
	mSlice1 := []*mockT{m, m, m}
	mSlice2 := []*mockT{m, m, m}
	mSlice3 := []*mockT{m, m}
	mArray1 := [3]*mockT{m, m, m}
	mArray2 := [3]*mockT{m, m, m}
	mArray3 := [3]*mockT{m}
	f1 := func() {}
	var f2 func()
	var i1 interface{} = mSlice1
	var i2 interface{} = mSlice2
	var i3 interface{}

	// These tests should all succeed.
	m.CheckPass(t, func() { T.Equal(nil, nil) })
	m.CheckPass(t, func() { T.Equal(true, true) })
	m.CheckPass(t, func() { T.Equal(int(1), int(1)) })
	m.CheckPass(t, func() { T.Equal(int8(1), int8(1)) })
	m.CheckPass(t, func() { T.Equal(int16(1), int16(1)) })
	m.CheckPass(t, func() { T.Equal(int32(1), int32(1)) })
	m.CheckPass(t, func() { T.Equal(int64(1), int64(1)) })
	m.CheckPass(t, func() { T.Equal(uint(1), uint(1)) })
	m.CheckPass(t, func() { T.Equal(uint8(1), uint8(1)) })
	m.CheckPass(t, func() { T.Equal(uint16(1), uint16(1)) })
	m.CheckPass(t, func() { T.Equal(uint32(1), uint32(1)) })
	m.CheckPass(t, func() { T.Equal(uint64(1), uint64(1)) })
	m.CheckPass(t, func() { T.Equal(uintptr(1), uintptr(1)) })
	m.CheckPass(t, func() { T.Equal(float32(1), float32(1)) })
	m.CheckPass(t, func() { T.Equal(float64(1), float64(1)) })
	m.CheckPass(t, func() { T.Equal("1", "1") })
	m.CheckPass(t, func() { T.Equal(nilPtr, nil) })
	m.CheckPass(t, func() { T.Equal(strSlice1, strSlice2) })
	m.CheckPass(t, func() { T.Equal(strMap1, strMap2) })
	m.CheckPass(t, func() { T.Equal(myStr1, myStr2) })
	m.CheckPass(t, func() { T.Equal(m, m) })
	m.CheckPass(t, func() { T.Equal(mSlice1, mSlice2) })
	m.CheckPass(t, func() { T.Equal(mArray1, mArray2) })
	m.CheckPass(t, func() { T.Equal(f1, f1) })
	m.CheckPass(t, func() { T.Equal(&i1, &i2) })

	// Expected failure conditions.
	m.CheckFail(t, func() { T.Equal(&mockT{}, "A") })
	m.CheckFail(t, func() { T.Equal(&mockT{}, nil) })
	m.CheckFail(t, func() { T.Equal(nil, &mockT{}) })
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
	m.CheckFail(t, func() { T.Equal(mSlice1, mSlice3) })
	m.CheckFail(t, func() { T.Equal(mArray1, mArray3) })
	m.CheckFail(t, func() { T.Equal(f1, f2) })
	m.CheckFail(t, func() { T.Equal(&i1, &i3) })

	m.CheckFail(t, func() { T.Equal(strSlice1, strSlice3) })
	m.CheckFail(t, func() { T.Equal(strMap1, strMap3) })
}

func TestTestLibNotEqual(t *testing.T) {
	m := &mockT{}
	T := NewTestLib(m)

	var nilPtr *mockT
	strSlice1 := []string{"A", "B", "C"}
	strSlice2 := []string{"A", "B", "C"}
	strSlice3 := []string{"X", "B", "C"}
	strMap1 := map[string]int{"A": 1, "B": 2, "C": 3}
	strMap2 := map[string]int{"C": 3, "B": 2, "A": 1}
	strMap3 := map[string]int{"C": 3, "B": 2, "A": -1}
	mSlice1 := []*mockT{m, m, m}
	mSlice2 := []*mockT{m, m, m}
	mSlice3 := []*mockT{m, m}

	mArray1 := [3]*mockT{m, m, m}
	mArray2 := [3]*mockT{m, m, m}
	mArray3 := [3]*mockT{m}
	f1 := func() {}
	var f2 func()
	var i1 interface{} = mSlice1
	var i2 interface{} = mSlice2
	var i3 interface{}

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
	m.CheckFail(t, func() { T.NotEqual(nilPtr, nil) })
	m.CheckFail(t, func() { T.NotEqual(strSlice1, strSlice2) })
	m.CheckFail(t, func() { T.NotEqual(strMap1, strMap2) })
	m.CheckFail(t, func() { T.NotEqual(m, m) })
	m.CheckFail(t, func() { T.NotEqual(mSlice1, mSlice2) })
	m.CheckFail(t, func() { T.NotEqual(mArray1, mArray2) })
	m.CheckFail(t, func() { T.NotEqual(f1, f1) })
	m.CheckFail(t, func() { T.NotEqual(&i1, &i2) })

	// Non failure cases.
	m.CheckPass(t, func() { T.NotEqual(&mockT{}, "A") })
	m.CheckPass(t, func() { T.NotEqual(&mockT{}, nil) })
	m.CheckPass(t, func() { T.NotEqual(nil, &mockT{}) })
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
	m.CheckPass(t, func() { T.NotEqual(mSlice1, mSlice3) })
	m.CheckPass(t, func() { T.NotEqual(mArray1, mArray3) })
	m.CheckPass(t, func() { T.NotEqual(f1, f2) })
	m.CheckPass(t, func() { T.NotEqual(&i1, &i3) })
	m.CheckPass(t, func() { T.NotEqual(strSlice1, strSlice3) })
	m.CheckPass(t, func() { T.NotEqual(strMap1, strMap3) })
}
