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

package common

// Sensitivity Sensitivity of the threshold.
// With `low` sensitivity, high statistical confidence is used. Brief violations (for example, due to a surge in load) won't trigger alerts.
// With `high` sensitivity, no statistical confidence is used. Each violation triggers an alert.
type Sensitivity string

// Sensitivities offers the known enum values
var Sensitivities = struct {
	High   Sensitivity
	Low    Sensitivity
	Medium Sensitivity
}{
	"HIGH",
	"LOW",
	"MEDIUM",
}
