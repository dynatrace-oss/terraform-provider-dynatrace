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

package notificationalertingprofile

type RiskLevel string

var RiskLevels = struct {
	Critical RiskLevel
	High     RiskLevel
	Low      RiskLevel
	Medium   RiskLevel
}{
	"CRITICAL",
	"HIGH",
	"LOW",
	"MEDIUM",
}

type TriggerEvent string

var TriggerEvents = struct {
	NewManagementZoneAffected TriggerEvent
	SecurityProblemOpened     TriggerEvent
}{
	"NEW_MANAGEMENT_ZONE_AFFECTED",
	"SECURITY_PROBLEM_OPENED",
}
