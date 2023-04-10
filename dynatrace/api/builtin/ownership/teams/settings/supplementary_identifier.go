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

package teams

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SupplementaryIdentifiers []*SupplementaryIdentifier

func (me *SupplementaryIdentifiers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"supplementary_identifier": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(SupplementaryIdentifier).Schema()},
		},
	}
}

func (me SupplementaryIdentifiers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("supplementary_identifier", me)
}

func (me *SupplementaryIdentifiers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("supplementary_identifier", me)
}

type SupplementaryIdentifier struct {
	SupplementaryIdentifier string `json:"supplementaryIdentifier"` // Supplementary Identifier
}

func (me *SupplementaryIdentifier) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"supplementary_identifier": {
			Type:        schema.TypeString,
			Description: "Supplementary Identifier",
			Required:    true,
		},
	}
}

func (me *SupplementaryIdentifier) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"supplementary_identifier": me.SupplementaryIdentifier,
	})
}

func (me *SupplementaryIdentifier) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"supplementary_identifier": &me.SupplementaryIdentifier,
	})
}
