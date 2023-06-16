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

package match

type Comparison string

var Comparisons = struct {
	Equals           Comparison
	Contains         Comparison
	StartsWith       Comparison
	EndsWith         Comparison
	DoesNotEqual     Comparison
	DoesNotContain   Comparison
	DoesNotStartWith Comparison
	DoesNotEndWith   Comparison
}{
	Equals:           Comparison("EQUALS"),
	Contains:         Comparison("CONTAINS"),
	StartsWith:       Comparison("STARTS_WITH"),
	EndsWith:         Comparison("ENDS_WITH"),
	DoesNotEqual:     Comparison("DOES_NOT_EQUAL"),
	DoesNotContain:   Comparison("DOES_NOT_CONTAIN"),
	DoesNotStartWith: Comparison("DOES_NOT_START_WITH"),
	DoesNotEndWith:   Comparison("DOES_NOT_END_WITH"),
}
