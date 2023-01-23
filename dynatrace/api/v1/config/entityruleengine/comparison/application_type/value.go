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

package application_type

// Value The value to compare to.
type Value string

func (atcv *Value) String() string {
	return string(*atcv)
}

func (atcv Value) Ref() *Value {
	return &atcv
}

// Values offers the known enum values
var Values = struct {
	AgentlessMonitoring Value
	Amp                 Value
	AutoInjected        Value
	Default             Value
	SaasVendor          Value
}{
	"AGENTLESS_MONITORING",
	"AMP",
	"AUTO_INJECTED",
	"DEFAULT",
	"SAAS_VENDOR",
}
