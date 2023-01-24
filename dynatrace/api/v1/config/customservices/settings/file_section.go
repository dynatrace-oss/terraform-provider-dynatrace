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

type FileSection struct {
	Name  *string
	Match *FileNameMatcher
}

func (me *FileSection) IsEmpty() bool {
	if me.Name != nil && len(*me.Name) > 0 {
		return false
	}
	if me.Match != nil && len(*me.Match) > 0 {
		return false
	}
	return true
}

func (me *FileSection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The full name of the file / the name to match the file name with",
			Required:    true,
		},
		"match": {
			Type:        schema.TypeString,
			Description: "Matcher applying to the file name (ENDS_WITH, EQUALS or STARTS_WITH). Default value is ENDS_WITH (if applicable)",
			Optional:    true,
		},
	}
}

func (me *FileSection) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("match", me.Match); err != nil {
		return err
	}
	return nil
}

func (me *FileSection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("match"); ok {
		me.Match = FileNameMatcher(value.(string)).Ref()
	}
	return nil
}
