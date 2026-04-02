/*
 * @license
 * Copyright 2026 Dynatrace LLC
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

package hcl

import "regexp"

var credPlaceholderPattern = regexp.MustCompile(`^\*{3}\d{3}\*{3}$`)

// ReplaceCredPlaceholders replaces every credential matching "***...***" in apiValue with the value from stateValue.
func ReplaceCredPlaceholders(stateValue, apiValue any) any {
	switch stateValueCast := stateValue.(type) {
	case map[string]any:
		apiMap, ok := apiValue.(map[string]any)
		if !ok {
			return apiValue
		}
		for key, apiMapValue := range apiMap {
			if stateMapValue, exists := stateValueCast[key]; exists {
				apiMap[key] = ReplaceCredPlaceholders(stateMapValue, apiMapValue)
			}
		}
	case []any:
		apiSlice, ok := apiValue.([]any)
		if !ok || len(stateValueCast) != len(apiSlice) {
			return apiValue
		}
		for i := range apiSlice {
			apiSlice[i] = ReplaceCredPlaceholders(stateValueCast[i], apiSlice[i])
		}
	case string:
		apiStr, ok := apiValue.(string)
		if ok && credPlaceholderPattern.MatchString(apiStr) {
			return stateValueCast
		}
	}
	return apiValue
}
