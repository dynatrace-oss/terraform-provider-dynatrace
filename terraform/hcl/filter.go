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

import "reflect"

// FilterEmpty removes entries from a slice that match the provided default value.
// This is a generic workaround for https://github.com/hashicorp/terraform-plugin-sdk/issues/895
// which creates phantom empty entries in TypeSet blocks during updates.
// There may occur two different types of empty set items.
// 1. A zero value (e.g. FieldExtractionEntry)
// 2. A default value with all default fields set (e.g. FieldExtractionEntry{ExtractionType: "field"}) including values set by `HandlePreconditions`
//
// Usage example:
//
//	*me = hcl.FilterEmpty(*me, FieldExtractionEntry{ExtractionType: "field"}) // contains "Default: ..." values
func FilterEmpty[T any](entries []*T, defaultValue T) []*T {
	result := make([]*T, 0, len(entries))
	var zeroValue T
	defaultValues := []T{zeroValue}
	pointVal := &defaultValue

	if v, ok := any(pointVal).(Preconditioner); ok {
		// if the value with the defaults is in the state file, it also went through `HandlePreconditions`
		_ = v.HandlePreconditions()
		defaultValues = append(defaultValues, *pointVal)
	} else {
		defaultValues = append(defaultValues, defaultValue)
	}

	for _, entry := range entries {
		isEmpty := false
		for _, possibleDefault := range defaultValues {
			if reflect.DeepEqual(*entry, possibleDefault) {
				isEmpty = true
				break
			}
		}
		if !isEmpty {
			result = append(result, entry)
		}
	}
	return result
}
