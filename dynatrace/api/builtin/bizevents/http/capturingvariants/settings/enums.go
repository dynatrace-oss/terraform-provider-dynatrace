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

package capturingvariants

type ContentTypeMatcher string

var ContentTypeMatchers = struct {
	Contains   ContentTypeMatcher
	EndsWith   ContentTypeMatcher
	Equals     ContentTypeMatcher
	StartsWith ContentTypeMatcher
}{
	"CONTAINS",
	"ENDS_WITH",
	"EQUALS",
	"STARTS_WITH",
}

type Parser string

var Parsers = struct {
	Json       Parser
	Raw        Parser
	Text       Parser
	Urlencoded Parser
	Xml        Parser
}{
	"JSON",
	"Raw",
	"Text",
	"URL encoded",
	"XML",
}
