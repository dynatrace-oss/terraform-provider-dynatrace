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

type Loader interface {
	Load(data []byte) error
}

type Storer interface {
	Store() ([]byte, error)
}

func ToJSON[T Settings](v T) ([]byte, error) {
	if storer, ok := any(v).(Storer); ok {
		return storer.Store()
	}

	insertAfterField := getInsertAfterField(v)
	scopeField, scopeTag := getScopeField(v)
	if scopeField.IsValid() || insertAfterField != nil {
		var data []byte
		var err error
		if data, err = json.Marshal(v); err != nil {
			return nil, err
		}

		m := map[string]json.RawMessage{}
		if err = json.Unmarshal(data, &m); err != nil {
			return nil, err
		}
		if scopeField.IsValid() {
			scope := GetScope(v)
			if data, err = json.Marshal(scope); err != nil {
				return nil, err
			}
			m[scopeTag] = data
		}
		if insertAfter := GetInsertAfter(v); insertAfter != nil {
			insertAfterBytes, _ := json.Marshal(*insertAfter)
			m["insertAfter"] = insertAfterBytes
		}
		return json.MarshalIndent(m, "", "  ")
	}
	return json.Marshal(v)
}

func FromJSON[T Settings](data []byte, v T) error {
	if loader, ok := any(v).(Loader); ok {
		return loader.Load(data)
	}
	scopeField, scopeTag := getScopeField(v)
	insertAfterField := getInsertAfterField(v)
	if scopeField.IsValid() || insertAfterField != nil {
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		m := map[string]any{}
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		if scopeField.IsValid() {
			if scope, found := m[scopeTag]; found {
				SetScope(v, scope.(string))
			}
		}
		if insertAfterField != nil {
			if insertAfter, found := m["insertAfter"]; found {
				if sInsertAfter, ok := insertAfter.(string); ok {
					SetInsertAfter(v, sInsertAfter)
				}
			}
		}
		return nil
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
