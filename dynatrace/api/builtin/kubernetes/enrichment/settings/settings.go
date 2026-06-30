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

package enrichment

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Rules Rules   `json:"rules,omitempty"` // Kubernetes Telemetry Enrichment empowers you to effectively tag your telemetry data using Kubernetes namespace labels and annotations. Additionally, it enables you to tag it for cost allocation and permission purposes.\n\n  Enrichment Options:\n\n  - **Enrich telemetry with label/annotation directly:** Tag your telemetry data with existing Kubernetes namespace labels or annotations. These will be made available as domain-specific fields (e.g., `k8s.namespace.label.your_key`). This allows for flexible pipeline routing, bucket selection, segmentation, and filtering.\n\n  - **Security Context and Cost Allocation:** Leverage existing Kubernetes namespace labels or annotations  as the basis for security context or cost allocation. This provides granular control over permissions and facilitates chargeback functionalities.\n\n  Additional Information:\n\n  - Only namespace-level labels or annotations can be used as source.\n\n  - You can define up to 20 enrichment rules.\n\n  - New rules may take up to 45 minutes to take effect.\n\n  - Pod restarts are required after the 45 mins to ensure the changes take effect.\n\n  To learn more, please refer to our [documentation](https://dt-url.net/pn22sye).
	Scope *string `json:"-" scope:"scope"` // The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rules": {
			Type:        schema.TypeList,
			Description: "Kubernetes Telemetry Enrichment empowers you to effectively tag your telemetry data using Kubernetes namespace labels and annotations. Additionally, it enables you to tag it for cost allocation and permission purposes.\n\n  Enrichment Options:\n\n  - **Enrich telemetry with label/annotation directly:** Tag your telemetry data with existing Kubernetes namespace labels or annotations. These will be made available as domain-specific fields (e.g., `k8s.namespace.label.your_key`). This allows for flexible pipeline routing, bucket selection, segmentation, and filtering.\n\n  - **Security Context and Cost Allocation:** Leverage existing Kubernetes namespace labels or annotations  as the basis for security context or cost allocation. This provides granular control over permissions and facilitates chargeback functionalities.\n\n  Additional Information:\n\n  - Only namespace-level labels or annotations can be used as source.\n\n  - You can define up to 20 enrichment rules.\n\n  - New rules may take up to 45 minutes to take effect.\n\n  - Pod restarts are required after the 45 mins to ensure the changes take effect.\n\n  To learn more, please refer to our [documentation](https://dt-url.net/pn22sye).",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(Rules).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"rules": me.Rules,
		"scope": me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"rules": &me.Rules,
		"scope": &me.Scope,
	})
}
