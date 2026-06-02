/**
* @license
* Copyright 2026 Dynatrace LLC
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

package processgroupingrules

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type InstanceIdSource struct {
	AdvancedSettings *AdvancedSettings `json:"advancedSettings,omitempty"` // Set advanced options to customize delimiters and control how property values are processed.. Consider an environment with processes such as:\n  * `python myScript.py --env=prod12 --id=12`\n  * `python myScript.py --env=dev2 --id=2`\n  * etc.\n\n   To group production *(prod)* and development *(dev)* processes together you could use Command line property with:\n  * **Delimiter** from `--env=` to `--id` to extract `prod12 ` and `dev2 `\n  * Enable **Ignore numbers** to transform `prod12` to `prod*` and `dev2` to `dev*`.
	Name             *string           `json:"name,omitempty"`             // If Dynatrace detects this property at startup of a process, it will use its value to identify process groups more granular.
	Property         *string           `json:"property,omitempty"`         // 3.2.1. Property
}

func (me *InstanceIdSource) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"advanced_settings": {
			Type:        schema.TypeList,
			Description: "Set advanced options to customize delimiters and control how property values are processed.. Consider an environment with processes such as:\n  * `python myScript.py --env=prod12 --id=12`\n  * `python myScript.py --env=dev2 --id=2`\n  * etc.\n\n   To group production *(prod)* and development *(dev)* processes together you could use Command line property with:\n  * **Delimiter** from `--env=` to `--id` to extract `prod12 ` and `dev2 `\n  * Enable **Ignore numbers** to transform `prod12` to `prod*` and `dev2` to `dev*`.",
			Optional:    true, // nullable
			Elem:        &schema.Resource{Schema: new(AdvancedSettings).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "If Dynatrace detects this property at startup of a process, it will use its value to identify process groups more granular.",
			Optional:    true, // precondition
		},
		"property": {
			Type:        schema.TypeString,
			Description: "3.2.1. Property",
			Optional:    true, // nullable
		},
	}
}

func (me *InstanceIdSource) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"advanced_settings": me.AdvancedSettings,
		"name":              me.Name,
		"property":          me.Property,
	})
}

func (me *InstanceIdSource) HandlePreconditions() error {
	if (me.Name != nil) && (me.Property == nil || !slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(*me.Property))) {
		return fmt.Errorf("'name' must not be specified unless 'property' is one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	if (me.Name == nil) && (me.Property != nil && slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(*me.Property))) {
		return fmt.Errorf("'name' must be specified when 'property' is one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	return nil
}

func (me *InstanceIdSource) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"advanced_settings": &me.AdvancedSettings,
		"name":              &me.Name,
		"property":          &me.Property,
	})
}
