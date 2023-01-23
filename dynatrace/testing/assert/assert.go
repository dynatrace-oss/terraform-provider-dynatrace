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

package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

type Assert interface {
	Errorf(format string, args ...any)
	Error(error)
	Fail()
	True(v bool)
	Nil(v any)
	Equals(expected any, actual any)
	Equalsf(expected any, actual any, format string, args ...any)
	Success(err error)
}

func New(t *testing.T) Assert {
	a := &assert{t}
	log.SetOutput(a)
	return a
}

type assert struct {
	t *testing.T
}

func (a *assert) Write(p []byte) (n int, err error) {
	a.t.Helper()
	a.t.Log(strings.TrimSuffix(string(p), "\n"))
	return 0, nil
}

func (a *assert) Success(err error) {
	a.t.Helper()
	if err != nil {
		a.t.Error(err)
	}
}

func (a *assert) Errorf(format string, args ...any) {
	a.t.Helper()
	a.t.Errorf(format, args...)
	a.t.Fail()
}

func (a *assert) Error(err error) {
	a.t.Helper()
	a.t.Error(err)
	a.t.Fail()
}

func (a *assert) Fail() {
	a.t.Helper()
	a.t.Fail()
}

func (a *assert) True(v bool) {
	a.t.Helper()
	if !v {
		a.Fail()
	}
}

func (a *assert) Nil(v any) {
	a.t.Helper()
	if v != nil {
		a.Fail()
	}
}

func (a *assert) Equals(expected any, actual any) {
	a.t.Helper()
	if ve, ok := expected.(Equaler); ok {
		if !ve.Equals(actual) {
			a.Error(errors.New("different based on Equals"))
		}
	} else if ve, ok := expected.(ExtendedEqualer); ok {
		msg, equals := ve.Equals(actual)
		if !equals {
			a.Error(errors.New(msg))
		}
	} else if res := equals(reflect.ValueOf(expected), reflect.ValueOf(actual), ""); res != "" {
		a.Errorf("%s", res)
	}
}

func Equals(expected any, actual any) (string, bool) {
	res := equals(reflect.ValueOf(expected), reflect.ValueOf(actual), "")
	return res, res == ""
}

func (a *assert) Equalsf(expected any, actual any, format string, args ...any) {
	a.t.Helper()
	if res := equals(reflect.ValueOf(expected), reflect.ValueOf(actual), ""); res != "" {
		a.Errorf("%s: %s", fmt.Sprintf(format, args...), res)
	}
}

type Lister interface {
	List() []any
}

type BreadCrumbs string

func (bc BreadCrumbs) Format(s string, v ...any) string {
	msg := s
	if len(v) > 0 {
		msg = fmt.Sprintf(s, v...)
	}
	if len(bc) == 0 {
		return msg
	}
	return fmt.Sprintf("%s: %s", bc, msg)
}

func (bc BreadCrumbs) Dot(s string) BreadCrumbs {
	if len(bc) == 0 {
		return BreadCrumbs(s)
	}
	return BreadCrumbs(fmt.Sprintf("%s.%s", bc, s))
}

func (bc BreadCrumbs) Index(idx int) BreadCrumbs {
	return BreadCrumbs(fmt.Sprintf("%s[%d]", bc, idx))
}

func (bc BreadCrumbs) Key(s string) BreadCrumbs {
	return BreadCrumbs(fmt.Sprintf("%s[\"%s\"]", bc, s))
}

type Equaler interface {
	Equals(any) bool
}

type ExtendedEqualer interface {
	Equals(any) (string, bool)
}

var EqualerType = reflect.TypeOf((*Equaler)(nil)).Elem()
var ExtendedEqualerType = reflect.TypeOf((*ExtendedEqualer)(nil)).Elem()
var StringType = reflect.TypeOf("")

