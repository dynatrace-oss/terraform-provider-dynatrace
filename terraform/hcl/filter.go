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
//
// The defaultValue should represent what a phantom entry looks like — typically a
// zero-valued struct with any schema Default values applied. For example:
//
//	*me = hcl.FilterEmpty(*me, FieldExtractionEntry{ExtractionType: "field"})
func FilterEmpty[T Preconditioner](entries []T, defaultValue T) []T {
	result := make([]T, 0, len(entries))
	basicDefault := reflect.New(reflect.TypeOf(defaultValue).Elem()).Interface().(T)
	_ = defaultValue.HandlePreconditions()
	for _, entry := range entries {
		if reflect.DeepEqual(entry, basicDefault) || reflect.DeepEqual(entry, defaultValue) {
			continue
		}
		result = append(result, entry)
	}
	return result
}

func FilterEmptyWithBasic[T any](entries []*T, defaultValue T, basicValue T) []*T {
	result := make([]*T, 0, len(entries))
	for _, entry := range entries {
		if reflect.DeepEqual(*entry, basicValue) || reflect.DeepEqual(*entry, defaultValue) {
			continue
		}
		result = append(result, entry)
	}
	return result
}
