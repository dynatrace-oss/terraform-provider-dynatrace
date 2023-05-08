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

package loadactions

type LoadKpm string

var LoadKpms = struct {
	CumulativeLayoutShift  LoadKpm
	DomInteractive         LoadKpm
	FirstInputDelay        LoadKpm
	LargestContentfulPaint LoadKpm
	LoadEventEnd           LoadKpm
	LoadEventStart         LoadKpm
	ResponseEnd            LoadKpm
	ResponseStart          LoadKpm
	SpeedIndex             LoadKpm
	UserActionDuration     LoadKpm
	VisuallyComplete       LoadKpm
}{
	"CUMULATIVE_LAYOUT_SHIFT",
	"DOM_INTERACTIVE",
	"FIRST_INPUT_DELAY",
	"LARGEST_CONTENTFUL_PAINT",
	"LOAD_EVENT_END",
	"LOAD_EVENT_START",
	"RESPONSE_END",
	"RESPONSE_START",
	"SPEED_INDEX",
	"USER_ACTION_DURATION",
	"VISUALLY_COMPLETE",
}
