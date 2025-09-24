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

package sitereliabilityguardian

type ComparisonOperator string

var ComparisonOperators = struct {
	GreaterThanOrEqual ComparisonOperator
	LessThanOrEqual    ComparisonOperator
}{
	"GREATER_THAN_OR_EQUAL",
	"LESS_THAN_OR_EQUAL",
}

type EventKind string

var EventKinds = struct {
	BizEvent  EventKind
	SdlcEvent EventKind
}{
	"BIZ_EVENT",
	"SDLC_EVENT",
}

type ObjectiveType string

var ObjectiveTypes = struct {
	Dql          ObjectiveType
	ReferenceSlo ObjectiveType
}{
	"DQL",
	"REFERENCE_SLO",
}
