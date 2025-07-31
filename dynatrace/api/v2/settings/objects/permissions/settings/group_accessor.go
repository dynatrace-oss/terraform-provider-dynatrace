/**
* @license
* Copyright 2025 Dynatrace LLC
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

package settings

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Groups []*GroupAccessor

func (_ *Groups) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "Group that is to be granted read or write permissions",
			Elem:        &schema.Resource{Schema: new(GroupAccessor).Schema()},
		},
	}
}

func (g *Groups) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("group", g)
}

func (g *Groups) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("group", g)
}

type GroupAccessor struct {
	ID     string
	Access string
}

func (_ *GroupAccessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The UUID of the group, conveniently retrieved via the `id` attribute provided by the data source `dynatrace_iam_group",
			Required:    true,
		},
		"access": {
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("Valid values: `%s`, `%s`", HCLAccessorRead, HCLAccessorWrite),
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{HCLAccessorRead, HCLAccessorWrite}, false),
		},
	}
}

func (ga *GroupAccessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":     ga.ID,
		"access": ga.Access,
	})
}

func (ga *GroupAccessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":     &ga.ID,
		"access": &ga.Access,
	})
}
