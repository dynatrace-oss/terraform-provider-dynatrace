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

package workloaddetection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KubernetesOpenShiftCAWD struct {
	Enabled bool            `json:"enabled"` // This setting is enabled (`true`) or disabled (`false`)
	Filters FilterComplexes `json:"filters"` // Define rules to merge similar Kubernetes workloads into process groups. \n\n You can use workload properties like namespace name, base pod name or container name as well as the [environment variables DT_RELEASE_STAGE and DT_RELEASE_PRODUCT](https://dt-url.net/sb02v2a) for grouping processes of similar workloads. The first applicable rule will be applied. If no rule matches, “Namespace name” + “Base pod name” + “Container name” is used as fallback.
}

func (me *KubernetesOpenShiftCAWD) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"filters": {
			Type:        schema.TypeList,
			Description: "Define rules to merge similar Kubernetes workloads into process groups. \n\n You can use workload properties like namespace name, base pod name or container name as well as the [environment variables DT_RELEASE_STAGE and DT_RELEASE_PRODUCT](https://dt-url.net/sb02v2a) for grouping processes of similar workloads. The first applicable rule will be applied. If no rule matches, “Namespace name” + “Base pod name” + “Container name” is used as fallback.",
			Required:    true,

			Elem:     &schema.Resource{Schema: new(FilterComplexes).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
	}
}

func (me *KubernetesOpenShiftCAWD) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"filters": me.Filters,
	})
}

func (me *KubernetesOpenShiftCAWD) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled": &me.Enabled,
		"filters": &me.Filters,
	})
}
