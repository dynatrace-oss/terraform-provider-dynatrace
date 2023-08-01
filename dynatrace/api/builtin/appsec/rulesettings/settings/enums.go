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

package rulesettings

type ConditionOperator string

var ConditionOperators = struct {
	Equals    ConditionOperator
	NotEquals ConditionOperator
}{
	"EQUALS",
	"NOT_EQUALS",
}

type MonitoringMode string

var MonitoringModes = struct {
	MonitoringOff MonitoringMode
	MonitoringOn  MonitoringMode
}{
	"MONITORING_OFF",
	"MONITORING_ON",
}

type Property string

var Properties = struct {
	HostTag        Property
	ManagementZone Property
	ProcessTag     Property
}{
	"HOST_TAG",
	"MANAGEMENT_ZONE",
	"PROCESS_TAG",
}
