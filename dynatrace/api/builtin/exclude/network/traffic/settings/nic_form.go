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

package traffic

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NicForms []*NicForm

func (me *NicForms) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"nic_form": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(NicForm).Schema()},
		},
	}
}

func (me NicForms) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("nic_form", me)
}

func (me *NicForms) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("nic_form", me)
}

type NicForm struct {
	Interface string     `json:"interface"` // Network interface
	Os        OsTypeEnum `json:"os"`        // Possible Values: `OS_TYPE_AIX`, `OS_TYPE_DARWIN`, `OS_TYPE_HPUX`, `OS_TYPE_LINUX`, `OS_TYPE_SOLARIS`, `OS_TYPE_UNKNOWN`, `OS_TYPE_WINDOWS`, `OS_TYPE_ZOS`
}

func (me *NicForm) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"interface": {
			Type:        schema.TypeString,
			Description: "Network interface",
			Required:    true,
		},
		"os": {
			Type:        schema.TypeString,
			Description: "Possible Values: `OS_TYPE_AIX`, `OS_TYPE_DARWIN`, `OS_TYPE_HPUX`, `OS_TYPE_LINUX`, `OS_TYPE_SOLARIS`, `OS_TYPE_UNKNOWN`, `OS_TYPE_WINDOWS`, `OS_TYPE_ZOS`",
			Required:    true,
		},
	}
}

func (me *NicForm) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"interface": me.Interface,
		"os":        me.Os,
	})
}

func (me *NicForm) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"interface": &me.Interface,
		"os":        &me.Os,
	})
}
