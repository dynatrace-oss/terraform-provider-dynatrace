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

package scheduling_rules

type RuleType string

var RuleTypes = struct {
	RRule          RuleType
	Grouping       RuleType
	FixedOffset    RuleType
	RelativeOffset RuleType
}{
	"rrule",
	"grouping",
	"fixed_offset",
	"relative_offset",
}

type Frequency string

var Frequencies = struct {
	Yearly   Frequency
	Monthly  Frequency
	Weekly   Frequency
	Daily    Frequency
	Hourly   Frequency
	Minutely Frequency
	Secondly Frequency
}{
	"YEARLY",
	"MONTHLY",
	"WEEKLY",
	"DAILY",
	"HOURLY",
	"MINUTELY",
	"SECONDLY",
}

type Direction string

var Directions = struct {
	Next     Direction
	Previous Direction
}{
	"next",
	"previous",
}
