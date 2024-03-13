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
)

func SetInsertAfter(settings Settings, insertAfter string) {
	if settings == nil {
		return
	}
	pInsertAfterField := getInsertAfterField(settings)
	if pInsertAfterField == nil {
		return
	}
	insertAfterField := *pInsertAfterField
	if !insertAfterField.IsValid() {
		return
	}
	if insertAfterField.Type() == stringType {
		insertAfterField.Set(reflect.ValueOf(insertAfter))
	} else if insertAfterField.Type() == stringPointerType {
		insertAfterField.Set(reflect.ValueOf(&insertAfter))
	}
}

func GetInsertAfter(settings Settings) *string {
	if settings == nil {
		return nil
	}
	pinsertAfterField := getInsertAfterField(settings)
	if pinsertAfterField == nil {
		return nil
	}
	insertAfterField := *pinsertAfterField
	if !insertAfterField.IsValid() {
		return nil
	}
	if insertAfterField.Type() == stringType {
		untypedinsertAfterValue := insertAfterField.Interface()
		if untypedinsertAfterValue == nil {
			return nil
		}
		if pinsertAfterValue, ok := untypedinsertAfterValue.(*string); ok {
			if pinsertAfterValue == nil {
				return nil
			}
			insertAfterValue := *pinsertAfterValue
			return &insertAfterValue
		}
		if ppinsertAfterValue, ok := untypedinsertAfterValue.(**string); ok {
			if ppinsertAfterValue == nil {
				return nil
			}
			pinsertAfterValue := *ppinsertAfterValue
			if pinsertAfterValue == nil {
				return nil
			}
			insertAfterValue := *pinsertAfterValue
			return &insertAfterValue
		}
		if insertAfterValue, ok := untypedinsertAfterValue.(string); ok {
			return &insertAfterValue
		}
	}
	return nil
}

func getInsertAfterField(settings Settings) *reflect.Value {
	rv := unref(reflect.ValueOf(settings))
	t := rv.Type()
	for idx := 0; idx < t.NumField(); idx++ {
		tField := t.Field(idx)
		if tField.Name == "InsertAfter" {
			vField := rv.Field(idx)
			return &vField
		}

	}
	return nil
}
