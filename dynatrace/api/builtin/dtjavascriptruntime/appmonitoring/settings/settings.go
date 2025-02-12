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

type Settings struct {
	AppMonitoring     AppMonitorings           `json:"appMonitoring,omitempty"`     // You can override the default monitoring setting for each app separately
	DefaultLogLevel   DefaultLogLevel          `json:"defaultLogLevel"`             // Possible Values: `debug`, `error`, `info`, `off`, `warn`
	DefaultTraceLevel *DefaultServerlessTraces `json:"defaultTraceLevel,omitempty"` // Possible Values: `off`, `on`
}

func (me *Settings) Name() string {
	return "app_monitoring"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app_monitoring": {
			Type:        schema.TypeList,
			Description: "You can override the default monitoring setting for each app separately",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(AppMonitorings).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"default_log_level": {
			Type:        schema.TypeString,
			Description: "Possible Values: `debug`, `error`, `info`, `off`, `warn`",
			Required:    true,
		},
		"default_trace_level": {
			Type:        schema.TypeString,
			Description: "Possible Values: `off`, `on`",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"app_monitoring":      me.AppMonitoring,
		"default_log_level":   me.DefaultLogLevel,
		"default_trace_level": me.DefaultTraceLevel,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"app_monitoring":      &me.AppMonitoring,
		"default_log_level":   &me.DefaultLogLevel,
		"default_trace_level": &me.DefaultTraceLevel,
	})
}
