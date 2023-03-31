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

package custommetrics

type Operator string

var Operators = struct {
	Equals               Operator
	GreaterThan          Operator
	GreaterThanOrEqualTo Operator
	In                   Operator
	IsNotNull            Operator
	IsNull               Operator
	LessThan             Operator
	LessThanOrEqualTo    Operator
	Like                 Operator
	NotEqual             Operator
	NotLike              Operator
	StartsWith           Operator
}{
	"EQUALS",
	"GREATER_THAN",
	"GREATER_THAN_OR_EQUAL_TO",
	"IN",
	"IS_NOT_NULL",
	"IS_NULL",
	"LESS_THAN",
	"LESS_THAN_OR_EQUAL_TO",
	"LIKE",
	"NOT_EQUAL",
	"NOT_LIKE",
	"STARTS_WITH",
}

type ValueType string

var ValueTypes = struct {
	Counter ValueType
	Field   ValueType
}{
	"COUNTER",
	"FIELD",
}
