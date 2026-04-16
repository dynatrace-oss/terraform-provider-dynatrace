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

package opt

// Int returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int(v *int) int {
	if v == nil {
		return *new(int)
	}
	return *v
}

// String returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func String(v *string) string {
	if v == nil {
		return *new(string)
	}
	return *v

}

// Bool returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Bool(v *bool) bool {
	if v == nil {
		return *new(bool)
	}
	return *v

}

// Int8 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int8(v *int8) int8 {
	if v == nil {
		return *new(int8)
	}
	return *v

}

// Int16 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int16(v *int16) int16 {
	if v == nil {
		return *new(int16)
	}
	return *v

}

// Int32 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int32(v *int32) int32 {
	if v == nil {
		return *new(int32)
	}
	return *v
}

// Int64 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int64(v *int64) int64 {
	if v == nil {
		return *new(int64)
	}
	return *v

}

// Uint returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Uint(v *uint) uint {
	if v == nil {
		return *new(uint)
	}
	return *v

}

// UInt8 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func UInt8(v *uint8) uint8 {
	if v == nil {
		return *new(uint8)
	}
	return *v

}

// UInt16 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func UInt16(v *uint16) uint16 {
	if v == nil {
		return *new(uint16)
	}
	return *v

}

// UInt32 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func UInt32(v *uint32) uint32 {
	if v == nil {
		return *new(uint32)
	}
	return *v

}

// UInt64 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func UInt64(v *uint64) uint64 {
	if v == nil {
		return *new(uint64)
	}
	return *v

}

// Float32 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Float32(v *float32) float32 {
	if v == nil {
		return *new(float32)
	}
	return *v
}

// Float64 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Float64(v *float64) float64 {
	if v == nil {
		return *new(float64)
	}
	return *v
}
