package alerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PredefinedEventFilter Configuration of a predefined event filter.
type PredefinedEventFilter struct {
	EventType EventType                  `json:"eventType"` // The type of the predefined event.
	Negate    bool                       `json:"negate"`    // The alert triggers when the problem of specified severity arises while the specified event **is** happening (`false`) or while the specified event is **not** happening (`true`).   For example, if you chose the Slowdown (`PERFORMANCE`) severity and Unexpected high traffic (`APPLICATION_UNEXPECTED_HIGH_LOAD`) event with **negate** set to `true`, the alerting profile will trigger only when the slowdown problem is raised while there is no unexpected high traffic event.  Consider the following use case as an example. The Slowdown (`PERFORMANCE`) severity rule is set. Depending on the configuration of the event filter (Unexpected high traffic (`APPLICATION_UNEXPECTED_HIGH_LOAD`) event is used as an example), the behavior of the alerting profile is one of the following:* **negate** is set to `false`: The alert triggers when the slowdown problem is raised while unexpected high traffic event is happening.  * **negate** is set to `true`: The alert triggers when the slowdown problem is raised while there is no unexpected high traffic event.  * no event rule is set: The alert triggers when the slowdown problem is raised, regardless of any events.
	Unknowns  map[string]json.RawMessage `json:"-"`
}

func (me *PredefinedEventFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"event_type": {
			Type:        schema.TypeString,
			Description: "The type of the predefined event. Possible values are `APPLICATION_ERROR_RATE_INCREASED`, `APPLICATION_SLOWDOWN`, `APPLICATION_UNEXPECTED_HIGH_LOAD`, `APPLICATION_UNEXPECTED_LOW_LOAD`, `AWS_LAMBDA_HIGH_ERROR_RATE`, `CUSTOM_APPLICATION_ERROR_RATE_INCREASED`, `CUSTOM_APPLICATION_SLOWDOWN`, `CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD`, `CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD`, `CUSTOM_APP_CRASH_RATE_INCREASED`, `DATABASE_CONNECTION_FAILURE`, `DATA_CENTER_SERVICE_PERFORMANCE_DEGRADATION`, `DATA_CENTER_SERVICE_UNAVAILABLE`, `EBS_VOLUME_HIGH_LATENCY`, `EC2_HIGH_CPU`, `ELB_HIGH_BACKEND_ERROR_RATE`, `ENTERPRICE_APPLICATION_PERFORMANCE_DEGRADATION`, `ENTERPRISE_APPLICATION_UNAVAILABLE`, `ESXI_GUEST_ACTIVE_SWAP_WAIT`, `ESXI_GUEST_CPU_LIMIT_REACHED`, `ESXI_HOST_CPU_SATURATION`, `ESXI_HOST_DATASTORE_LOW_DISK_SPACE`, `ESXI_HOST_DISK_QUEUE_SLOW`, `ESXI_HOST_DISK_SLOW`, `ESXI_HOST_MEMORY_SATURATION`, `ESXI_HOST_NETWORK_PROBLEMS`, `ESXI_HOST_OVERLOADED_STORAGE`, `ESXI_VM_IMPACT_HOST_CPU_SATURATION`, `ESXI_VM_IMPACT_HOST_MEMORY_SATURATION`, `EXTERNAL_SYNTHETIC_TEST_OUTAGE`, `EXTERNAL_SYNTHETIC_TEST_SLOWDOWN`, `HOST_OF_SERVICE_UNAVAILABLE`, `HTTP_CHECK_GLOBAL_OUTAGE`, `HTTP_CHECK_LOCAL_OUTAGE`, `HTTP_CHECK_TEST_LOCATION_SLOWDOWN`, `MOBILE_APPLICATION_ERROR_RATE_INCREASED`, `MOBILE_APPLICATION_SLOWDOWN`, `MOBILE_APPLICATION_UNEXPECTED_HIGH_LOAD`, `MOBILE_APPLICATION_UNEXPECTED_LOW_LOAD`, `MOBILE_APP_CRASH_RATE_INCREASED`, `MONITORING_UNAVAILABLE`, `OSI_DISK_LOW_INODES`, `OSI_GRACEFULLY_SHUTDOWN`, `OSI_HIGH_CPU`, `OSI_HIGH_MEMORY`, `OSI_LOW_DISK_SPACE`, `OSI_NIC_DROPPED_PACKETS_HIGH`, `OSI_NIC_ERRORS_HIGH`, `OSI_NIC_UTILIZATION_HIGH`, `OSI_SLOW_DISK`, `OSI_UNEXPECTEDLY_UNAVAILABLE`, `PGI_OF_SERVICE_UNAVAILABLE`, `PGI_UNAVAILABLE`, `PG_LOW_INSTANCE_COUNT`, `PROCESS_CRASHED`, `PROCESS_HIGH_GC_ACTIVITY`, `PROCESS_MEMORY_RESOURCE_EXHAUSTED`, `PROCESS_NA_HIGH_CONN_FAIL_RATE`, `PROCESS_NA_HIGH_LOSS_RATE`, `PROCESS_THREADS_RESOURCE_EXHAUSTED`, `RDS_HIGH_CPU`, `RDS_HIGH_LATENCY`, `RDS_LOW_MEMORY`, `RDS_LOW_STORAGE_SPACE`, `RDS_OF_SERVICE_UNAVAILABLE`, `RDS_RESTART_SEQUENCE`, `SERVICE_ERROR_RATE_INCREASED`, `SERVICE_SLOWDOWN`, `SERVICE_UNEXPECTED_HIGH_LOAD`, `SERVICE_UNEXPECTED_LOW_LOAD`, `SYNTHETIC_GLOBAL_OUTAGE`, `SYNTHETIC_LOCAL_OUTAGE`, `SYNTHETIC_NODE_OUTAGE`, `SYNTHETIC_PRIVATE_LOCATION_OUTAGE` and `SYNTHETIC_TEST_LOCATION_SLOWDOWN`",
			Required:    true,
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "The alert triggers when the problem of specified severity arises while the specified event **is** happening (`false`) or while the specified event is **not** happening (`true`).   For example, if you chose the Slowdown (`PERFORMANCE`) severity and Unexpected high traffic (`APPLICATION_UNEXPECTED_HIGH_LOAD`) event with **negate** set to `true`, the alerting profile will trigger only when the slowdown problem is raised while there is no unexpected high traffic event.  Consider the following use case as an example. The Slowdown (`PERFORMANCE`) severity rule is set. Depending on the configuration of the event filter (Unexpected high traffic (`APPLICATION_UNEXPECTED_HIGH_LOAD`) event is used as an example), the behavior of the alerting profile is one of the following:* **negate** is set to `false`: The alert triggers when the slowdown problem is raised while unexpected high traffic event is happening.  * **negate** is set to `true`: The alert triggers when the slowdown problem is raised while there is no unexpected high traffic event.  * no event rule is set: The alert triggers when the slowdown problem is raised, regardless of any events",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *PredefinedEventFilter) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"event_type": me.EventType,
		"negate":     me.Negate,
	})
}

