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

package anomalydetectors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HostMetadataConditionTypes []*HostMetadataConditionType

func (me *HostMetadataConditionTypes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_metadata_condition": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(HostMetadataConditionType).Schema()},
		},
	}
}

func (me HostMetadataConditionTypes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("host_metadata_condition", me)
}

func (me *HostMetadataConditionTypes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("host_metadata_condition", me)
}

type HostMetadataConditionType struct {
	HostMetadataCondition *HostMetadataCondition `json:"hostMetadataCondition"`
}

func (me *HostMetadataConditionType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_metadata_condition": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HostMetadataCondition).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *HostMetadataConditionType) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host_metadata_condition": me.HostMetadataCondition,
	})
}

func (me *HostMetadataConditionType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host_metadata_condition": &me.HostMetadataCondition,
	})
}
