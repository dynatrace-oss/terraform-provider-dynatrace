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
	"encoding/json"
	"reflect"
)

func ToJSON[T Settings](v T) ([]byte, error) {
	if storer, ok := any(v).(Storer); ok {
		return storer.Store()
	}
	return json.Marshal(v)
}

func FromJSON[T Settings](data []byte, v T) error {
	if loader, ok := any(v).(Loader); ok {
		return loader.Load(data)
	}
	return json.Unmarshal(data, v)
}

func Clone[T Settings](v T) (T, error) {
	var err error
	var data []byte
	var proto T
	clone := reflect.New(reflect.ValueOf(proto).Type().Elem()).Interface().(T)
	if data, err = ToJSON(v); err != nil {
		return clone, err
	}
	if err = FromJSON(data, clone); err != nil {
		return clone, err
	}
	return clone, nil
}