func (me *PredefinedEventFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "event_type")
		delete(me.Unknowns, "negate")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("event_type"); ok {
		me.EventType = EventType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		me.Negate = value.(bool)
	}
	return nil
}

func (me *PredefinedEventFilter) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(me.EventType)
		if err != nil {
			return nil, err
		}
		m["eventType"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(me.Negate)
		if err != nil {
			return nil, err
		}
		m["negate"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *PredefinedEventFilter) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["eventType"]; found {
		if err := json.Unmarshal(v, &me.EventType); err != nil {
			return err
		}
	}
	if v, found := m["negate"]; found {
		if err := json.Unmarshal(v, &me.Negate); err != nil {
			return err
		}
	}

	delete(m, "eventType")
	delete(m, "negate")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}

// EventType The type of the predefined event.
type EventType string

// EventTypes offers the known enum values
var EventTypes = struct {
	ApplicationErrorRateIncreased               EventType
	ApplicationSlowdown                         EventType
	ApplicationUnexpectedHighLoad               EventType
	ApplicationUnexpectedLowLoad                EventType
	AWSLambdaHighErrorRate                      EventType
	CustomApplicationErrorRateIncreased         EventType
	CustomApplicationSlowdown                   EventType
	CustomApplicationUnexpectedHighLoad         EventType
	CustomApplicationUnexpectedLowLoad          EventType
	CustomAppCrashRateIncreased                 EventType
	DatabaseConnectionFailure                   EventType
	DataCenterServicePerformanceDegradation     EventType
	DataCenterServiceUnavailable                EventType
	EbsVolumeHighLatency                        EventType
	EC2HighCPU                                  EventType
	ElbHighBackendErrorRate                     EventType
	EnterpriceApplicationPerformanceDegradation EventType
	EnterpriseApplicationUnavailable            EventType
	ESXIGuestActiveSwapWait                     EventType
	ESXIGuestCPULimitReached                    EventType
	ESXIHostCPUSaturation                       EventType
	ESXIHostDatastoreLowDiskSpace               EventType
	ESXIHostDiskQueueSlow                       EventType
	ESXIHostDiskSlow                            EventType
	ESXIHostMemorySaturation                    EventType
	ESXIHostNetworkProblems                     EventType
	ESXIHostOverloadedStorage                   EventType
	ESXIVMImpactHostCPUSaturation               EventType
	ESXIVMImpactHostMemorySaturation            EventType
	ExternalSyntheticTestOutage                 EventType
	ExternalSyntheticTestSlowdown               EventType
	HostOfServiceUnavailable                    EventType
	HTTPCheckGlobalOutage                       EventType
	HTTPCheckLocalOutage                        EventType
	HTTPCheckTestLocationSlowdown               EventType
	MobileApplicationErrorRateIncreased         EventType
	MobileApplicationSlowdown                   EventType
	MobileApplicationUnexpectedHighLoad         EventType
	MobileApplicationUnexpectedLowLoad          EventType
	MobileAppCrashRateIncreased                 EventType
	MonitoringUnavailable                       EventType
	OsiDiskLowInodes                            EventType
	OsiGracefullyShutdown                       EventType
	OsiHighCPU                                  EventType
	OsiHighMemory                               EventType
	OsiLowDiskSpace                             EventType
	OsiNicDroppedPacketsHigh                    EventType
	OsiNicErrorsHigh                            EventType
	OsiNicUtilizationHigh                       EventType
	OsiSlowDisk                                 EventType
	OsiUnexpectedlyUnavailable                  EventType
	PgiOfServiceUnavailable                     EventType
	PgiUnavailable                              EventType
	PgLowInstanceCount                          EventType
	ProcessCrashed                              EventType
	ProcessHighGcActivity                       EventType
	ProcessMemoryResourceExhausted              EventType
	ProcessNaHighConnFailRate                   EventType
	ProcessNaHighLossRate                       EventType
	ProcessThreadsResourceExhausted             EventType
	RdsHighCPU                                  EventType
	RdsHighLatency                              EventType
	RdsLowMemory                                EventType
	RdsLowStorageSpace                          EventType
	RdsOfServiceUnavailable                     EventType
	RdsRestartSequence                          EventType
	ServiceErrorRateIncreased                   EventType
	ServiceSlowdown                             EventType
	ServiceUnexpectedHighLoad                   EventType
	ServiceUnexpectedLowLoad                    EventType
	SyntheticGlobalOutage                       EventType
	SyntheticLocalOutage                        EventType
	SyntheticNodeOutage                         EventType
	SyntheticPrivateLocationOutage              EventType
	SyntheticTestLocationSlowdown               EventType
}{
	"APPLICATION_ERROR_RATE_INCREASED",
	"APPLICATION_SLOWDOWN",
	"APPLICATION_UNEXPECTED_HIGH_LOAD",
	"APPLICATION_UNEXPECTED_LOW_LOAD",
	"AWS_LAMBDA_HIGH_ERROR_RATE",
	"CUSTOM_APPLICATION_ERROR_RATE_INCREASED",
	"CUSTOM_APPLICATION_SLOWDOWN",
	"CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD",
	"CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD",
	"CUSTOM_APP_CRASH_RATE_INCREASED",
	"DATABASE_CONNECTION_FAILURE",
	"DATA_CENTER_SERVICE_PERFORMANCE_DEGRADATION",
	"DATA_CENTER_SERVICE_UNAVAILABLE",
	"EBS_VOLUME_HIGH_LATENCY",
	"EC2_HIGH_CPU",
	"ELB_HIGH_BACKEND_ERROR_RATE",
	"ENTERPRICE_APPLICATION_PERFORMANCE_DEGRADATION",
	"ENTERPRISE_APPLICATION_UNAVAILABLE",
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
	"MOBILE_APPLICATION_ERROR_RATE_INCREASED",
	"MOBILE_APPLICATION_SLOWDOWN",
	"MOBILE_APPLICATION_UNEXPECTED_HIGH_LOAD",
	"MOBILE_APPLICATION_UNEXPECTED_LOW_LOAD",
	"MOBILE_APP_CRASH_RATE_INCREASED",
	"MONITORING_UNAVAILABLE",
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
	"PGI_OF_SERVICE_UNAVAILABLE",
	"PGI_UNAVAILABLE",
	"PG_LOW_INSTANCE_COUNT",
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
