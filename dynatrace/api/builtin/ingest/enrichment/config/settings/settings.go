/**
* @license
* Copyright 2026 Dynatrace LLC
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

package config

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Condition       *string `json:"condition,omitempty"`       // Optional DQL condition to filter signals.
	InsertAfter     *string `json:"-"`                         // Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched
	Scope           *string `json:"-" scope:"scope"`           // The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP, AWS_ACCOUNT, AZURE_MICROSOFT_RESOURCES_SUBSCRIPTIONS, GCP_PROJECT). Omit this property if you want to cover the whole environment.
	Target          string  `json:"target"`                    // The field or tag to populate with the source value.
	Type            Type    `json:"type"`                      // Defines the enrichment type for this rule. Possible values: `AWS_TAG`, `AZURE_TAG`, `CUSTOM`, `GCP_LABEL`, `GCP_TAG`, `HOST_PROCESS_PROPERTY`, `K8S_NAMESPACE_ANNOTATION`, `K8S_NAMESPACE_LABEL`, `K8S_POD_ANNOTATION`, `K8S_POD_LABEL`, `K8S_WORKLOAD_ANNOTATION`, `K8S_WORKLOAD_LABEL`
	ValueExtraction *string `json:"valueExtraction,omitempty"` // An optional DPL expression to extract the enrichment value. This is only allowed when the enrichment type is set to 'Host/process property'.
	ValueSource     string  `json:"valueSource"`               // Defines the value source for the enrichment. Depending on the enrichment type, this is a key (for Kubernetes labels/annotations, cloud tags/labels, and host process properties) or a literal (for type 'Custom').
}

func (me *Settings) Name() string {
	return "ingest_enrichment_config"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "Optional DQL condition to filter signals.",
			Optional:    true, // nullable
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Computed:    true,
			Optional:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, KUBERNETES_CLUSTER, HOST_GROUP, AWS_ACCOUNT, AZURE_MICROSOFT_RESOURCES_SUBSCRIPTIONS, GCP_PROJECT). Omit this property if you want to cover the whole environment.",
			ForceNew:    true,
			Optional:    true,
			Default:     "environment",
		},
		"target": {
			Type:        schema.TypeString,
			Description: "The field or tag to populate with the source value.",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Defines the enrichment type for this rule. Possible values: `AWS_TAG`, `AZURE_TAG`, `CUSTOM`, `GCP_LABEL`, `GCP_TAG`, `HOST_PROCESS_PROPERTY`, `K8S_NAMESPACE_ANNOTATION`, `K8S_NAMESPACE_LABEL`, `K8S_POD_ANNOTATION`, `K8S_POD_LABEL`, `K8S_WORKLOAD_ANNOTATION`, `K8S_WORKLOAD_LABEL`",
			Required:    true,
		},
		"value_extraction": {
			Type:        schema.TypeString,
			Description: "An optional DPL expression to extract the enrichment value. This is only allowed when the enrichment type is set to 'Host/process property'.",
			Optional:    true, // nullable
		},
		"value_source": {
			Type:        schema.TypeString,
			Description: "Defines the value source for the enrichment. Depending on the enrichment type, this is a key (for Kubernetes labels/annotations, cloud tags/labels, and host process properties) or a literal (for type 'Custom').",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition":        me.Condition,
		"insert_after":     me.InsertAfter,
		"scope":            me.Scope,
		"target":           me.Target,
		"type":             me.Type,
		"value_extraction": me.ValueExtraction,
		"value_source":     me.ValueSource,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition":        &me.Condition,
		"insert_after":     &me.InsertAfter,
		"scope":            &me.Scope,
		"target":           &me.Target,
		"type":             &me.Type,
		"value_extraction": &me.ValueExtraction,
		"value_source":     &me.ValueSource,
	})
}
