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

// NewInt produces a pointer to the given primitive value
func NewInt(v int) *int { return &v }

// String returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func String(v *string) string {
	if v == nil {
		return *new(string)
	}
	return *v

}

// NewString produces a pointer to the given primitive value
func NewString(v string) *string { return &v }

// Bool returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Bool(v *bool) bool {
	if v == nil {
		return *new(bool)
	}
	return *v

}

// NewBool produces a pointer to the given primitive value
func NewBool(v bool) *bool { return &v }

// NewInt8 produces a pointer to the given primitive value
func NewInt8(v int8) *int8 { return &v }

// NewInt16 produces a pointer to the given primitive value
func NewInt16(v int16) *int16 { return &v }

// Int32 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int32(v *int32) int32 {
	if v == nil {
		return *new(int32)
	}
	return *v
}

// NewInt32 produces a pointer to the given primitive value
func NewInt32(v int32) *int32 { return &v }

// Int64 returns the underlying value of the given pointer or if that pointer is `nil` a zero value
func Int64(v *int64) int64 {
	if v == nil {
		return *new(int64)
	}
	return *v

}

// NewInt64 produces a pointer to the given primitive value
func NewInt64(v int64) *int64 { return &v }

// NewUint produces a pointer to the given primitive value
func NewUint(v uint) *uint { return &v }

// NewUInt8 produces a pointer to the given primitive value
func NewUInt8(v uint8) *uint8 { return &v }

// NewUInt16 produces a pointer to the given primitive value
func NewUInt16(v uint16) *uint16 { return &v }

// NewUInt32 produces a pointer to the given primitive value
func NewUInt32(v uint32) *uint32 { return &v }

// NewUInt64 produces a pointer to the given primitive value
func NewUInt64(v uint64) *uint64 { return &v }

// NewFloat32 produces a pointer to the given primitive value
func NewFloat32(v float32) *float32 { return &v }

// NewFloat64 produces a pointer to the given primitive value
func NewFloat64(v float64) *float64 { return &v }
