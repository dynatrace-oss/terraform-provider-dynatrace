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

package metricevents

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/metricevents/settings/dimensions"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/metricevents/settings/scope"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/metricevents/settings/strategy"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MetricEvent The configuration of the metric event.
type MetricEvent struct {
	MetricID            *string                     `json:"metricId"`                      // The ID of the metric evaluated by the metric event.
	AggregationType     *AggregationType            `json:"aggregationType,omitempty"`     // How the metric data points are aggregated for the evaluation.   The timeseries must support this aggregation.
	Description         string                      `json:"description"`                   // The description of the metric event.
	Name                string                      `json:"name"`                          // The name of the metric event displayed in the UI.
	WarningReason       *WarningReason              `json:"warningReason,omitempty"`       // The reason of a warning set on the config.  The `NONE` means config has no warnings.
	MetricDimensions    dimensions.Dimensions       `json:"metricDimensions,omitempty"`    // Defines the dimensions of the metric to alert on. The filters are combined by conjunction.
	DisabledReason      *DisabledReason             `json:"disabledReason,omitempty"`      // The reason of automatic disabling.  The `NONE` means config was not disabled automatically.
	Enabled             bool                        `json:"enabled"`                       // The metric event is enabled (`true`) or disabled (`false`).
	AlertingScope       scope.AlertingScopes        `json:"alertingScope,omitempty"`       // Defines the scope of the metric event. Only one filter is allowed per filter type, except for tags, where up to 3 are allowed. The filters are combined by conjunction.
	MonitoringStrategy  strategy.MonitoringStrategy `json:"monitoringStrategy"`            // A monitoring strategy for a metric event config. This is the base version of the monitoring strategy, depending on the type,  the actual JSON may contain additional fields.
	PrimaryDimensionKey *string                     `json:"primaryDimensionKey,omitempty"` // Defines which dimension key should be used for the **alertingScope**.
	Severity            *Severity                   `json:"severity,omitempty"`            // The type of the event to trigger on the threshold violation.  The `CUSTOM_ALERT` type is not correlated with other alerts. The `INFO` type does not open a problem.
	MetricSelector      *string                     `json:"metricSelector,omitempty"`      // The metric selector that should be executed
	Unknowns            map[string]json.RawMessage  `json:"-"`
}

func (me *MetricEvent) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"metric_selector"},
			Description:   "The ID of the metric evaluated by the metric event",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the metric event displayed in the UI",
		},
		"description": {
			Type:             schema.TypeString,
			Required:         true,
			Description:      "The description of the metric event",
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"aggregation_type": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"metric_selector"},
			Description:   "How the metric data points are aggregated for the evaluation. The timeseries must support this aggregation",
		},
		"metric_selector": {
			Type:             schema.TypeString,
			Optional:         true,
			ConflictsWith:    []string{"metric_id", "scopes", "aggregation_type"},
			Description:      "The metric selector that should be executed",
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"warning_reason": {
			Type:     schema.TypeString,
			Optional: true,

			Deprecated:  "This property is not meant to be configured from the outside. It will get removed completely in future versions",
			Description: "The reason of a warning set on the config. The `NONE` means config has no warnings. The other supported value is `TOO_MANY_DIMS`",
		},
		"dimensions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Defines the dimensions of the metric to alert on. The filters are combined by conjunction",
			Elem:        &schema.Resource{Schema: new(dimensions.Dimensions).Schema()},
		},
		"disabled_reason": {
			Type:        schema.TypeString,
			Optional:    true,
			Deprecated:  "This property is not meant to be configured from the outside. It will get removed completely in future versions",
			Description: "The reason of automatic disabling.  The `NONE` means config was not disabled automatically. Possible values are `METRIC_DEFINITION_INCONSISTENCY`, `NONE`, `TOO_MANY_DIMS` and `TOPX_FORCIBLY_DEACTIVATED`",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "The metric event is enabled (`true`) or disabled (`false`)",
		},
		"scopes": {
			Type:          schema.TypeList,
			ConflictsWith: []string{"metric_selector"},
			Optional:      true,
			Description:   "Defines the scope of the metric event. Only one filter is allowed per filter type, except for tags, where up to 3 are allowed. The filters are combined by conjunction",
			Elem:          &schema.Resource{Schema: new(scope.AlertingScopes).Schema()},
		},
		"strategy": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "A monitoring strategy for a metric event config. This is the base version of the monitoring strategy, depending on the type,  the actual JSON may contain additional fields",
			Elem:        &schema.Resource{Schema: new(strategy.Wrapper).Schema()},
		},
		"primary_dimension_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Defines which dimension key should be used for the **alertingScope**",
		},
		"severity": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The type of the event to trigger on the threshold violation.  The `CUSTOM_ALERT` type is not correlated with other alerts. The `INFO` type does not open a problem",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *MetricEvent) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"metric_selector":       me.MetricSelector,
		"metric_id":             me.MetricID,
		"name":                  me.Name,
		"description":           me.Description,
		"aggregation_type":      me.AggregationType,
		"enabled":               me.Enabled,
		"primary_dimension_key": me.PrimaryDimensionKey,
		"severity":              me.Severity,
		"dimensions":            me.MetricDimensions,
		"scopes":                me.AlertingScope,
		"strategy":              &strategy.Wrapper{Strategy: me.MonitoringStrategy},
	})
}

