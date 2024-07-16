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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RequestConfigurations []*RequestConfiguration

func (me *RequestConfigurations) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"request_configuration": {
			Type:        schema.TypeSet,
			Description: "The configuration of a network availability monitor request",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(RequestConfiguration).Schema()},
		},
	}
}

func (me RequestConfigurations) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("request_configuration", me)
}

func (me *RequestConfigurations) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("request_configuration", me)
}

type RequestConfiguration struct {
	Constraints Constraints `json:"constraints,omitempty"` // Request constraints
}

func (me *RequestConfiguration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"constraints": {
			Type:        schema.TypeList,
			Description: "Request constraints",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Constraints).Schema()},
		},
	}
}

func (me RequestConfiguration) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"constraints": me.Constraints,
	})
}

func (me *RequestConfiguration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"constraints": &me.Constraints,
	})
}
