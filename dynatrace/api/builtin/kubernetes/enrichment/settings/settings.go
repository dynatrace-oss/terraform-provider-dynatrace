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
	Rules Rules   `json:"rules,omitempty"` // Dynatrace allows to use metadata defined on Kubernetes nodes, namespaces, and pods to set security and cost allocation attributes and dimensions for metrics, events, log, spans, and entities associated with the respective Kubernetes resource. \n\n The following annotation keys are considered: \n * `metadata.dynatrace.com/dt.security_context` \n * `metadata.dynatrace.com/dt.cost.product` \n * `metadata.dynatrace.com/dt.cost.costcenter` \n\n Pod annotations determine the attributes of data associated with the pod itself, and containers belonging to the pod. \n\nNamespace annotations determine the attributes of data associated with the namespace itself, workloads, services, and - if not overwritten on pod level - pods, and containers belonging to the namespace. \n\nNode annotations determine the attributes of data associated with only the node. \n\n Depending on your specific use case and environment, you have the following enrichment options: \n\n **Manual annotation:** \n\n Use the aforementioned annotation keys when annotating your namespaces and pods to enrich your Kubernetes data with security and cost allocation attributes.\n\nWith Dynatrace Operator version 1.3.0, the aforementioned namespace annotations are copied down to pods in the namespace, if they are not yet set on the respective pod. \n\n **Rule-based annotation:**\n\nIf you already have labels or annotations defined on your namespaces, and you want to reuse them for enrichment, you can do so with the help of rules definable here. \n\n**Example:**\n\n * Namespace label:\n   * `label/example: test-value`\n\n * Rule: \n   * `Label` \n `label/example --> dt.security_context`\n\n * Pod annotation: \n   * `metadata.dynatrace.com/dt.security_context: test-value`\n\nA maximum of 5 rules can be defined. The first applicable rule will be applied. Preexisting annotations will not be overwritten. For a detailed description of this feature, have a look at our [documentation](https://dt-url.net/pn22sye).
	Scope *string `json:"-" scope:"scope"` // The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rules": {
			Type:        schema.TypeList,
			Description: "Dynatrace allows to use metadata defined on Kubernetes nodes, namespaces, and pods to set security and cost allocation attributes and dimensions for metrics, events, log, spans, and entities associated with the respective Kubernetes resource. \n\n The following annotation keys are considered: \n * `metadata.dynatrace.com/dt.security_context` \n * `metadata.dynatrace.com/dt.cost.product` \n * `metadata.dynatrace.com/dt.cost.costcenter` \n\n Pod annotations determine the attributes of data associated with the pod itself, and containers belonging to the pod. \n\nNamespace annotations determine the attributes of data associated with the namespace itself, workloads, services, and - if not overwritten on pod level - pods, and containers belonging to the namespace. \n\nNode annotations determine the attributes of data associated with only the node. \n\n Depending on your specific use case and environment, you have the following enrichment options: \n\n **Manual annotation:** \n\n Use the aforementioned annotation keys when annotating your namespaces and pods to enrich your Kubernetes data with security and cost allocation attributes.\n\nWith Dynatrace Operator version 1.3.0, the aforementioned namespace annotations are copied down to pods in the namespace, if they are not yet set on the respective pod. \n\n **Rule-based annotation:**\n\nIf you already have labels or annotations defined on your namespaces, and you want to reuse them for enrichment, you can do so with the help of rules definable here. \n\n**Example:**\n\n * Namespace label:\n   * `label/example: test-value`\n\n * Rule: \n   * `Label` \n `label/example --> dt.security_context`\n\n * Pod annotation: \n   * `metadata.dynatrace.com/dt.security_context: test-value`\n\nA maximum of 5 rules can be defined. The first applicable rule will be applied. Preexisting annotations will not be overwritten. For a detailed description of this feature, have a look at our [documentation](https://dt-url.net/pn22sye).",
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
