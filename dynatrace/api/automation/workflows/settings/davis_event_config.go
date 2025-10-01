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

package workflows

import (
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DavisEventConfig struct {
	EntityTagsMatch EntityTagsMatch        `json:"entityTagsMatch"`                // Possible values: `all` and `any`
	EntityTags      map[string]StringArray `json:"entityTags"`                     // key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Multiple values can be provided separated by whitespace (e.g. \"val1 val2\") and will be parsed as multiple tag values. Omit this attribute if all entities should match
	OnProblemClose  bool                   `json:"onProblemClose" default:"false"` // If set to `true` closing a problem also is considered an event that triggers the execution
	Types           []string               `json:"types" flags:"uniqueitems"`      // The types of davis events to trigger an execution
	Names           DavisEventNames        `json:"names"`
}

func (me *DavisEventConfig) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_tags_match": {
			Type:         schema.TypeString,
			Description:  "Specifies whether all or just any of the configured entity tags need to match. Possible values: `all` and `any`. Omit this attribute if all entities should match",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"all", "any"}, false),
		},
		"entity_tags": {
			Type:        schema.TypeMap,
			Description: "key/value pairs for entity tags to match for. For tags that don't require a value, just specify an empty string as value. Multiple values can be provided separated by whitespace (e.g. \"val1 val2\") and will be parsed as multiple tag values. Omit this attribute if all entities should match",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"on_problem_close": {
			Type:        schema.TypeBool,
			Description: "If set to `true` closing a problem also is considered an event that triggers the execution",
			Optional:    true,
			Default:     false,
		},
		"types": {
			Type:        schema.TypeSet,
			Description: "The types of davis events to trigger an execution. Possible values are `CUSTOM_ANNOTATION`, `APPLICATION_UNEXPECTED_HIGH_LOAD`, `APPLICATION_UNEXPECTED_LOW_LOAD`, `APPLICATION_OVERLOAD_PREVENTION`, `APPLICATION_SLOWDOWN`, `AVAILABILITY_EVENT`, `LOG_AVAILABILITY`, `EC2_HIGH_CPU`, `RDS_BACKUP_COMPLETED`, `RDS_BACKUP_STARTED`, `SYNTHETIC_GLOBAL_OUTAGE`, `SYNTHETIC_LOCAL_OUTAGE`, `SYNTHETIC_TEST_LOCATION_SLOWDOWN`, `CUSTOM_CONFIGURATION`, `PROCESS_NA_HIGH_CONN_FAIL_RATE`, `OSI_HIGH_CPU`, `CUSTOM_ALERT`, `CUSTOM_APP_CRASH_RATE_INCREASED`, `CUSTOM_APPLICATION_ERROR_RATE_INCREASED`, `CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD`, `CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD`, `CUSTOM_APPLICATION_OVERLOAD_PREVENTION`, `CUSTOM_APPLICATION_SLOWDOWN`, `PGI_CUSTOM_AVAILABILITY`, `PGI_CUSTOM_ERROR`, `CUSTOM_INFO`, `PGI_CUSTOM_PERFORMANCE`, `CUSTOM_DEPLOYMENT`, `DEPLOYMENT_CHANGED_CHANGE`, `DEPLOYMENT_CHANGED_NEW`, `DEPLOYMENT_CHANGED_REMOVED`, `EBS_VOLUME_HIGH_LATENCY`, `ERROR_EVENT`, `LOG_ERROR`, `ESXI_HOST_CONNECTION_FAILED`, `ESXI_HOST_CONNECTION_LOST`, `ESXI_GUEST_CPU_LIMIT_REACHED`, `ESXI_GUEST_ACTIVE_SWAP_WAIT`, `ESXI_HOST_CPU_SATURATION`, `ESXI_HOST_MEMORY_SATURATION`, `ESXI_HOST_MAINTENANCE`, `ESXI_HOST_NETWORK_PROBLEMS`, `ESXI_HOST_NO_CONNECTION`, `ESXI_HOST_SHUTDOWN`, `ESXI_HOST_DISK_SLOW`, `ESXI_HOST_UP`, `ESXI_HOST_TIMEOUT`, `ESXI_VM_IMPACT_HOST_CPU_SATURATION`, `ESXI_VM_IMPACT_HOST_MEMORY_SATURATION`, `DATABASE_CONNECTION_FAILURE`, `RDS_AZ_FAILOVER_COMPLETED`, `RDS_AZ_FAILOVER_STARTED`, `SERVICE_ERROR_RATE_INCREASED`, `RDS_HIGH_LATENCY`, `OSI_NIC_UTILIZATION_HIGH`, `OSI_NIC_ERRORS_HIGH`, `PGI_HAPROXY_QUEUED_REQUESTS_HIGH`, `PGI_RMQ_HIGH_FILE_DESC_USAGE`, `PGI_RMQ_HIGH_MEM_USAGE`, `PGI_RMQ_HIGH_PROCESS_USAGE`, `PGI_RMQ_HIGH_SOCKETS_USAGE`, `OSI_NIC_DROPPED_PACKETS_HIGH`, `PGI_MYSQL_SLOW_QUERIES_RATE_HIGH`, `PGI_KEYSTONE_SLOW`, `PGI_HAPROXY_SESSION_USAGE_HIGH`, `HOST_LOG_AVAILABILITY`, `HOST_LOG_ERROR`, `OSI_GRACEFULLY_SHUTDOWN`, `HOST_LOG_MATCHED`, `OSI_UNEXPECTEDLY_UNAVAILABLE`, `HOST_LOG_PERFORMANCE`, `HOST_OF_SERVICE_UNAVAILABLE`, `HTTP_CHECK_GLOBAL_OUTAGE`, `HTTP_CHECK_LOCAL_OUTAGE`, `HTTP_CHECK_TEST_LOCATION_SLOWDOWN`, `ESXI_HOST_DISK_QUEUE_SLOW`, `LOG_MATCHED`, `APPLICATION_ERROR_RATE_INCREASED`, `APPLICATION_JS_FRAMEWORK_DETECTED`, `AWS_LAMBDA_HIGH_ERROR_RATE`, `ELB_HIGH_BACKEND_ERROR_RATE`, `ELB_HIGH_FRONTEND_ERROR_RATE`, `ELB_HIGH_UNHEALTHY_HOST_RATE`, `PROCESS_HIGH_GC_ACTIVITY`, `ESXI_HOST_DATASTORE_LOW_DISK_SPACE`, `OSI_DOCKER_DEVICEMAPPER_LOW_DATA_SPACE`, `OSI_LOW_DISK_SPACE`, `OSI_DOCKER_DEVICEMAPPER_LOW_METADATA_SPACE`, `OSI_DISK_LOW_INODES`, `PGI_RMQ_LOW_DISK_SPACE`, `RDS_LOW_STORAGE_SPACE`, `MARKED_FOR_TERMINATION`, `PROCESS_MEMORY_RESOURCE_EXHAUSTED`, `OSI_HIGH_MEMORY`, `MOBILE_APP_CRASH_RATE_INCREASED`, `MOBILE_APPLICATION_ERROR_RATE_INCREASED`, `MOBILE_APPLICATION_OVERLOAD_PREVENTION`, `MOBILE_APPLICATION_SLOWDOWN`, `MOBILE_APPLICATION_UNEXPECTED_HIGH_LOAD`, `MOBILE_APPLICATION_UNEXPECTED_LOW_LOAD`, `MONITORING_UNAVAILABLE`, `PROCESS_NA_HIGH_LOSS_RATE`, `PGI_KEYSTONE_UNHEALTHY`, `ESXI_HOST_OVERLOADED_STORAGE`, `PERFORMANCE_EVENT`, `LOG_PERFORMANCE`, `PGI_LOG_AVAILABILITY`, `PGI_CRASHED_INFO`, `PROCESS_CRASHED`, `PGI_LOG_ERROR`, `PG_LOW_INSTANCE_COUNT`, `PGI_LOG_MATCHED`, `PGI_MEMDUMP`, `PGI_LOG_PERFORMANCE`, `PROCESS_RESTART`, `PGI_UNAVAILABLE`, `RDS_HIGH_CPU`, `RDS_LOW_MEMORY`, `RDS_OF_SERVICE_UNAVAILABLE`, `RESOURCE_CONTENTION_EVENT`, `SERVICE_SLOWDOWN`, `RDS_RESTART`, `RDS_RESTART_SEQUENCE`, `PGI_OF_SERVICE_UNAVAILABLE`, `OSI_SLOW_DISK`, `SYNTHETIC_NODE_OUTAGE`, `SYNTHETIC_PRIVATE_LOCATION_OUTAGE`, `EXTERNAL_SYNTHETIC_TEST_OUTAGE`, `EXTERNAL_SYNTHETIC_TEST_SLOWDOWN`, `PROCESS_THREADS_RESOURCE_EXHAUSTED`, `SERVICE_UNEXPECTED_HIGH_LOAD`, `SERVICE_UNEXPECTED_LOW_LOAD`, `ESXI_VM_DISCONNECTED`, `OPENSTACK_VM_LAUNCH_FAILED`, `ESXI_HOST_VM_MOTION_LEFT`, `ESXI_HOST_VM_MOTION_ARRIVED`, `ESXI_VM_MOTION`, `OPENSTACK_VM_MOTION`, `ESXI_VM_POWER_OFF`, `ESXI_VM_SHUTDOWN`, `OPENSTACK_HOST_VM_SHUTDOWN`, `ESXI_VM_START`, `ESXI_HOST_VM_STARTED`, `OPENSTACK_HOST_VM_STARTED`",
			MinItems:    1,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"names": {
			Type:        schema.TypeList,
			Description: "The Davis Events to match on",
			MinItems:    1,
			MaxItems:    1,
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DavisEventNames).Schema("names")},
		},
	}
}