func equals(expected reflect.Value, actual reflect.Value, bc BreadCrumbs) string {
	if !expected.IsValid() && !actual.IsValid() {
		return ""
	}
	if expected.IsValid() && !actual.IsValid() {
		return bc.Format("expected: valid, actual: invalid")
	}
	if !expected.IsValid() && actual.IsValid() {
		return bc.Format("expected: invalid, actual: valid")
	}

	if expected.Type() != actual.Type() {
		return bc.Format("different type: expected: %T, actual: %T", expected.Interface(), actual.Interface())
	}
	switch expected.Kind() {
	case reflect.Struct:
		etype := expected.Type()
		for i := 0; i < etype.NumField(); i++ {
			if res := equals(expected.Field(i), actual.Field(i), bc.Dot(etype.Field(i).Name)); res != "" {
				return res
			}
		}
	case reflect.String:
		if expected.String() != actual.String() {
			return bc.Format("expected: %v, actual: %v", expected.String(), actual.String())
		}
	case reflect.Slice:
		if expected.Len() != actual.Len() {
			de, _ := json.Marshal(expected.Interface())
			if len(string(de)) < 100 {
				da, _ := json.Marshal(actual.Interface())
				return bc.Format("expected length: %d, actual length: %v, details:\n  expected: %s\n    actual: %s", expected.Len(), actual.Len(), string(de), string(da))
			} else {
				return bc.Format("expected length: %d, actual length: %v", expected.Len(), actual.Len())
			}
		}
		if expected.Type().Implements(EqualerType) {
			if !expected.Interface().(Equaler).Equals(actual.Interface()) {
				return bc.Format("different based on Equals")
			}
			return ""
		} else if expected.Type().Implements(ExtendedEqualerType) {
			msg, equals := expected.Interface().(ExtendedEqualer).Equals(actual.Interface())
			if !equals {
				return bc.Format(msg)
			}
			return ""
		}

		if expected.Type().Elem().Kind() == reflect.String {
			for i := 0; i < expected.Len(); i++ {
				found := false
				for j := 0; j < actual.Len(); j++ {
					if expected.Index(i).String() == actual.Index(j).String() {
						found = true
						break
					}
				}
				if !found {
					de, _ := json.Marshal(expected.Interface())
					da, _ := json.Marshal(actual.Interface())
					return bc.Format("Elements of String Slice don't match:\n  expected: %s\n    actual:%s", string(de), string(da))
				}
			}
			return ""
		}
		for i := 0; i < expected.Len(); i++ {
			if res := equals(expected.Index(i), actual.Index(i), bc.Index(i)); res != "" {
				return res
			}
		}
	case reflect.Pointer:
		return equals(expected.Elem(), actual.Elem(), bc)
	case reflect.Bool:
		if expected.Bool() != actual.Bool() {
			return bc.Format("expected: %v, actual: %v", expected.Bool(), actual.Bool())
		}
	case reflect.Interface:
		return equals(expected.Elem(), actual.Elem(), bc)
	case reflect.Map:
		if expected.Len() != actual.Len() {
			expectedKeys := []string{}
			for _, k := range expected.MapKeys() {
				expectedKeys = append(expectedKeys, k.String())
			}
			actualKeys := []string{}
			for _, k := range actual.MapKeys() {
				actualKeys = append(actualKeys, k.String())
			}
			if len(expectedKeys) > len(actualKeys) {
				for _, k := range expectedKeys {
					found := false
					for _, k2 := range actualKeys {
						if k == k2 {
							found = true
							break
						}
					}
					if !found {
						return bc.Format("[\"%s\"] expected but not found", k)
					}
				}
			} else {
				for _, k := range actualKeys {
					found := false
					for _, k2 := range expectedKeys {
						if k == k2 {
							found = true
							break
						}
					}
					if !found {
						return bc.Format("key \"%s\", not expected", k)
					}
				}

			}
			de, _ := json.Marshal(expected.Interface())
			if len(string(de)) < 100 {
				da, _ := json.Marshal(actual.Interface())
				return bc.Format("expected length: %d, actual length: %v, details:\n  expected: %s\n    actual: %s", expected.Len(), actual.Len(), string(de), string(da))
			} else {
				return bc.Format("expected length: %d, actual length: %v", expected.Len(), actual.Len())
			}
		}
		for _, mapKey := range expected.MapKeys() {
			if res := equals(expected.MapIndex(mapKey), expected.MapIndex(mapKey), bc.Key(mapKey.String())); res != "" {
				return res
			}
		}
	case reflect.Uint:
		if expected.Interface().(uint) != actual.Interface().(uint) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(uint), actual.Interface().(uint))
		}
	case reflect.Uint8:
		if expected.Interface().(uint8) != actual.Interface().(uint8) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(uint8), actual.Interface().(uint8))
		}
	case reflect.Uint16:
		if expected.Interface().(uint16) != actual.Interface().(uint16) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(uint16), actual.Interface().(uint16))
		}
	case reflect.Uint32:
		if expected.Interface().(uint32) != actual.Interface().(uint32) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(uint32), actual.Interface().(uint32))
		}
	case reflect.Uint64:
		if expected.Interface().(uint64) != actual.Interface().(uint64) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(uint64), actual.Interface().(uint64))
		}
	case reflect.Int:
		if expected.Interface().(int) != actual.Interface().(int) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(int), actual.Interface().(int))
		}
	case reflect.Int8:
		if expected.Interface().(int8) != actual.Interface().(int8) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(int8), actual.Interface().(int8))
		}
	case reflect.Int16:
		if expected.Interface().(int16) != actual.Interface().(int16) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(int16), actual.Interface().(int16))
		}
	case reflect.Int32:
		if expected.Interface().(int32) != actual.Interface().(int32) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(int32), actual.Interface().(int32))
		}
	case reflect.Int64:
		if expected.Interface().(int64) != actual.Interface().(int64) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(int64), actual.Interface().(int64))
		}
	case reflect.Float64:
		if expected.Interface().(float64) != actual.Interface().(float64) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(float64), actual.Interface().(float64))
		}
	case reflect.Float32:
		if expected.Interface().(float32) != actual.Interface().(float32) {
			return bc.Format("expected: %v, actual: %v", expected.Interface().(float32), actual.Interface().(float32))
		}
	default:
		return bc.Format("unsupported kind %v", kindStr(expected.Kind()))
	}
	return ""
}

