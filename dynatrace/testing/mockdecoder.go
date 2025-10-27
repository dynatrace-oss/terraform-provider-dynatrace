/*
 * @license
 * Copyright 2025 Dynatrace LLC
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package testing

import (
	"reflect"
)

type MockDecoder struct {
	Elements map[string]any
}

func (m MockDecoder) GetOkExists(key string) (any, bool) {
	panic("not implemented")
}

func (m MockDecoder) GetOk(key string) (any, bool) {
	panic("not implemented")
}

func (m MockDecoder) Get(key string) any {
	panic("not implemented")
}

func (m MockDecoder) Decode(key string, v any) error {
	panic("not implemented")
}

func (m MockDecoder) DecodeAll(m2 map[string]any) error {
	panic("not implemented")
}

func (m MockDecoder) DecodeAny(m2 map[string]any) (any, error) {
	panic("not implemented")
}

func (m MockDecoder) DecodeSlice(key string, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv.Elem().Set(reflect.ValueOf(m.Elements[key]))
	}
	return nil
}

func (m MockDecoder) Path() string {
	panic("not implemented")
}
