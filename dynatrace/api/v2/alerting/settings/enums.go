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

package alerting

type Operator string

var Operators = struct {
	BeginsWith   Operator
	EndsWith     Operator
	Contains     Operator
	RegexMatches Operator
	StringEquals Operator
}{
	Operator("BEGINS_WITH"),
	Operator("ENDS_WITH"),
	Operator("CONTAINS"),
	Operator("REGEX_MATCHES"),
	Operator("STRING_EQUALS"),
}

type EventFilterType string

var EventFilterTypes = struct {
	Predefined EventFilterType
	Custom     EventFilterType
}{
	EventFilterType("PREDEFINED"),
	EventFilterType("CUSTOM"),
}

type SeverityLevel string

var SeverityLevels = struct {
	Availability          SeverityLevel
	Custom                SeverityLevel
	Error                 SeverityLevel
	MonitoringUnavailable SeverityLevel
	Slowdown              SeverityLevel
	Resource              SeverityLevel
}{
	SeverityLevel("AVAILABILITY"),
	SeverityLevel("CUSTOM_ALERT"),
	SeverityLevel("ERRORS"),
	SeverityLevel("MONITORING_UNAVAILABLE"),
	SeverityLevel("PERFORMANCE"),
	SeverityLevel("RESOURCE_CONTENTION"),
}

type TagFilterIncludeMode string

var TagFilterIncludeModes = struct {
	None       TagFilterIncludeMode
	IncludeAny TagFilterIncludeMode
	IncludeAll TagFilterIncludeMode
}{
	TagFilterIncludeMode("NONE"),
	TagFilterIncludeMode("INCLUDE_ANY"),
	TagFilterIncludeMode("INCLUDE_ALL"),
}

type EventType string

var EventTypes = struct {
	AWSCPUSaturation                        EventType
	CPUSaturation                           EventType
	ELBHighBackendErrorRate                 EventType
	ConnectivityProblem                     EventType
	CustomAppCrashRateIncrease              EventType
	CustomAppErrorRateIncrease              EventType
	CustomAppSlowdown                       EventType
	CustomAppUnexpectedLowLoad              EventType
	CustomAppUnexpectedHighLoad             EventType
	DataCenterServicePerformanceDegredation EventType
	DataCenterServiceUnvailable             EventType
	ESXiGuestCPUSaturation                  EventType
	ESXiGuestMemorySaturation               EventType
	ESXiHostCPUSaturation                   EventType
	ESXiHostMemorySaturation                EventType
}{
	EventType("EC2_HIGH_CPU"),
	EventType("OSI_HIGH_CPU"),
	EventType("ELB_HIGH_BACKEND_ERROR_RATE"),
	EventType("PROCESS_NA_HIGH_CONN_FAIL_RATE"),
	EventType("CUSTOM_APP_CRASH_RATE_INCREASED"),
	EventType("CUSTOM_APPLICATION_ERROR_RATE_INCREASED"),
	EventType("CUSTOM_APPLICATION_SLOWDOWN"),
	EventType("CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD"),
	EventType("CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD"),
	EventType("DCRUM_SVC_PERFORMANCE_DEGRADATION"),
	EventType("DCRUM_SVC_LOW_AVAILABILITY"),
	EventType("ESXI_GUEST_CPU_LIMIT_REACHED"),
	EventType("ESXI_GUEST_ACTIVE_SWAP_WAIT"),
	EventType("ESXI_HOST_CPU_SATURATION"),
	EventType("ESXI_HOST_MEMORY_SATURATION"),
}
