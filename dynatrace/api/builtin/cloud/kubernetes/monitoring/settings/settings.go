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

package monitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	CloudApplicationPipelineEnabled bool              `json:"cloudApplicationPipelineEnabled"` // Monitor Kubernetes namespaces, services, workloads, and pods
	EventPatterns                   EventComplexTypes `json:"eventPatterns,omitempty"`         // Define Kubernetes event filters to ingest events into your environment. For more details, see the [documentation](https://dt-url.net/2201p0u).
	EventProcessingActive           bool              `json:"eventProcessingActive"`           // All events are monitored unless event filters are specified. All ingested events are subject to licensing by default.\n\nIf you have a DPS license see [licensing documentation](https://dt-url.net/cee34zj) for details.\n\nIf you have a non-DPS license see [DDUs for events](https://dt-url.net/5n03vcu) for details.
	FilterEvents                    *bool             `json:"filterEvents,omitempty"`          // Include only events specified by Events Field Selectors
	IncludeAllFdiEvents             *bool             `json:"includeAllFdiEvents,omitempty"`   // For a list of included events, see the [documentation](https://dt-url.net/l61d02no).
	OpenMetricsBuiltinEnabled       bool              `json:"openMetricsBuiltinEnabled"`       // Workload and node resource metrics are based on a subset of cAdvisor metrics. Depending on your Kubernetes cluster size, this may increase the CPU/memory resource consumption of your ActiveGate.
	OpenMetricsPipelineEnabled      bool              `json:"openMetricsPipelineEnabled"`      // For annotation guidance, see the [documentation](https://dt-url.net/g42i0ppw).
	PvcMonitoringEnabled            bool              `json:"pvcMonitoringEnabled"`            // To enable dashboards and alerts, add the [Kubernetes persistent volume claims](ui/hub/ext/com.dynatrace.extension.kubernetes-pvc) extension to your environment.
	Scope                           *string           `json:"-" scope:"scope"`                 // The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cloud_application_pipeline_enabled": {
			Type:        schema.TypeBool,
			Description: "Monitor Kubernetes namespaces, services, workloads, and pods",
			Required:    true,
		},
		"event_patterns": {
			Type:        schema.TypeList,
			Description: "Define Kubernetes event filters to ingest events into your environment. For more details, see the [documentation](https://dt-url.net/2201p0u).",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(EventComplexTypes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"event_processing_active": {
			Type:        schema.TypeBool,
			Description: "All events are monitored unless event filters are specified. All ingested events are subject to licensing by default.\n\nIf you have a DPS license see [licensing documentation](https://dt-url.net/cee34zj) for details.\n\nIf you have a non-DPS license see [DDUs for events](https://dt-url.net/5n03vcu) for details.",
			Required:    true,
		},
		"filter_events": {
			Type:        schema.TypeBool,
			Description: "Include only events specified by Events Field Selectors",
			Optional:    true, // precondition
		},
		"include_all_fdi_events": {
			Type:        schema.TypeBool,
			Description: "For a list of included events, see the [documentation](https://dt-url.net/l61d02no).",
			Optional:    true, // precondition
		},
		"open_metrics_builtin_enabled": {
			Type:        schema.TypeBool,
			Description: "Workload and node resource metrics are based on a subset of cAdvisor metrics. Depending on your Kubernetes cluster size, this may increase the CPU/memory resource consumption of your ActiveGate.",
			Required:    true,
		},
		"open_metrics_pipeline_enabled": {
			Type:        schema.TypeBool,
			Description: "For annotation guidance, see the [documentation](https://dt-url.net/g42i0ppw).",
			Required:    true,
		},
		"pvc_monitoring_enabled": {
			Type:        schema.TypeBool,
			Description: "To enable dashboards and alerts, add the [Kubernetes persistent volume claims](ui/hub/ext/com.dynatrace.extension.kubernetes-pvc) extension to your environment.",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cloud_application_pipeline_enabled": me.CloudApplicationPipelineEnabled,
		"event_patterns":                     me.EventPatterns,
		"event_processing_active":            me.EventProcessingActive,
		"filter_events":                      me.FilterEvents,
		"include_all_fdi_events":             me.IncludeAllFdiEvents,
		"open_metrics_builtin_enabled":       me.OpenMetricsBuiltinEnabled,
		"open_metrics_pipeline_enabled":      me.OpenMetricsPipelineEnabled,
		"pvc_monitoring_enabled":             me.PvcMonitoringEnabled,
		"scope":                              me.Scope,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.FilterEvents == nil) && (me.EventProcessingActive) {
		me.FilterEvents = opt.NewBool(false)
	}
	if (me.IncludeAllFdiEvents == nil) && (me.FilterEvents != nil && *me.FilterEvents) {
		me.IncludeAllFdiEvents = opt.NewBool(false)
	}
	// ---- EventPatterns EventComplexTypes -> {"expectedValue":true,"property":"filterEvents","type":"EQUALS"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cloud_application_pipeline_enabled": &me.CloudApplicationPipelineEnabled,
		"event_patterns":                     &me.EventPatterns,
		"event_processing_active":            &me.EventProcessingActive,
		"filter_events":                      &me.FilterEvents,
		"include_all_fdi_events":             &me.IncludeAllFdiEvents,
		"open_metrics_builtin_enabled":       &me.OpenMetricsBuiltinEnabled,
		"open_metrics_pipeline_enabled":      &me.OpenMetricsPipelineEnabled,
		"pvc_monitoring_enabled":             &me.PvcMonitoringEnabled,
		"scope":                              &me.Scope,
	})
}
