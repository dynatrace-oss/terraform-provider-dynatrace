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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KubernetesAppOptions struct {
	EnableKubernetesApp bool `json:"enableKubernetesApp"` // New Kubernetes experience
}

func (me *KubernetesAppOptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_kubernetes_app": {
			Type:        schema.TypeBool,
			Description: "New Kubernetes experience",
			Required:    true,
		},
	}
}

func (me *KubernetesAppOptions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enable_kubernetes_app": me.EnableKubernetesApp,
	})
}

func (me *KubernetesAppOptions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enable_kubernetes_app": &me.EnableKubernetesApp,
	})
}
