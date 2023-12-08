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

package kubernetes

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ActiveGateGroup                 *string           `json:"activeGateGroup,omitempty"`                 // ActiveGate Group
	AuthToken                       *string           `json:"authToken,omitempty"`                       // Create a bearer token for [Kubernetes](https://dt-url.net/og43szq \"Kubernetes\") or [OpenShift](https://dt-url.net/7l43xtp \"OpenShift\").
	CertificateCheckEnabled         *bool             `json:"certificateCheckEnabled,omitempty"`         // Require valid certificates for communication with API server (recommended)
	CloudApplicationPipelineEnabled *bool             `json:"cloudApplicationPipelineEnabled,omitempty"` // Monitor Kubernetes namespaces, services, workloads, and pods
	ClusterID                       *string           `json:"clusterId,omitempty"`                       // Unique ID of the cluster, the containerized ActiveGate is deployed to. Defaults to the UUID of the kube-system namespace. The cluster ID of containerized ActiveGates is shown on the Deployment status screen.
	ClusterIdEnabled                bool              `json:"clusterIdEnabled"`                          // For more information on local Kubernetes API monitoring, see the [documentation](https://dt-url.net/6q62uep).
	Enabled                         bool              `json:"enabled"`                                   // This setting is enabled (`true`) or disabled (`false`)
	EndpointUrl                     *string           `json:"endpointUrl,omitempty"`                     // Get the API URL for [Kubernetes](https://dt-url.net/kz23snj \"Kubernetes\") or [OpenShift](https://dt-url.net/d623xgw \"OpenShift\").
	EventPatterns                   EventComplexTypes `json:"eventPatterns,omitempty"`                   // Define Kubernetes event filters to ingest events into your environment. For more details, see the [documentation](https://dt-url.net/2201p0u).
	EventProcessingActive           *bool             `json:"eventProcessingActive,omitempty"`           // All events are monitored by default unless event filters are specified.\n\nKubernetes events are subject to Davis data units (DDU) licensing.\nSee [DDUs for events](https://dt-url.net/5n03vcu) for details.
	FilterEvents                    *bool             `json:"filterEvents,omitempty"`                    // Include only events specified by Events Field Selectors
	HostnameVerificationEnabled     *bool             `json:"hostnameVerificationEnabled,omitempty"`     // Verify hostname in certificate against Kubernetes API URL
	IncludeAllFdiEvents             *bool             `json:"includeAllFdiEvents,omitempty"`             // For a list of included events, see the [documentation](https://dt-url.net/l61d02no).
	Label                           string            `json:"label"`                                     // Renaming the cluster breaks configurations that are based on its name (e.g., management zones, and alerting).
	OpenMetricsBuiltinEnabled       *bool             `json:"openMetricsBuiltinEnabled,omitempty"`       // The workload resource metrics are based on a subset of cAdvisor metrics. Depending on your Kubernetes cluster size, this may increase the CPU/memory resource consumption of your ActiveGate.
	OpenMetricsPipelineEnabled      *bool             `json:"openMetricsPipelineEnabled,omitempty"`      // For annotation guidance, see the [documentation](https://dt-url.net/g42i0ppw).
	PvcMonitoringEnabled            *bool             `json:"pvcMonitoringEnabled,omitempty"`            // To enable dashboards and alerts, add the [Kubernetes persistent volume claims](ui/hub/ext/com.dynatrace.extension.kubernetes-pvc) extension to your environment.
	Scope                           *string           `json:"-" scope:"scope"`                           // The scope of this setting (KUBERNETES_CLUSTER)
}

func (me *Settings) SetScope(scope string) {
	me.Scope = &scope
}

