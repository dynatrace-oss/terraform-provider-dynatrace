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

package opentelemetrymetrics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Settings struct {
	AdditionalAttributes                   AdditionalAttributeItems `json:"additionalAttributes,omitempty"`         // When enabled, the attributes defined in the list below will be added as dimensions to ingested OTLP metrics if they are present in the OpenTelemetry resource or in the instrumentation scope.\n\n**Notes:**\n\n* Modifying this setting (renaming, disabling or removing attributes) will cause the metric to change. This may have an impact on existing dashboards, events and alerts that make use of these dimensions. In this case, they will need to be updated manually.\n\n* Dynatrace does not recommend changing/removing the attributes starting with \"dt.\". Dynatrace leverages these attributes to [Enrich metrics](https://www.dynatrace.com/support/help/extend-dynatrace/extend-metrics/reference/enrich-metrics).
	AdditionalAttributesToDimensionEnabled *bool                    `json:"additionalAttributesToDimensionEnabled"` // Add the resource and scope attributes configured below as dimensions
	MeterNameToDimensionEnabled            *bool                    `json:"meterNameToDimensionEnabled"`            // When enabled, the Meter name (also referred to as InstrumentationScope or InstrumentationLibrary in OpenTelemetry SDKs) and version will be added as dimensions (`otel.scope.name` and `otel.scope.version`) to ingested OTLP metrics.\n\n**Note:** Modifying this setting will cause the metric to change. This may have an impact on existing dashboards, events and alerts that make use of these dimensions. In this case, they will need to be updated manually.
	Scope                                  *string                  `json:"-" scope:"scope"`                        // The scope of this setting (environment-default). Omit this property if you want to cover the whole environment.
	ToDropAttributes                       DropAttributeItems       `json:"toDropAttributes,omitempty"`             // The attributes defined in the list below will be dropped from all ingested OTLP metrics.\n\nUpon ingest, the *Allow list: resource and scope attributes* above is applied first. Then, the *Deny list: all attributes* below is applied. The deny list therefore applies to all attributes from all sources (data points, scope and resource).\n\n**Notes:**\n\n* Modifying this setting (adding, renaming, disabling or removing attributes) will cause the metric to change. This may have an impact on existing dashboards, events and alerts that make use of these dimensions. In this case, they will need to be updated manually.\n\n* Dynatrace does not recommend including attributes starting with \"dt.\" to the deny list. Dynatrace leverages these attributes to [Enrich metrics](https://www.dynatrace.com/support/help/extend-dynatrace/extend-metrics/reference/enrich-metrics).
	Mode                                   Mode                     `json:"-"`
}

func (me *Settings) IsComputer() bool {
	return true
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"additional_attributes": {
			Type:        schema.TypeList,
			Description: "When enabled, the attributes defined in the list below will be added as dimensions to ingested OTLP metrics if they are present in the OpenTelemetry resource or in the instrumentation scope.\n\n**Notes:**\n\n* Modifying this setting (renaming, disabling or removing attributes) will cause the metric to change. This may have an impact on existing dashboards, events and alerts that make use of these dimensions. In this case, they will need to be updated manually.\n\n* Dynatrace does not recommend changing/removing the attributes starting with \"dt.\". Dynatrace leverages these attributes to [Enrich metrics](https://www.dynatrace.com/support/help/extend-dynatrace/extend-metrics/reference/enrich-metrics).",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(AdditionalAttributeItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"additional_attributes_to_dimension_enabled": {
			Type:        schema.TypeBool,
			Description: "Add the resource and scope attributes configured below as dimensions",
			Optional:    true,
			Computed:    true,
		},
		"meter_name_to_dimension_enabled": {
			Type:        schema.TypeBool,
			Description: "When enabled, the Meter name (also referred to as InstrumentationScope or InstrumentationLibrary in OpenTelemetry SDKs) and version will be added as dimensions (`otel.scope.name` and `otel.scope.version`) to ingested OTLP metrics.\n\n**Note:** Modifying this setting will cause the metric to change. This may have an impact on existing dashboards, events and alerts that make use of these dimensions. In this case, they will need to be updated manually",
			Optional:    true,
			Computed:    true,
		},
		"scope": {
			Type:         schema.TypeString,
			Description:  "The scope of this setting (environment-default). Omit this property if you want to cover the whole environment.",
			Optional:     true,
			Default:      "environment",
			ValidateFunc: validation.StringInSlice([]string{"environment"}, false),
			// ForceNew:    true,
		},
		"mode": {
			Type:         schema.TypeString,
			Description:  "Specifies whether the given attributes to enable (`additional_attributes`) and the attributes to drop (`to_drop_attributes`) will get applied explicitly (`EXPLICIT`) or additive (`ADDITIVE`).\n\nDefault behavior is `EXPLICIT` - in which case it is recommended to have just ONE instance of this resource\n\nWith mode `ADDITIVE` you're able to have multiple instances of this resource within the same Terraform Module.\n\n**Note:** Using `ADDITIVE` and `EXPLICIT` at the same time within differnt resource instances will lead to unexpected results.",
			Default:      "EXPLICIT",
			ValidateFunc: validation.StringInSlice([]string{"EXPLICIT", "ADDITIVE"}, false),
			Optional:     true,
		},
		"to_drop_attributes": {
			Type:        schema.TypeList,
			Description: "The attributes defined in the list below will be dropped from all ingested OTLP metrics.\n\nUpon ingest, the *Allow list: resource and scope attributes* above is applied first. Then, the *Deny list: all attributes* below is applied. The deny list therefore applies to all attributes from all sources (data points, scope and resource).\n\n**Notes:**\n\n* Modifying this setting (adding, renaming, disabling or removing attributes) will cause the metric to change. This may have an impact on existing dashboards, events and alerts that make use of these dimensions. In this case, they will need to be updated manually.\n\n* Dynatrace does not recommend including attributes starting with \"dt.\" to the deny list. Dynatrace leverages these attributes to [Enrich metrics](https://www.dynatrace.com/support/help/extend-dynatrace/extend-metrics/reference/enrich-metrics).",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(DropAttributeItems).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"additional_attributes":                      me.AdditionalAttributes,
		"additional_attributes_to_dimension_enabled": me.AdditionalAttributesToDimensionEnabled,
		"meter_name_to_dimension_enabled":            me.MeterNameToDimensionEnabled,
		"scope":                                      me.Scope,
		"to_drop_attributes":                         me.ToDropAttributes,
		"mode":                                       me.Mode,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"additional_attributes":                      &me.AdditionalAttributes,
		"additional_attributes_to_dimension_enabled": &me.AdditionalAttributesToDimensionEnabled,
		"meter_name_to_dimension_enabled":            &me.MeterNameToDimensionEnabled,
		"scope":                                      &me.Scope,
		"to_drop_attributes":                         &me.ToDropAttributes,
		"mode":                                       &me.Mode,
	})
}
