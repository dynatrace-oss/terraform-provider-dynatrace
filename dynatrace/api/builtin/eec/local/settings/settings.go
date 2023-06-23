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

package local

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool                `json:"enabled"`                      // This setting is enabled (`true`) or disabled (`false`)
	IngestActive       *bool               `json:"ingestActive,omitempty"`       // Enable local HTTP Metric, Log and Event Ingest API
	PerformanceProfile *PerformanceProfile `json:"performanceProfile,omitempty"` // Possible Values: `DEFAULT`, `HIGH`
	Scope              *string             `json:"-" scope:"scope"`              // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	StatsdActive       *bool               `json:"statsdActive,omitempty"`       // Enable Dynatrace StatsD
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"ingest_active": {
			Type:        schema.TypeBool,
			Description: "Enable local HTTP Metric, Log and Event Ingest API",
			Optional:    true,
		},
		"performance_profile": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DEFAULT`, `HIGH`",
			Optional:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"statsd_active": {
			Type:        schema.TypeBool,
			Description: "Enable Dynatrace StatsD",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":             me.Enabled,
		"ingest_active":       me.IngestActive,
		"performance_profile": me.PerformanceProfile,
		"scope":               me.Scope,
		"statsd_active":       me.StatsdActive,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"enabled":             &me.Enabled,
		"ingest_active":       &me.IngestActive,
		"performance_profile": &me.PerformanceProfile,
		"scope":               &me.Scope,
		"statsd_active":       &me.StatsdActive,
	})
	if me.IngestActive == nil && me.Enabled {
		me.IngestActive = opt.NewBool(false)
	}
	if me.StatsdActive == nil && me.Enabled {
		me.StatsdActive = opt.NewBool(false)
	}
	return err
}
