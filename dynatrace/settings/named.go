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

package settings

import (
	"reflect"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
)

func unref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Struct {
		return v
	}
	if v.Kind() == reflect.Ptr {
		return unref(v.Elem())
	}
	if v.Kind() == reflect.Interface {
		return unref(v.Elem())
	}
	panic(v.Kind())
}

var stringType = reflect.ValueOf("").Type()
var stringPointerType = reflect.ValueOf(opt.NewString("")).Type()
var stringSliceType = reflect.ValueOf([]string{}).Type()

func Name(v any) string {
	rv := unref(reflect.ValueOf(v))
	field := rv.FieldByName("Name")
	if field.IsValid() && field.Type() == stringType {
		return field.String()
	}
	if field.IsValid() && field.Type() == stringPointerType {
		return field.Elem().String()
	}

	method := rv.Addr().MethodByName("Name")
	if method.IsValid() && method.Type().NumIn() == 0 && method.Type().NumOut() == 1 && method.Type().Out(0) == stringType {
		results := method.Call([]reflect.Value{})
		return results[0].String()
	}
	field = rv.FieldByName("DisplayName")
	if field.IsValid() && field.Type() == stringType {
		return field.String()
	}
	field = rv.FieldByName("Label")
	if field.IsValid() && field.Type() == stringType {
		return field.String()
	}
	field = rv.FieldByName("Key")
	if field.IsValid() && field.Type() == stringType {
		return field.String()
	}
	panic(rv.Type().Name())
}
