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

package securityposturemanagement

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ConfigurationDatasetPipelineEnabled bool    `json:"configurationDatasetPipelineEnabled"` // Follow the [installation instructions](https://dt-url.net/4x23ut5) to deploy the Security Posture Management components.
	Scope                               *string `json:"-" scope:"scope"`                     // The scope of this setting (KUBERNETES_CLUSTER). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"configuration_dataset_pipeline_enabled": {
			Type:        schema.TypeBool,
			Description: "Follow the [installation instructions](https://dt-url.net/4x23ut5) to deploy the Security Posture Management components.",
			Required:    true,
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
		"configuration_dataset_pipeline_enabled": me.ConfigurationDatasetPipelineEnabled,
		"scope":                                  me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"configuration_dataset_pipeline_enabled": &me.ConfigurationDatasetPipelineEnabled,
		"scope":                                  &me.Scope,
	})
}
