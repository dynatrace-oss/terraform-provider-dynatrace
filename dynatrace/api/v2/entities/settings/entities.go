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

package entities

import (
	entity "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/entity/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Entities []*entity.Entity // A list of monitored entities.

func (me Entities) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity": {
			Type:        schema.TypeList,
			Description: "A list of monitored entities.",
			Elem:        &schema.Resource{Schema: new(entity.Entity).Schema()},
			Optional:    true,
		},
	}
}

func (me *Entities) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("entity", me); err != nil {
		return err
	}
	return nil
}

func (me *Entities) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("entity", me)
}
