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

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PredefinedEventFilter struct {
	EventType EventType `json:"eventType"` // Filter problems by a Dynatrace event type
	Negate    bool      `json:"negate"`    // Negate the given event type
}

func (me *PredefinedEventFilter) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the predefined event. Possible values are `APPLICATION_ERROR_RATE_INCREASED`, `APPLICATION_SLOWDOWN`, `APPLICATION_UNEXPECTED_HIGH_LOAD`, `APPLICATION_UNEXPECTED_LOW_LOAD`, `AWS_LAMBDA_HIGH_ERROR_RATE`, `CUSTOM_APPLICATION_ERROR_RATE_INCREASED`, `CUSTOM_APPLICATION_SLOWDOWN`, `CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD`, `CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD`, `CUSTOM_APP_CRASH_RATE_INCREASED`, `DATABASE_CONNECTION_FAILURE`, `DATA_CENTER_SERVICE_PERFORMANCE_DEGRADATION`, `DATA_CENTER_SERVICE_UNAVAILABLE`, `EBS_VOLUME_HIGH_LATENCY`, `EC2_HIGH_CPU`, `ELB_HIGH_BACKEND_ERROR_RATE`, `ENTERPRICE_APPLICATION_PERFORMANCE_DEGRADATION`, `ENTERPRISE_APPLICATION_UNAVAILABLE`, `ESXI_GUEST_ACTIVE_SWAP_WAIT`, `ESXI_GUEST_CPU_LIMIT_REACHED`, `ESXI_HOST_CPU_SATURATION`, `ESXI_HOST_DATASTORE_LOW_DISK_SPACE`, `ESXI_HOST_DISK_QUEUE_SLOW`, `ESXI_HOST_DISK_SLOW`, `ESXI_HOST_MEMORY_SATURATION`, `ESXI_HOST_NETWORK_PROBLEMS`, `ESXI_HOST_OVERLOADED_STORAGE`, `ESXI_VM_IMPACT_HOST_CPU_SATURATION`, `ESXI_VM_IMPACT_HOST_MEMORY_SATURATION`, `EXTERNAL_SYNTHETIC_TEST_OUTAGE`, `EXTERNAL_SYNTHETIC_TEST_SLOWDOWN`, `HOST_OF_SERVICE_UNAVAILABLE`, `HTTP_CHECK_GLOBAL_OUTAGE`, `HTTP_CHECK_LOCAL_OUTAGE`, `HTTP_CHECK_TEST_LOCATION_SLOWDOWN`, `MOBILE_APPLICATION_ERROR_RATE_INCREASED`, `MOBILE_APPLICATION_SLOWDOWN`, `MOBILE_APPLICATION_UNEXPECTED_HIGH_LOAD`, `MOBILE_APPLICATION_UNEXPECTED_LOW_LOAD`, `MOBILE_APP_CRASH_RATE_INCREASED`, `MONITORING_UNAVAILABLE`, `OSI_DISK_LOW_INODES`, `OSI_GRACEFULLY_SHUTDOWN`, `OSI_HIGH_CPU`, `OSI_HIGH_MEMORY`, `OSI_LOW_DISK_SPACE`, `OSI_NIC_DROPPED_PACKETS_HIGH`, `OSI_NIC_ERRORS_HIGH`, `OSI_NIC_UTILIZATION_HIGH`, `OSI_SLOW_DISK`, `OSI_UNEXPECTEDLY_UNAVAILABLE`, `PGI_OF_SERVICE_UNAVAILABLE`, `PGI_UNAVAILABLE`, `PG_LOW_INSTANCE_COUNT`, `PROCESS_CRASHED`, `PROCESS_HIGH_GC_ACTIVITY`, `PROCESS_MEMORY_RESOURCE_EXHAUSTED`, `PROCESS_NA_HIGH_CONN_FAIL_RATE`, `PROCESS_NA_HIGH_LOSS_RATE`, `PROCESS_THREADS_RESOURCE_EXHAUSTED`, `RDS_HIGH_CPU`, `RDS_HIGH_LATENCY`, `RDS_LOW_MEMORY`, `RDS_LOW_STORAGE_SPACE`, `RDS_OF_SERVICE_UNAVAILABLE`, `RDS_RESTART_SEQUENCE`, `SERVICE_ERROR_RATE_INCREASED`, `SERVICE_SLOWDOWN`, `SERVICE_UNEXPECTED_HIGH_LOAD`, `SERVICE_UNEXPECTED_LOW_LOAD`, `SYNTHETIC_GLOBAL_OUTAGE`, `SYNTHETIC_LOCAL_OUTAGE`, `SYNTHETIC_NODE_OUTAGE`, `SYNTHETIC_PRIVATE_LOCATION_OUTAGE` and `SYNTHETIC_TEST_LOCATION_SLOWDOWN`",
			Required:    true,
		},
		"negate": {
			Type:        schema.TypeBool,
			Description: "The alert triggers when the problem of specified severity arises while the specified event **is** happening (`false`) or while the specified event is **not** happening (`true`).   For example, if you chose the Slowdown (`PERFORMANCE`) severity and Unexpected high traffic (`APPLICATION_UNEXPECTED_HIGH_LOAD`) event with **negate** set to `true`, the alerting profile will trigger only when the slowdown problem is raised while there is no unexpected high traffic event.  Consider the following use case as an example. The Slowdown (`PERFORMANCE`) severity rule is set. Depending on the configuration of the event filter (Unexpected high traffic (`APPLICATION_UNEXPECTED_HIGH_LOAD`) event is used as an example), the options of the alerting profile is one of the following:* **negate** is set to `false`: The alert triggers when the slowdown problem is raised while unexpected high traffic event is happening.  * **negate** is set to `true`: The alert triggers when the slowdown problem is raised while there is no unexpected high traffic event.  * no event rule is set: The alert triggers when the slowdown problem is raised, regardless of any events",
			Optional:    true,
		},
	}
}

func (me *PredefinedEventFilter) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"type":   string(me.EventType),
		"negate": me.Negate,
	})
}

func (me *PredefinedEventFilter) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("type"); ok {
		me.EventType = EventType(value.(string))
	}
	if value, ok := decoder.GetOk("negate"); ok {
		me.Negate = value.(bool)
	}
	return nil
}
