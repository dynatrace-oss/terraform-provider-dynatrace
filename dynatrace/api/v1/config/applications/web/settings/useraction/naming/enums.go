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

package naming

type Operator string

var Operators = struct {
	Contains                    Operator
	EndsWith                    Operator
	Equals                      Operator
	IsEmpty                     Operator
	IsNotEmpty                  Operator
	MatchesRegularExpression    Operator
	NotContains                 Operator
	NotEndsWith                 Operator
	NotEquals                   Operator
	NotMatchesRegularExpression Operator
	NotStartsWith               Operator
	StartsWith                  Operator
}{
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
	"IS_EMPTY",
	"IS_NOT_EMPTY",
	"MATCHES_REGULAR_EXPRESSION",
	"NOT_CONTAINS",
	"NOT_ENDS_WITH",
	"NOT_EQUALS",
	"NOT_MATCHES_REGULAR_EXPRESSION",
	"NOT_STARTS_WITH",
	"STARTS_WITH",
}

type Input string

var Inputs = struct {
	ElementIdentifier Input
	InputType         Input
	MetaData          Input
	PageTitle         Input
	PageURL           Input
	SourceURL         Input
	TopXHRURL         Input
	XHRURL            Input
}{
	"ELEMENT_IDENTIFIER",
	"INPUT_TYPE",
	"METADATA",
	"PAGE_TITLE",
	"PAGE_URL",
	"SOURCE_URL",
	"TOP_XHR_URL",
	"XHR_URL",
}

type ProcessingPart string

var ProcessingParts = struct {
	All    ProcessingPart
	Anchor ProcessingPart
	Path   ProcessingPart
}{
	"ALL",
	"ANCHOR",
	"PATH",
}

type ProcessingStepType string

var ProcessingStepTypes = struct {
	ExtractByRegularExpression   ProcessingStepType
	Replacement                  ProcessingStepType
	ReplaceIDs                   ProcessingStepType
	ReplaceWithPattern           ProcessingStepType
	ReplaceWithRegularExpression ProcessingStepType
	SubString                    ProcessingStepType
}{
	"EXTRACT_BY_REGULAR_EXPRESSION",
	"REPLACEMENT",
	"REPLACE_IDS",
	"REPLACE_WITH_PATTERN",
	"REPLACE_WITH_REGULAR_EXPRESSION",
	"SUBSTRING",
}

type PatternSearchType string

var PatternSearchTypes = struct {
	First PatternSearchType
	Last  PatternSearchType
}{
	"FIRST",
	"LAST",
}