func (me *DavisEventConfig) MarshalHCL(properties hcl.Properties) error {
	if err := me.MarshalEntityTagsHCL(properties); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"entity_tags_match": me.EntityTagsMatch,
		"on_problem_close":  me.OnProblemClose,
		"types":             me.Types,
		"names":             me.Names,
	})
}

func (me *DavisEventConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := me.UnmarshalEntityTagsHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"entity_tags_match": &me.EntityTagsMatch,
		"on_problem_close":  &me.OnProblemClose,
		"types":             &me.Types,
		"names":             &me.Names,
	})
}

func (me *DavisEventConfig) MarshalEntityTagsHCL(properties hcl.Properties) error {
	entityTagsMap := map[string]string{}
	for k, v := range me.EntityTags {
		if len(k) == 0 {
			continue
		}
		if len(v) == 0 {
			continue
		}
		entityTagsMap[k] = strings.Join([]string(v), " ")
	}
	if len(entityTagsMap) > 0 {
		if err := properties.Encode("entity_tags", entityTagsMap); err != nil {
			return err
		}
	}
	return nil
}

func (me *DavisEventConfig) UnmarshalEntityTagsHCL(decoder hcl.Decoder) error {
	entityTagsMap := map[string]string{}
	if err := decoder.Decode("entity_tags", &entityTagsMap); err != nil {
		return err
	}
	for k, v := range entityTagsMap {
		if len(k) == 0 {
			continue
		}
		if me.EntityTags == nil {
			me.EntityTags = map[string]StringArray{}
		}
		parts := strings.Split(v, " ")
		var sa StringArray
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			sa = append(sa, p)
		}
		me.EntityTags[k] = sa
	}
	return nil
}
