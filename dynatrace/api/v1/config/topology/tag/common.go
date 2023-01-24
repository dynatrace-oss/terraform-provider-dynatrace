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

package tag

import "strings"

func TagSubsetCheck(source []Tag, input []Tag) bool {
	for _, inputTag := range input {
		found := false
		for _, restTag := range source {
			if restTag.Key == inputTag.Key {
				if restTag.Value == nil && inputTag.Value == nil {
					found = true
					break
				} else if restTag.Value != nil && inputTag.Value != nil && *restTag.Value == *inputTag.Value {
					found = true
					break
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// StringsToTags processes the slice of string tags into a slice of tag structs
// Arguments: slice of string tags, pointer to slice of tag structs
func StringsToTags(tagList []any, tags *[]Tag) {
	for _, iTag := range tagList {
		var tag Tag
		if strings.Contains(iTag.(string), "=") {
			tagSplit := strings.Split(iTag.(string), "=")
			tag.Key = tagSplit[0]
			tag.Value = &tagSplit[1]
		} else {
			tag.Key = iTag.(string)
		}
		*tags = append(*tags, tag)
	}
}
