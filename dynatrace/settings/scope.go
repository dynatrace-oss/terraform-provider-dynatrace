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

type ScopeAware interface {
	SetScope(string)
	GetScope() string
}

func SetScope(settings Settings, scope string) {
	if scopeAware, ok := any(settings).(ScopeAware); ok {
		scopeAware.SetScope(scope)
		return
	}
	scopeField, _ := getScopeField(settings)
	if !scopeField.IsValid() {
		return
	}
	if scopeField.Type() == stringType {
		scopeField.Set(reflect.ValueOf(scope))
	} else if scopeField.Type() == stringPointerType {
		scopeField.Set(reflect.ValueOf(&scope))
	}
}

func GetScope(settings Settings) string {
	if scopeAware, ok := any(settings).(ScopeAware); ok {
		return scopeAware.GetScope()
	}

	scopeField, _ := getScopeField(settings)
	if !scopeField.IsValid() {
		return "environment"
	}
	if scopeField.Type() == stringType {
		return scopeField.Interface().(string)
	}
	if scopeField.Type() == stringPointerType {
		pScope := scopeField.Interface().(*string)
		return *pScope
	}
	return "environment"
}

func getScopeField(settings Settings) (reflect.Value, string) {
	rv := unref(reflect.ValueOf(settings))
	t := rv.Type()
	for idx := 0; idx < t.NumField(); idx++ {
		tField := t.Field(idx)
		tag := tField.Tag
		scopeTag := tag.Get("scope")
		if len(scopeTag) > 0 {
			return rv.Field(idx), scopeTag
		}
	}
	return reflect.Value{}, ""
}
