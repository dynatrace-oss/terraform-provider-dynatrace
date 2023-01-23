/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package opt_test

import (
	"math/rand"
	"testing"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/google/uuid"
)

// TestInt validates `opt.Int(*int)` and `opt.NewInt(int)`
func TestInt(t *testing.T) {

	var want int

	if got := opt.Int(nil); got != want {
		t.Errorf("Int(nil) failed, got: %v, want: %v.", got, want)
	}

	want = rand.Int()
	if got := opt.Int(&want); got != want {
		t.Errorf("Int(*int) failed, got: %v, want: %v.", got, want)
	}

	want = rand.Int()
	if got := *opt.NewInt(want); got != want {
		t.Errorf("NewInt(int) failed, got: %v, want: %v.", got, want)
	}
}

// TestString validates `opt.String(*string)` and `opt.NewString(string)`
func TestString(t *testing.T) {

	var want string

	if got := opt.String(nil); got != want {
		t.Errorf("String(nil) failed, got: %v, want: %v.", got, want)
	}

	want = uuid.New().String()
	if got := opt.String(&want); got != want {
		t.Errorf("String(*string) failed, got: %v, want: %v.", got, want)
	}

	want = uuid.New().String()
	if got := *opt.NewString(want); got != want {
		t.Errorf("NewString(string) failed, got: %v, want: %v.", got, want)
	}
}

// TestBool validates `opt.Bool(*bool)` and `opt.NewBool(bool)`
func TestBool(t *testing.T) {

	var want bool

	if got := opt.Bool(nil); got != want {
		t.Errorf("Bool(nil) failed, got: %v, want: %v.", got, want)
	}

	want = (rand.Intn(2) == 0)
	if got := opt.Bool(&want); got != want {
		t.Errorf("Bool(*bool) failed, got: %v, want: %v.", got, want)
	}

	want = (rand.Intn(2) == 0)
	if got := *opt.NewBool(want); got != want {
		t.Errorf("NewBool(bool) failed, got: %v, want: %v.", got, want)
	}
}

// TestInt8 validates `opt.Int8(*int8)` and `opt.NewInt8(int8)`
func TestInt8(t *testing.T) {

	var want int8

	if got := opt.Int8(nil); got != want {
		t.Errorf("Int8(nil) failed, got: %v, want: %v.", got, want)
	}

	want = int8(rand.Int())
	if got := opt.Int8(&want); got != want {
		t.Errorf("Int8(*int8) failed, got: %v, want: %v.", got, want)
	}

	want = int8(rand.Int())
	if got := *opt.NewInt8(want); got != want {
		t.Errorf("NewInt8(int8) failed, got: %v, want: %v.", got, want)
	}
}

// TestInt16 validates `opt.Int16(*int16)` and `opt.NewInt16(int16)`
func TestInt16(t *testing.T) {

	var want int16

	if got := opt.Int16(nil); got != want {
		t.Errorf("Int16(nil) failed, got: %v, want: %v.", got, want)
	}

	want = int16(rand.Int())
	if got := opt.Int16(&want); got != want {
		t.Errorf("Int16(*int16) failed, got: %v, want: %v.", got, want)
	}

	want = int16(rand.Int())
	if got := *opt.NewInt16(want); got != want {
		t.Errorf("NewInt16(int16) failed, got: %v, want: %v.", got, want)
	}
}

// TestInt32 validates `opt.Int32(*int32)` and `opt.NewInt32(int32)`
func TestInt32(t *testing.T) {

	var want int32

	if got := opt.Int32(nil); got != want {
		t.Errorf("Int32(nil) failed, got: %v, want: %v.", got, want)
	}

	want = int32(rand.Int())
	if got := opt.Int32(&want); got != want {
		t.Errorf("Int32(*int32) failed, got: %v, want: %v.", got, want)
	}

	want = int32(rand.Int())
	if got := *opt.NewInt32(want); got != want {
		t.Errorf("NewInt32(int32) failed, got: %v, want: %v.", got, want)
	}
}

// TestInt64 validates `opt.Int64(*int64)` and `opt.NewInt64(int64)`
func TestInt64(t *testing.T) {

	var want int64

	if got := opt.Int64(nil); got != want {
		t.Errorf("Int64(nil) failed, got: %v, want: %v.", got, want)
	}

	want = int64(rand.Int())
	if got := opt.Int64(&want); got != want {
		t.Errorf("Int64(*int64) failed, got: %v, want: %v.", got, want)
	}

	want = int64(rand.Int())
	if got := *opt.NewInt64(want); got != want {
		t.Errorf("NewInt64(int64) failed, got: %v, want: %v.", got, want)
	}
}

