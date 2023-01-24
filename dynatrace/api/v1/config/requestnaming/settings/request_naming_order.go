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

package requestnaming

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Order struct {
	Values []Ref `json:"values"`
}

func (me *Order) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ids": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The IDs of the request namings in the order they should be taken into consideration",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *Order) MarshalHCL(properties hcl.Properties) error {
	refs := []any{}
	for _, ref := range me.Values {
		refs = append(refs, ref.ID)
	}
	if err := properties.Encode("ids", refs); err != nil {
		return err
	}
	return nil
}

func (me *Order) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Values = []Ref{}
	values, ok := decoder.GetOk("ids")
	if ok {
		vals := values.([]any)
		for _, val := range vals {
			me.Values = append(me.Values, Ref{ID: val.(string)})
		}
	}
	return nil
}

type Ref struct {
	ID string `json:"id"`
}