func kindStr(v reflect.Kind) string {
	switch v {
	case reflect.Invalid:
		return "Invalid"
	case reflect.Bool:
		return "Bool"
	case reflect.Int:
		return "Int"
	case reflect.Int8:
		return "Int8"
	case reflect.Int16:
		return "Int16"
	case reflect.Int32:
		return "Int32"
	case reflect.Int64:
		return "Int64"
	case reflect.Uint:
		return "Uint"
	case reflect.Uint8:
		return "Uint8"
	case reflect.Uint16:
		return "Uint16"
	case reflect.Uint32:
		return "Uint32"
	case reflect.Uint64:
		return "Uint64"
	case reflect.Uintptr:
		return "Uintptr"
	case reflect.Float32:
		return "Float32"
	case reflect.Float64:
		return "Float64"
	case reflect.Complex64:
		return "Complex64"
	case reflect.Complex128:
		return "Complex128"
	case reflect.Array:
		return "Array"
	case reflect.Chan:
		return "Chan"
	case reflect.Func:
		return "Func"
	case reflect.Interface:
		return "Interface"
	case reflect.Map:
		return "Map"
	case reflect.Pointer:
		return "Pointer"
	case reflect.Slice:
		return "Slice"
	case reflect.String:
		return "String"
	case reflect.Struct:
		return "Struct"
	case reflect.UnsafePointer:
		return "UnsafePointer"
	default:
		return "unknown kind"
	}
}
