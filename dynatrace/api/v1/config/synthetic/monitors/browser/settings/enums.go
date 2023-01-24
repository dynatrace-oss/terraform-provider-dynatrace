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

package browser

// Orientation The orientation of the device— `portrait` or `landscape`
type Orientation string

// Orientations offers the known enum values
var Orientations = struct {
	Portrait  Orientation
	Landscape Orientation
}{
	`portrait`,
	`landscape`,
}

// ScriptType The type of monitor. Possible values are `clickpath` for clickpath monitors and `availability` for single-URL browser monitors. These monitors are only allowed to have one event of the `navigate` type
type ScriptType string

// ScriptTypes offers the known enum values
var ScriptTypes = struct {
	ClickPath    ScriptType
	Availability ScriptType
}{
	"clickpath",
	"availability",
}
