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

package manualinsertion

type CacheDurationType string

var CacheDurationTypes = struct {
	CacheDurationType1   CacheDurationType
	CacheDurationType12  CacheDurationType
	CacheDurationType144 CacheDurationType
	CacheDurationType24  CacheDurationType
	CacheDurationType3   CacheDurationType
	CacheDurationType6   CacheDurationType
	CacheDurationType72  CacheDurationType
}{
	"1",
	"12",
	"144",
	"24",
	"3",
	"6",
	"72",
}

type CodeSnippetType string

var CodeSnippetTypes = struct {
	Deferred      CodeSnippetType
	Synchronously CodeSnippetType
}{
	"DEFERRED",
	"SYNCHRONOUSLY",
}

type ScriptExecutionAttribute string

var ScriptExecutionAttributes = struct {
	Async ScriptExecutionAttribute
	Defer ScriptExecutionAttribute
	None  ScriptExecutionAttribute
}{
	"async",
	"defer",
	"none",
}
