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

package detectionrules

type Base string

var Bases = struct {
	FileName Base
	Fqcn     Base
	Package  Base
}{
	"FILE_NAME",
	"FQCN",
	"PACKAGE",
}

type Matcher string

var Matchers = struct {
	BeginsWith Matcher
	Contains   Matcher
}{
	"BEGINS_WITH",
	"CONTAINS",
}

type Technology string

var Technologies = struct {
	Dotnet Technology
	Go     Technology
	Java   Technology
	Nodejs Technology
	Php    Technology
	Python Technology
}{
	"dotNet",
	"Go",
	"Java",
	"Nodejs",
	"PHP",
	"Python",
}
