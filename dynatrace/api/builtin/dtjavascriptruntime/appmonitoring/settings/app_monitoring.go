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

package appmonitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AppMonitorings []*AppMonitoring

func (me *AppMonitorings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_monitoring": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(AppMonitoring).Schema()},
		},
	}
}

func (me AppMonitorings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("app_monitoring", me)
}

func (me *AppMonitorings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("app_monitoring", me)
}

type AppMonitoring struct {
	AppID            string                    `json:"appId"`                      // App ID
	CustomLogLevel   OverrideDefaultLogLevel   `json:"customLogLevel"`             // Possible Values: `debug`, `error`, `info`, `off`, `useDefault`, `warn`
	CustomTraceLevel *OverrideServerlessTraces `json:"customTraceLevel,omitempty"` // Possible Values: `off`, `on`, `useDefault`
}

func (me *AppMonitoring) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_id": {
			Type:        schema.TypeString,
			Description: "App ID",
			Required:    true,
		},
		"custom_log_level": {
			Type:        schema.TypeString,
			Description: "Possible Values: `debug`, `error`, `info`, `off`, `useDefault`, `warn`",
			Required:    true,
		},
		"custom_trace_level": {
			Type:        schema.TypeString,
			Description: "Possible Values: `off`, `on`, `useDefault`",
			Optional:    true,
		},
	}
}

func (me *AppMonitoring) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"app_id":             me.AppID,
		"custom_log_level":   me.CustomLogLevel,
		"custom_trace_level": me.CustomTraceLevel,
	})
}

func (me *AppMonitoring) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"app_id":             &me.AppID,
		"custom_log_level":   &me.CustomLogLevel,
		"custom_trace_level": &me.CustomTraceLevel,
	})
}
