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

import "reflect"

func SetInsertBefore(settings Settings, insertBefore string) {
	if settings == nil {
		return
	}
	pInsertBeforeField := getInsertBeforeField(settings)
	if pInsertBeforeField == nil {
		return
	}
	insertBeforeField := *pInsertBeforeField
	if !insertBeforeField.IsValid() {
		return
	}
	if insertBeforeField.Type() == stringType {
		insertBeforeField.Set(reflect.ValueOf(insertBefore))
	} else if insertBeforeField.Type() == stringPointerType {
		insertBeforeField.Set(reflect.ValueOf(&insertBefore))
	}
}

func GetInsertBefore(settings Settings) *string {
	if settings == nil {
		return nil
	}
	pInsertBeforeField := getInsertBeforeField(settings)
	if pInsertBeforeField == nil {
		return nil
	}
	insertBeforeField := *pInsertBeforeField
	if !insertBeforeField.IsValid() {
		return nil
	}
	if insertBeforeField.Type() == stringType {
		untypedInsertBeforeValue := insertBeforeField.Interface()
		if untypedInsertBeforeValue == nil {
			return nil
		}
		if pInsertBeforeValue, ok := untypedInsertBeforeValue.(*string); ok {
			if pInsertBeforeValue == nil {
				return nil
			}
			insertBeforeValue := *pInsertBeforeValue
			return &insertBeforeValue
		}
		if ppInsertBeforeValue, ok := untypedInsertBeforeValue.(**string); ok {
			if ppInsertBeforeValue == nil {
				return nil
			}
			pInsertBeforeValue := *ppInsertBeforeValue
			if pInsertBeforeValue == nil {
				return nil
			}
			insertBeforeValue := *pInsertBeforeValue
			return &insertBeforeValue
		}
		if insertBeforeValue, ok := untypedInsertBeforeValue.(string); ok {
			return &insertBeforeValue
		}
	}
	return nil
}

func getInsertBeforeField(settings Settings) *reflect.Value {
	rv := unref(reflect.ValueOf(settings))
	t := rv.Type()
	for idx := 0; idx < t.NumField(); idx++ {
		tField := t.Field(idx)
		if tField.Name == "InsertBefore" {
			vField := rv.Field(idx)
			return &vField
		}

	}
	return nil
}
