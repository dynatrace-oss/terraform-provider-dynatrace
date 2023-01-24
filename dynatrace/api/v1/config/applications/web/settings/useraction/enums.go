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

package useraction

type MatchEntity string

var MatchEntities = struct {
	ActionName         MatchEntity
	CSSSelector        MatchEntity
	JavaScriptVariable MatchEntity
	MetaTag            MatchEntity
	PagePath           MatchEntity
	PageTitle          MatchEntity
	PageURL            MatchEntity
	URLAnchor          MatchEntity
	XHRURL             MatchEntity
}{
	"ActionName",
	"CssSelector",
	"JavaScriptVariable",
	"MetaTag",
	"PagePath",
	"PageTitle",
	"PageUrl",
	"UrlAnchor",
	"XhrUrl",
}

type ActionType string

var ActionTypes = struct {
	Custom ActionType
	Load   ActionType
	XHR    ActionType
}{
	"Custom",
	"Load",
	"Xhr",
}
