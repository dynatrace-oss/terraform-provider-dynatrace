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

import "strings"

func JoinID(id string, context string) string {
	return id + "#MULTI#PART#ID#" + context
}

func SplitID(id string) (string, string, bool) {
	if strings.Contains(id, "#MULTI#PART#ID#") {
		parts := strings.Split(id, "#MULTI#PART#ID#")
		return parts[0], parts[1], true
	}
	return id, id, false
}
