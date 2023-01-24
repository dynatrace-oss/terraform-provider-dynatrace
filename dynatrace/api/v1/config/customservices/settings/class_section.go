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

package customservices

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ClassSection struct {
	Name  *string
	Match *ClassNameMatcher
}

func (me *ClassSection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The full name of the class / the name to match the class name with",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeString,
			Description: "Matcher applying to the class name (ENDS_WITH, EQUALS or STARTS_WITH). STARTS_WITH can only be used if there is at least one annotation defined. Default value is EQUALS",
			Optional:    true,
			Default:     "EQUALS",
		},
	}
}

func (me *ClassSection) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("match", me.Match); err != nil {
		return err
	}
	return nil
}

func (me *ClassSection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("match"); ok {
		me.Match = ClassNameMatcher(value.(string)).Ref()
	}
	return nil
}
