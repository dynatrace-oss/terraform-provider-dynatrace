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

package profile

type AlertingProfileEventFilterType string

var AlertingProfileEventFilterTypes = struct {
	Custom     AlertingProfileEventFilterType
	Predefined AlertingProfileEventFilterType
}{
	"CUSTOM",
	"PREDEFINED",
}

type ComparisonOperator string

var ComparisonOperators = struct {
	BeginsWith   ComparisonOperator
	Contains     ComparisonOperator
	EndsWith     ComparisonOperator
	RegexMatches ComparisonOperator
	StringEquals ComparisonOperator
}{
	"BEGINS_WITH",
	"CONTAINS",
	"ENDS_WITH",
	"REGEX_MATCHES",
	"STRING_EQUALS",
}

type EventType string

var EventTypes = struct {
	ApplicationErrorRateIncreased       EventType
	ApplicationSlowdown                 EventType
	ApplicationUnexpectedHighLoad       EventType
	ApplicationUnexpectedLowLoad        EventType
	AwsLambdaHighErrorRate              EventType
	CustomAppCrashRateIncreased         EventType
	CustomApplicationErrorRateIncreased EventType
	CustomApplicationSlowdown           EventType
	CustomApplicationUnexpectedHighLoad EventType
	CustomApplicationUnexpectedLowLoad  EventType
	DatabaseConnectionFailure           EventType
	EbsVolumeHighLatency                EventType
	Ec2HighCpu                          EventType
	ElbHighBackendErrorRate             EventType
	EsxiGuestActiveSwapWait             EventType
	EsxiGuestCpuLimitReached            EventType
	EsxiHostCpuSaturation               EventType
	EsxiHostDatastoreLowDiskSpace       EventType
	EsxiHostDiskQueueSlow               EventType
	EsxiHostDiskSlow                    EventType
	EsxiHostMemorySaturation            EventType
	EsxiHostNetworkProblems             EventType
	EsxiHostOverloadedStorage           EventType
	EsxiVmImpactHostCpuSaturation       EventType
	EsxiVmImpactHostMemorySaturation    EventType
	ExternalSyntheticTestOutage         EventType
	ExternalSyntheticTestSlowdown       EventType
	HostOfServiceUnavailable            EventType
	HttpCheckGlobalOutage               EventType
	HttpCheckLocalOutage                EventType
	HttpCheckTestLocationSlowdown       EventType
	MobileAppCrashRateIncreased         EventType
	MobileApplicationErrorRateIncreased EventType
	MobileApplicationSlowdown           EventType
	MobileApplicationUnexpectedHighLoad EventType
	MobileApplicationUnexpectedLowLoad  EventType
	MonitoringUnavailable               EventType
	MultiProtocolGlobalOutage           EventType
	MultiProtocolLocalOutage            EventType
	MultiProtocolLocationSlowdown       EventType
	OsiDiskLowInodes                    EventType
	OsiGracefullyShutdown               EventType
	OsiHighCpu                          EventType
	OsiHighMemory                       EventType
	OsiLowDiskSpace                     EventType
	OsiNicDroppedPacketsHigh            EventType
	OsiNicErrorsHigh                    EventType
	OsiNicUtilizationHigh               EventType
	OsiSlowDisk                         EventType
	OsiUnexpectedlyUnavailable          EventType
	PgLowInstanceCount                  EventType
	PgiOfServiceUnavailable             EventType
	PgiUnavailable                      EventType
	ProcessCrashed                      EventType
	ProcessHighGcActivity               EventType
	ProcessMemoryResourceExhausted      EventType
	ProcessNaHighConnFailRate           EventType
	ProcessNaHighLossRate               EventType
	ProcessThreadsResourceExhausted     EventType
	RdsHighCpu                          EventType
	RdsHighLatency                      EventType
	RdsLowMemory                        EventType
	RdsLowStorageSpace                  EventType
	RdsOfServiceUnavailable             EventType
	RdsRestartSequence                  EventType
	ServiceErrorRateIncreased           EventType
	ServiceSlowdown                     EventType
	ServiceUnexpectedHighLoad           EventType
	ServiceUnexpectedLowLoad            EventType
	SyntheticGlobalOutage               EventType
	SyntheticLocalOutage                EventType
	SyntheticNodeOutage                 EventType
	SyntheticPrivateLocationOutage      EventType
	SyntheticTestLocationSlowdown       EventType
}{
	"APPLICATION_ERROR_RATE_INCREASED",
	"APPLICATION_SLOWDOWN",
	"APPLICATION_UNEXPECTED_HIGH_LOAD",
	"APPLICATION_UNEXPECTED_LOW_LOAD",
	"AWS_LAMBDA_HIGH_ERROR_RATE",
	"CUSTOM_APP_CRASH_RATE_INCREASED",
	"CUSTOM_APPLICATION_ERROR_RATE_INCREASED",
	"CUSTOM_APPLICATION_SLOWDOWN",
	"CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD",
	"CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD",
	"DATABASE_CONNECTION_FAILURE",
	"EBS_VOLUME_HIGH_LATENCY",
	"EC2_HIGH_CPU",
	"ELB_HIGH_BACKEND_ERROR_RATE",
	"ESXI_GUEST_ACTIVE_SWAP_WAIT",
	"ESXI_GUEST_CPU_LIMIT_REACHED",
	"ESXI_HOST_CPU_SATURATION",
	"ESXI_HOST_DATASTORE_LOW_DISK_SPACE",
	"ESXI_HOST_DISK_QUEUE_SLOW",
	"ESXI_HOST_DISK_SLOW",
	"ESXI_HOST_MEMORY_SATURATION",
	"ESXI_HOST_NETWORK_PROBLEMS",
	"ESXI_HOST_OVERLOADED_STORAGE",
	"ESXI_VM_IMPACT_HOST_CPU_SATURATION",
	"ESXI_VM_IMPACT_HOST_MEMORY_SATURATION",
	"EXTERNAL_SYNTHETIC_TEST_OUTAGE",
	"EXTERNAL_SYNTHETIC_TEST_SLOWDOWN",
	"HOST_OF_SERVICE_UNAVAILABLE",
	"HTTP_CHECK_GLOBAL_OUTAGE",
	"HTTP_CHECK_LOCAL_OUTAGE",
	"HTTP_CHECK_TEST_LOCATION_SLOWDOWN",
	"MOBILE_APP_CRASH_RATE_INCREASED",
	"MOBILE_APPLICATION_ERROR_RATE_INCREASED",
	"MOBILE_APPLICATION_SLOWDOWN",
	"MOBILE_APPLICATION_UNEXPECTED_HIGH_LOAD",
	"MOBILE_APPLICATION_UNEXPECTED_LOW_LOAD",
	"MONITORING_UNAVAILABLE",
	"MULTI_PROTOCOL_GLOBAL_OUTAGE",
	"MULTI_PROTOCOL_LOCAL_OUTAGE",
	"MULTI_PROTOCOL_LOCATION_SLOWDOWN",
	"OSI_DISK_LOW_INODES",
	"OSI_GRACEFULLY_SHUTDOWN",
	"OSI_HIGH_CPU",
	"OSI_HIGH_MEMORY",
	"OSI_LOW_DISK_SPACE",
	"OSI_NIC_DROPPED_PACKETS_HIGH",
	"OSI_NIC_ERRORS_HIGH",
	"OSI_NIC_UTILIZATION_HIGH",
	"OSI_SLOW_DISK",
	"OSI_UNEXPECTEDLY_UNAVAILABLE",
	"PG_LOW_INSTANCE_COUNT",
	"PGI_OF_SERVICE_UNAVAILABLE",
	"PGI_UNAVAILABLE",
	"PROCESS_CRASHED",
	"PROCESS_HIGH_GC_ACTIVITY",
	"PROCESS_MEMORY_RESOURCE_EXHAUSTED",
	"PROCESS_NA_HIGH_CONN_FAIL_RATE",
	"PROCESS_NA_HIGH_LOSS_RATE",
	"PROCESS_THREADS_RESOURCE_EXHAUSTED",
	"RDS_HIGH_CPU",
	"RDS_HIGH_LATENCY",
	"RDS_LOW_MEMORY",
	"RDS_LOW_STORAGE_SPACE",
	"RDS_OF_SERVICE_UNAVAILABLE",
	"RDS_RESTART_SEQUENCE",
	"SERVICE_ERROR_RATE_INCREASED",
	"SERVICE_SLOWDOWN",
	"SERVICE_UNEXPECTED_HIGH_LOAD",
	"SERVICE_UNEXPECTED_LOW_LOAD",
	"SYNTHETIC_GLOBAL_OUTAGE",
	"SYNTHETIC_LOCAL_OUTAGE",
	"SYNTHETIC_NODE_OUTAGE",
	"SYNTHETIC_PRIVATE_LOCATION_OUTAGE",
	"SYNTHETIC_TEST_LOCATION_SLOWDOWN",
}

type SeverityLevel string

var SeverityLevels = struct {
	Availability          SeverityLevel
	CustomAlert           SeverityLevel
	Errors                SeverityLevel
	MonitoringUnavailable SeverityLevel
	Performance           SeverityLevel
	ResourceContention    SeverityLevel
}{
	"AVAILABILITY",
	"CUSTOM_ALERT",
	"ERRORS",
	"MONITORING_UNAVAILABLE",
	"PERFORMANCE",
	"RESOURCE_CONTENTION",
}

type TagFilterIncludeMode string

var TagFilterIncludeModes = struct {
	IncludeAll TagFilterIncludeMode
	IncludeAny TagFilterIncludeMode
	None       TagFilterIncludeMode
}{
	"INCLUDE_ALL",
	"INCLUDE_ANY",
	"NONE",
}
