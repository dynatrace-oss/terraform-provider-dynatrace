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

package logagentfeatureflags

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	JournaldLogDetector     *bool   `json:"JournaldLogDetector,omitempty"` // Enable OneAgent to collect logs from Journald on Linux systems. \nThis setting enables:\n* Detection and to have logs ingested matching ingest rule is required.
	NewContainerLogDetector bool    `json:"NewContainerLogDetector"`       // Enable OneAgent to collect all container logs in Kubernetes environments. \nThis setting enables:\n* Detection and collection of logs from short-lived containers and processes in Kubernetes.\n* Detection and collection of logs from any processes in containers in Kubernetes. Up until now only processes detected by OneAgent are covered with the Log module.\n* Log events decoration according to semantic dictionary.\n **Note:** The matcher \"Deployment name\" in the log sources configuration will be ignored and needs to be replaced with \"Workload name\", requires **Dynatrace Operator 1.4.2+**.\n\n For more details, check our [documentation](https://dt-url.net/jn02ey0).
	Scope                   *string `json:"-" scope:"scope"`               // The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.
	UserAndEventData        *bool   `json:"UserAndEventData,omitempty"`    // Enable OneAgent to collect data from Event Logs in the User Data and Event Data sections.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"journald_log_detector": {
			Type:        schema.TypeBool,
			Description: "Enable OneAgent to collect logs from Journald on Linux systems. \nThis setting enables:\n* Detection and to have logs ingested matching ingest rule is required.",
			Optional:    true,
		},
		"new_container_log_detector": {
			Type:        schema.TypeBool,
			Description: "Enable OneAgent to collect all container logs in Kubernetes environments. \nThis setting enables:\n* Detection and collection of logs from short-lived containers and processes in Kubernetes.\n* Detection and collection of logs from any processes in containers in Kubernetes. Up until now only processes detected by OneAgent are covered with the Log module.\n* Log events decoration according to semantic dictionary.\n **Note:** The matcher \"Deployment name\" in the log sources configuration will be ignored and needs to be replaced with \"Workload name\", requires **Dynatrace Operator 1.4.2+**.\n\n For more details, check our [documentation](https://dt-url.net/jn02ey0).",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"user_and_event_data": {
			Type:        schema.TypeBool,
			Description: "Enable OneAgent to collect data from Event Logs in the User Data and Event Data sections.",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"journald_log_detector":      me.JournaldLogDetector,
		"new_container_log_detector": me.NewContainerLogDetector,
		"scope":                      me.Scope,
		"user_and_event_data":        me.UserAndEventData,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"journald_log_detector":      &me.JournaldLogDetector,
		"new_container_log_detector": &me.NewContainerLogDetector,
		"scope":                      &me.Scope,
		"user_and_event_data":        &me.UserAndEventData,
	})
}
