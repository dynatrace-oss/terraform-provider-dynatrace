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

package custominjectionrules

type Operator string

var Operators = struct {
	Allpages Operator
	Contains Operator
	Ends     Operator
	Equals   Operator
	Starts   Operator
}{
	"AllPages",
	"Contains",
	"Ends",
	"Equals",
	"Starts",
}

type Rule string

var Rules = struct {
	Afterspecifichtml  Rule
	Automatic          Rule
	Beforespecifichtml Rule
	Donotinject        Rule
}{
	"AfterSpecificHtml",
	"Automatic",
	"BeforeSpecificHtml",
	"DoNotInject",
}