// TestUInt8 validates `opt.UInt8(*uint8)` and `opt.NewUInt8(uint8)`
func TestUInt8(t *testing.T) {

	var want uint8

	if got := opt.UInt8(nil); got != want {
		t.Errorf("UInt8(nil) failed, got: %v, want: %v.", got, want)
	}

	want = uint8(rand.Int())
	if got := opt.UInt8(&want); got != want {
		t.Errorf("UInt8(*uint8) failed, got: %v, want: %v.", got, want)
	}

	want = uint8(rand.Int())
	if got := *opt.NewUInt8(want); got != want {
		t.Errorf("NewUInt8(uint8) failed, got: %v, want: %v.", got, want)
	}
}

// TestUInt16 validates `opt.UInt16(*uint16)` and `opt.NewUInt16(uint16)`
func TestUInt16(t *testing.T) {

	var want uint16

	if got := opt.UInt16(nil); got != want {
		t.Errorf("UInt16(nil) failed, got: %v, want: %v.", got, want)
	}

	want = uint16(rand.Int())
	if got := opt.UInt16(&want); got != want {
		t.Errorf("UInt16(*uint16) failed, got: %v, want: %v.", got, want)
	}

	want = uint16(rand.Int())
	if got := *opt.NewUInt16(want); got != want {
		t.Errorf("NewUInt16(uint16) failed, got: %v, want: %v.", got, want)
	}
}

// TestUInt32 validates `opt.UInt32(*uint32)` and `opt.NewUInt32(uint32)`
func TestUInt32(t *testing.T) {

	var want uint32

	if got := opt.UInt32(nil); got != want {
		t.Errorf("UInt32(nil) failed, got: %v, want: %v.", got, want)
	}

	want = uint32(rand.Int())
	if got := opt.UInt32(&want); got != want {
		t.Errorf("UInt32(*uint32) failed, got: %v, want: %v.", got, want)
	}

	want = uint32(rand.Int())
	if got := *opt.NewUInt32(want); got != want {
		t.Errorf("NewUInt32(uint32) failed, got: %v, want: %v.", got, want)
	}
}

// TestUInt64 validates `opt.UInt64(*uint64)` and `opt.NewUInt64(uint64)`
func TestUInt64(t *testing.T) {

	var want uint64

	if got := opt.UInt64(nil); got != want {
		t.Errorf("UInt64(nil) failed, got: %v, want: %v.", got, want)
	}

	want = uint64(rand.Int())
	if got := opt.UInt64(&want); got != want {
		t.Errorf("UInt64(*uint64) failed, got: %v, want: %v.", got, want)
	}

	want = uint64(rand.Int())
	if got := *opt.NewUInt64(want); got != want {
		t.Errorf("NewUInt64(uint64) failed, got: %v, want: %v.", got, want)
	}
}

// TestFloat32 validates `opt.Float32(*float32)` and `opt.NewFloat32(float32)`
func TestFloat32(t *testing.T) {

	var want float32

	if got := opt.Float32(nil); got != want {
		t.Errorf("Float32(nil) failed, got: %v, want: %v.", got, want)
	}

	want = rand.Float32()
	if got := opt.Float32(&want); got != want {
		t.Errorf("Float32(*float32) failed, got: %v, want: %v.", got, want)
	}

	want = rand.Float32()
	if got := *opt.NewFloat32(want); got != want {
		t.Errorf("NewFloat32(float32) failed, got: %v, want: %v.", got, want)
	}
}

// TestFloat64 validates `opt.Float64(*float64)` and `opt.NewFloat64(float64)`
func TestFloat64(t *testing.T) {

	var want float64

	if got := opt.Float64(nil); got != want {
		t.Errorf("Float64(nil) failed, got: %v, want: %v.", got, want)
	}

	want = rand.Float64()
	if got := opt.Float64(&want); got != want {
		t.Errorf("Float64(*float64) failed, got: %v, want: %v.", got, want)
	}

	want = rand.Float64()
	if got := *opt.NewFloat64(want); got != want {
		t.Errorf("NewFloat64(float64) failed, got: %v, want: %v.", got, want)
	}
}
