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

package errors

type CustomErrorRuleKeyMatcher string

var CustomErrorRuleKeyMatchers = struct {
	BeginsWith CustomErrorRuleKeyMatcher
	Contains   CustomErrorRuleKeyMatcher
	EndsWith   CustomErrorRuleKeyMatcher
	Equals     CustomErrorRuleKeyMatcher
}{
	"BEGINS_WITH",
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
}

type CustomErrorRuleValueMatcher string

var CustomErrorRuleValueMatchers = struct {
	BeginsWith CustomErrorRuleValueMatcher
	Contains   CustomErrorRuleValueMatcher
	EndsWith   CustomErrorRuleValueMatcher
	Equals     CustomErrorRuleValueMatcher
}{
	"BEGINS_WITH",
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
}

type HTTPErrorRuleFilter string

var HTTPErrorRuleFilters = struct {
	BeginsWith HTTPErrorRuleFilter
	Contains   HTTPErrorRuleFilter
	EndsWith   HTTPErrorRuleFilter
	Equals     HTTPErrorRuleFilter
}{
	"BEGINS_WITH",
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
}
