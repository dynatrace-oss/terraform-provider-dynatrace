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

package detection

// Mode How to detect response time degradation: automatically, or based on fixed thresholds, or do not detect.
type Mode string

// Modes offers the known enum values
var Modes = struct {
	DetectAutomatically        Mode
	DetectUsingFixedThresholds Mode
	DontDetect                 Mode
}{
	"DETECT_AUTOMATICALLY",
	"DETECT_USING_FIXED_THRESHOLDS",
	"DONT_DETECT",
}