func (me *Settings) GetScope() string {
	if me.Scope == nil {
		return ""
	}
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"active_gate_group": {
			Type:        schema.TypeString,
			Description: "ActiveGate Group",
			Optional:    true, // nullable & precondition
		},
		"auth_token": {
			Type:        schema.TypeString,
			Description: "Create a bearer token for [Kubernetes](https://dt-url.net/og43szq \"Kubernetes\") or [OpenShift](https://dt-url.net/7l43xtp \"OpenShift\").",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"certificate_check_enabled": {
			Type:        schema.TypeBool,
			Description: "Require valid certificates for communication with API server (recommended)",
			Optional:    true, // precondition
		},
		"cloud_application_pipeline_enabled": {
			Type:        schema.TypeBool,
			Description: "Monitor Kubernetes namespaces, services, workloads, and pods",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"cluster_id": {
			Type:        schema.TypeString,
			Description: "Unique ID of the cluster, the containerized ActiveGate is deployed to. Defaults to the UUID of the kube-system namespace. The cluster ID of containerized ActiveGates is shown on the Deployment status screen.",
			Optional:    true, // precondition
		},
		"cluster_id_enabled": {
			Type:        schema.TypeBool,
			Description: "For more information on local Kubernetes API monitoring, see the [documentation](https://dt-url.net/6q62uep).",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"endpoint_url": {
			Type:        schema.TypeString,
			Description: "Get the API URL for [Kubernetes](https://dt-url.net/kz23snj \"Kubernetes\") or [OpenShift](https://dt-url.net/d623xgw \"OpenShift\").",
			Optional:    true, // precondition
		},
		"event_patterns": {
			Type:        schema.TypeList,
			Description: "Define Kubernetes event filters to ingest events into your environment. For more details, see the [documentation](https://dt-url.net/2201p0u).",
			Optional:    true, // precondition & minobjects == 0
			Elem:        &schema.Resource{Schema: new(EventComplexTypes).Schema()},
			MinItems:    1,
			MaxItems:    1,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"event_processing_active": {
			Type:        schema.TypeBool,
			Description: "All events are monitored by default unless event filters are specified.\n\nKubernetes events are subject to Davis data units (DDU) licensing.\nSee [DDUs for events](https://dt-url.net/5n03vcu) for details.",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"filter_events": {
			Type:        schema.TypeBool,
			Description: "Include only events specified by Events Field Selectors",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"hostname_verification_enabled": {
			Type:        schema.TypeBool,
			Description: "Verify hostname in certificate against Kubernetes API URL",
			Optional:    true, // precondition
		},
		"include_all_fdi_events": {
			Type:        schema.TypeBool,
			Description: "For a list of included events, see the [documentation](https://dt-url.net/l61d02no).",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"label": {
			Type:        schema.TypeString,
			Description: "Renaming the cluster breaks configurations that are based on its name (e.g., management zones, and alerting).",
			Required:    true,
		},
		"open_metrics_builtin_enabled": {
			Type:        schema.TypeBool,
			Description: "The workload resource metrics are based on a subset of cAdvisor metrics. Depending on your Kubernetes cluster size, this may increase the CPU/memory resource consumption of your ActiveGate.",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"open_metrics_pipeline_enabled": {
			Type:        schema.TypeBool,
			Description: "For annotation guidance, see the [documentation](https://dt-url.net/g42i0ppw).",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"pvc_monitoring_enabled": {
			Type:        schema.TypeBool,
			Description: "To enable dashboards and alerts, add the [Kubernetes persistent volume claims](ui/hub/ext/com.dynatrace.extension.kubernetes-pvc) extension to your environment.",
			Optional:    true,
			Deprecated:  "This field has been moved to a new schema, please utilize the resource `dynatrace_k8s_monitoring` to configure this field.",
		},
		"scope": {
			Type:             schema.TypeString,
			Description:      "The scope of this setting (KUBERNETES_CLUSTER)",
			Optional:         true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return new == "" },
			ForceNew:         true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(), map[string]any{
			"active_gate_group":                  me.ActiveGateGroup,
			"auth_token":                         me.AuthToken,
			"certificate_check_enabled":          me.CertificateCheckEnabled,
			"cloud_application_pipeline_enabled": me.CloudApplicationPipelineEnabled,
			"cluster_id":                         me.ClusterID,
			"cluster_id_enabled":                 me.ClusterIdEnabled,
			"enabled":                            me.Enabled,
			"endpoint_url":                       me.EndpointUrl,
			"event_patterns":                     me.EventPatterns,
			"event_processing_active":            me.EventProcessingActive,
			"filter_events":                      me.FilterEvents,
			"hostname_verification_enabled":      me.HostnameVerificationEnabled,
			"include_all_fdi_events":             me.IncludeAllFdiEvents,
			"label":                              me.Label,
			"open_metrics_builtin_enabled":       me.OpenMetricsBuiltinEnabled,
			"open_metrics_pipeline_enabled":      me.OpenMetricsPipelineEnabled,
			"pvc_monitoring_enabled":             me.PvcMonitoringEnabled,
			"scope":                              me.Scope,
		},
	))
}

func (me *Settings) HandlePreconditions() error {
	if me.CertificateCheckEnabled == nil && !me.ClusterIdEnabled {
		me.CertificateCheckEnabled = opt.NewBool(false)
	}
	if me.HostnameVerificationEnabled == nil && !me.ClusterIdEnabled {
		me.HostnameVerificationEnabled = opt.NewBool(false)
	}
	if me.ActiveGateGroup == nil && !me.ClusterIdEnabled {
		return fmt.Errorf("'active_gate_group' must be specified if 'cluster_id_enabled' is set to '%v'", me.ClusterIdEnabled)
	}
	if me.AuthToken == nil && !me.ClusterIdEnabled {
		return fmt.Errorf("'auth_token' must be specified if 'cluster_id_enabled' is set to '%v'", me.ClusterIdEnabled)
	}
	if me.ClusterID == nil && me.ClusterIdEnabled {
		return fmt.Errorf("'cluster_id' must be specified if 'cluster_id_enabled' is set to '%v'", me.ClusterIdEnabled)
	}
	if me.EndpointUrl == nil && !me.ClusterIdEnabled {
		return fmt.Errorf("'endpoint_url' must be specified if 'cluster_id_enabled' is set to '%v'", me.ClusterIdEnabled)
	}
	// ---- EventPatterns EventComplexTypes -> {"expectedValue":true,"property":"filterEvents","type":"EQUALS"}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"active_gate_group":                  &me.ActiveGateGroup,
		"auth_token":                         &me.AuthToken,
		"certificate_check_enabled":          &me.CertificateCheckEnabled,
		"cloud_application_pipeline_enabled": &me.CloudApplicationPipelineEnabled,
		"cluster_id":                         &me.ClusterID,
		"cluster_id_enabled":                 &me.ClusterIdEnabled,
		"enabled":                            &me.Enabled,
		"endpoint_url":                       &me.EndpointUrl,
		"event_patterns":                     &me.EventPatterns,
		"event_processing_active":            &me.EventProcessingActive,
		"filter_events":                      &me.FilterEvents,
		"hostname_verification_enabled":      &me.HostnameVerificationEnabled,
		"include_all_fdi_events":             &me.IncludeAllFdiEvents,
		"label":                              &me.Label,
		"open_metrics_builtin_enabled":       &me.OpenMetricsBuiltinEnabled,
		"open_metrics_pipeline_enabled":      &me.OpenMetricsPipelineEnabled,
		"pvc_monitoring_enabled":             &me.PvcMonitoringEnabled,
		"scope":                              &me.Scope,
	})
}

const credsNotProvided = "REST API didn't provide credential data"

func (me *Settings) FillDemoValues() []string {
	if !me.ClusterIdEnabled {
		me.AuthToken = opt.NewString("################")
		return []string{credsNotProvided}
	}
	return nil
}
