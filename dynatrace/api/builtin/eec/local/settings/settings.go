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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool                `json:"enabled"`                      // This setting is enabled (`true`) or disabled (`false`)
	IngestActive       *bool               `json:"ingestActive,omitempty"`       // Enable local HTTP Metric, Log and Event Ingest API
	PerformanceProfile *PerformanceProfile `json:"performanceProfile,omitempty"` // Select performance profile for Extension Execution Controller [Documentation](https://www.dynatrace.com/support/help/shortlink/extensions-concepts#resource-consumption \"More about performance profiles\"). Possible values: `DEFAULT`, `HIGH`
	Scope              *string             `json:"-" scope:"scope"`              // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
	StatsdActive       *bool               `json:"statsdActive,omitempty"`       // This is applicable only to non-containerized Linux and Windows hosts
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
			Optional:    true, // precondition
		},
		"performance_profile": {
			Type:        schema.TypeString,
			Description: "Select performance profile for Extension Execution Controller [Documentation](https://www.dynatrace.com/support/help/shortlink/extensions-concepts#resource-consumption \"More about performance profiles\"). Possible values: `DEFAULT`, `HIGH`",
			Optional:    true, // precondition
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			ForceNew:    true,
			Optional:    true,
			Default:     "environment",
		},
		"statsd_active": {
			Type:        schema.TypeBool,
			Description: "This is applicable only to non-containerized Linux and Windows hosts",
			Optional:    true, // precondition
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

func (me *Settings) HandlePreconditions() error {
	if (me.IngestActive == nil) && (me.Enabled) {
		me.IngestActive = new(false)
	}
	if (me.StatsdActive == nil) && (me.Enabled) {
		me.StatsdActive = new(false)
	}
	if (me.IngestActive != nil) && (!me.Enabled) {
		return fmt.Errorf("'ingest_active' must not be specified unless 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.PerformanceProfile != nil) && (!me.Enabled) {
		return fmt.Errorf("'performance_profile' must not be specified unless 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.PerformanceProfile == nil) && (me.Enabled) {
		return fmt.Errorf("'performance_profile' must be specified when 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	if (me.StatsdActive != nil) && (!me.Enabled) {
		return fmt.Errorf("'statsd_active' must not be specified unless 'enabled' is set to 'true'; got 'enabled'='%v'", me.Enabled)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":             &me.Enabled,
		"ingest_active":       &me.IngestActive,
		"performance_profile": &me.PerformanceProfile,
		"scope":               &me.Scope,
		"statsd_active":       &me.StatsdActive,
	})
}