func (me *MetricEvent) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "metricId")
		delete(me.Unknowns, "aggregationType")
		delete(me.Unknowns, "description")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "warningReason")
		delete(me.Unknowns, "metricDimensions")
		delete(me.Unknowns, "disabledReason")
		delete(me.Unknowns, "enabled")
		delete(me.Unknowns, "alertingScope")
		delete(me.Unknowns, "monitoringStrategy")
		delete(me.Unknowns, "primaryDimensionKey")
		delete(me.Unknowns, "severity")
		delete(me.Unknowns, "metadata")
		delete(me.Unknowns, "metric_selector")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}

	if value, ok := decoder.GetOk("metric_selector"); ok {
		me.MetricSelector = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("metric_id"); ok {
		me.MetricID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = value.(string)
	}
	if value, ok := decoder.GetOk("aggregation_type"); ok {
		me.AggregationType = AggregationType(value.(string)).Ref()
	}
	// if value, ok := decoder.GetOk("warning_reason"); ok {
	// 	me.WarningReason = WarningReason(value.(string)).Ref()
	// }
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("primary_dimension_key"); ok {
		me.PrimaryDimensionKey = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("severity"); ok {
		me.Severity = Severity(value.(string)).Ref()
	}
	// if value, ok := decoder.GetOk("disabled_reason"); ok {
	// 	me.DisabledReason = DisabledReason(value.(string)).Ref()
	// }
	if _, ok := decoder.GetOk("strategy.#"); ok {
		cfg := new(strategy.Wrapper)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "strategy", 0)); err != nil {
			return err
		}
		me.MonitoringStrategy = cfg.Strategy
	}
	if _, ok := decoder.GetOk("dimensions.#"); ok {
		me.MetricDimensions = dimensions.Dimensions{}
		if err := me.MetricDimensions.UnmarshalHCL(hcl.NewDecoder(decoder, "dimensions", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("scopes.#"); ok {
		me.AlertingScope = scope.AlertingScopes{}
		if err := me.AlertingScope.UnmarshalHCL(hcl.NewDecoder(decoder, "scopes", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *MetricEvent) MarshalJSON() ([]byte, error) {
	delete(me.Unknowns, "id")
	delete(me.Unknowns, "metadata")
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"metricId":            me.MetricID,
		"aggregationType":     me.AggregationType,
		"description":         me.Description,
		"name":                me.Name,
		"warningReason":       me.WarningReason,
		"primaryDimensionKey": me.PrimaryDimensionKey,
		"severity":            me.Severity,
		"disabledReason":      me.DisabledReason,
		"enabled":             me.Enabled,
		"metricDimensions":    me.MetricDimensions,
		"alertingScope":       me.AlertingScope,
		"monitoringStrategy":  me.MonitoringStrategy,
		"metricSelector":      me.MetricSelector,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *MetricEvent) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	delete(me.Unknowns, "id")
	delete(me.Unknowns, "metadata")
	wrapper := strategy.Wrapper{}
	if err := properties.UnmarshalAll(map[string]any{
		"metricId":            &me.MetricID,
		"aggregationType":     &me.AggregationType,
		"description":         &me.Description,
		"name":                &me.Name,
		"warningReason":       &me.WarningReason,
		"primaryDimensionKey": &me.PrimaryDimensionKey,
		"severity":            &me.Severity,
		"disabledReason":      &me.DisabledReason,
		"enabled":             &me.Enabled,
		"metricDimensions":    &me.MetricDimensions,
		"alertingScope":       &me.AlertingScope,
		"monitoringStrategy":  &wrapper,
		"metricSelector":      &me.MetricSelector,
	}); err != nil {
		return err
	}
	me.MonitoringStrategy = wrapper.Strategy
	return nil
}
