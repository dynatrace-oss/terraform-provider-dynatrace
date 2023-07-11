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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DavisProblemCategories struct {
	MonitoringUnavailable bool `json:"monitoringUnavailable"` //
	Availability          bool `json:"availability"`          //
	Error                 bool `json:"error"`                 //
	Slowdown              bool `json:"slowdown"`              //
	Resource              bool `json:"resource"`              //
	Custom                bool `json:"custom"`                //
	Info                  bool `json:"info"`                  //
}

func (me *DavisProblemCategories) Schema(prefix string) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"monitoring_unavailable": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
		"availability": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
		"error": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
		"slowdown": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
		"resource": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
		"custom": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
		"info": {
			Type:        schema.TypeBool,
			Description: "",
			Optional:    true,
			Default:     false,
		},
	}
}

func (me *DavisProblemCategories) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"monitoring_unavailable": me.MonitoringUnavailable,
		"availability":           me.Availability,
		"error":                  me.Error,
		"slowdown":               me.Slowdown,
		"resource":               me.Resource,
		"custom":                 me.Custom,
		"info":                   me.Info,
	})
}

func (me *DavisProblemCategories) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"monitoring_unavailable": &me.MonitoringUnavailable,
		"availability":           &me.Availability,
		"error":                  &me.Error,
		"slowdown":               &me.Slowdown,
		"resource":               &me.Resource,
		"custom":                 &me.Custom,
		"info":                   &me.Info,
	})
}
