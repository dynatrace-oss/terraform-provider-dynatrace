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

// ReplaceCredPlaceholders replaces every credential matching "***...***" in apiState with the value from tfState.
func ReplaceCredPlaceholders(tfState, apiState any) any {
	switch aVal := tfState.(type) {
	case map[string]any:
		bMap, ok := apiState.(map[string]any)
		if !ok {
			return apiState
		}
		for key, bVal := range bMap {
			if aInnerVal, exists := aVal[key]; exists {
				bMap[key] = ReplaceCredPlaceholders(aInnerVal, bVal)
			}
		}
	case []any:
		bSlice, ok := apiState.([]any)
		if !ok || len(aVal) != len(bSlice) {
			return apiState
		}
		for i := range bSlice {
			bSlice[i] = ReplaceCredPlaceholders(aVal[i], bSlice[i])
		}
	case string:
		bStr, ok := apiState.(string)
		if ok && credPlaceholderPattern.MatchString(bStr) {
			return aVal
		}
	}
	return apiState
}
