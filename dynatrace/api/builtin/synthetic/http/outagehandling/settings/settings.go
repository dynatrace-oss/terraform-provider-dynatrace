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

package outagehandling

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	GlobalConsecutiveOutageCountThreshold *int    `json:"globalConsecutiveOutageCountThreshold,omitempty"` // (Field has overlap with `dynatrace_http_monitor`) Alert if all locations are unable to access my web application
	GlobalOutages                         bool    `json:"globalOutages"`                                   // (Field has overlap with `dynatrace_http_monitor`) Generate a problem and send an alert when the monitor is unavailable at all configured locations.
	LocalConsecutiveOutageCountThreshold  *int    `json:"localConsecutiveOutageCountThreshold,omitempty"`  // (Field has overlap with `dynatrace_http_monitor`) are unable to access my web application
	LocalLocationOutageCountThreshold     *int    `json:"localLocationOutageCountThreshold,omitempty"`     // (Field has overlap with `dynatrace_http_monitor`) Alert if at least
	LocalOutages                          bool    `json:"localOutages"`                                    // (Field has overlap with `dynatrace_http_monitor`) Generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location.
	Scope                                 *string `json:"-" scope:"scope"`                                 // The scope of this setting (HTTP_CHECK). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"global_consecutive_outage_count_threshold": {
			Type:        schema.TypeInt,
			Description: "(Field has overlap with `dynatrace_http_monitor`) Alert if all locations are unable to access my web application",
			Optional:    true,
		},
		"global_outages": {
			Type:        schema.TypeBool,
			Description: "(Field has overlap with `dynatrace_http_monitor`) Generate a problem and send an alert when the monitor is unavailable at all configured locations.",
			Required:    true,
		},
		"local_consecutive_outage_count_threshold": {
			Type:        schema.TypeInt,
			Description: "(Field has overlap with `dynatrace_http_monitor`) are unable to access my web application",
			Optional:    true,
		},
		"local_location_outage_count_threshold": {
			Type:        schema.TypeInt,
			Description: "(Field has overlap with `dynatrace_http_monitor`) Alert if at least",
			Optional:    true,
		},
		"local_outages": {
			Type:        schema.TypeBool,
			Description: "(Field has overlap with `dynatrace_http_monitor`) Generate a problem and send an alert when the monitor is unavailable for one or more consecutive runs at any location.",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HTTP_CHECK). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"global_consecutive_outage_count_threshold": me.GlobalConsecutiveOutageCountThreshold,
		"global_outages": me.GlobalOutages,
		"local_consecutive_outage_count_threshold": me.LocalConsecutiveOutageCountThreshold,
		"local_location_outage_count_threshold":    me.LocalLocationOutageCountThreshold,
		"local_outages":                            me.LocalOutages,
		"scope":                                    me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"global_consecutive_outage_count_threshold": &me.GlobalConsecutiveOutageCountThreshold,
		"global_outages": &me.GlobalOutages,
		"local_consecutive_outage_count_threshold": &me.LocalConsecutiveOutageCountThreshold,
		"local_location_outage_count_threshold":    &me.LocalLocationOutageCountThreshold,
		"local_outages":                            &me.LocalOutages,
		"scope":                                    &me.Scope,
	})
}
