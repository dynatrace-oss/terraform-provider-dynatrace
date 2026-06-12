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

type GroupIdSource struct {
	AdvancedSettings *AdvancedSettings `json:"advancedSettings,omitempty"` // Set advanced options to customize delimiters and control how property values are processed.. Consider an environment with processes such as:\n  * `python myScript.py --env=prod12 --id=12`\n  * `python myScript.py --env=dev2 --id=2`\n  * etc.\n\n   To group production *(prod)* and development *(dev)* processes together you could use Command line property with:\n  * **Delimiter** from `--env=` to `--id` to extract `prod12 ` and `dev2 `\n  * Enable **Ignore numbers** to transform `prod12` to `prod*` and `dev2` to `dev*`.
	ID               *string           `json:"id,omitempty"`               // This identifier is used by Dynatrace to recognize this process group.
	Name             *string           `json:"name,omitempty"`             // If Dynatrace detects this property at startup of a process, it will use its value to identify process groups.
	Property         *string           `json:"property,omitempty"`         // 3.1.2. Property
	StandaloneRule   bool              `json:"standaloneRule"`             // Valid only for **deep monitored** processes.. If this option is selected, the default Dynatrace behavior is disabled for the detected processes. Only this rule is used to separate the process group.\n\n  If this option is not selected, this rule contributes to the default Dynatrace process group detection. \n\n  [See our help page for examples.](https://dt-url.net/1722wrz)
	Type             IdType            `json:"type"`                       // Pick which property should be used to identify your process group. You can pick a custom variable or pick an existing process property. Possible values: `CUSTOM`, `EXISTING`
}

func (me *GroupIdSource) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"advanced_settings": {
			Type:        schema.TypeList,
			Description: "Set advanced options to customize delimiters and control how property values are processed.. Consider an environment with processes such as:\n  * `python myScript.py --env=prod12 --id=12`\n  * `python myScript.py --env=dev2 --id=2`\n  * etc.\n\n   To group production *(prod)* and development *(dev)* processes together you could use Command line property with:\n  * **Delimiter** from `--env=` to `--id` to extract `prod12 ` and `dev2 `\n  * Enable **Ignore numbers** to transform `prod12` to `prod*` and `dev2` to `dev*`.",
			Optional:    true, // nullable & precondition
			Elem:        &schema.Resource{Schema: new(AdvancedSettings).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "This identifier is used by Dynatrace to recognize this process group.",
			Optional:    true, // precondition
		},
		"name": {
			Type:        schema.TypeString,
			Description: "If Dynatrace detects this property at startup of a process, it will use its value to identify process groups.",
			Optional:    true, // precondition
		},
		"property": {
			Type:        schema.TypeString,
			Description: "3.1.2. Property",
			Optional:    true, // precondition
		},
		"standalone_rule": {
			Type:        schema.TypeBool,
			Description: "Valid only for **deep monitored** processes.. If this option is selected, the default Dynatrace behavior is disabled for the detected processes. Only this rule is used to separate the process group.\n\n  If this option is not selected, this rule contributes to the default Dynatrace process group detection. \n\n  [See our help page for examples.](https://dt-url.net/1722wrz)",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Pick which property should be used to identify your process group. You can pick a custom variable or pick an existing process property. Possible values: `CUSTOM`, `EXISTING`",
			Required:    true,
		},
	}
}

func (me *GroupIdSource) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"advanced_settings": me.AdvancedSettings,
		"id":                me.ID,
		"name":              me.Name,
		"property":          me.Property,
		"standalone_rule":   me.StandaloneRule,
		"type":              me.Type,
	})
}

func (me *GroupIdSource) HandlePreconditions() error {
	if (me.AdvancedSettings != nil) && (string(me.Type) != "EXISTING") {
		return fmt.Errorf("'advanced_settings' must not be specified unless 'type' is set to 'EXISTING'; got 'type'='%v'", me.Type)
	}
	if (me.ID != nil) && (string(me.Type) != "CUSTOM") {
		return fmt.Errorf("'id' must not be specified unless 'type' is set to 'CUSTOM'; got 'type'='%v'", me.Type)
	}
	if (me.ID == nil) && (string(me.Type) == "CUSTOM") {
		return fmt.Errorf("'id' must be specified when 'type' is set to 'CUSTOM'; got 'type'='%v'", me.Type)
	}
	if (me.Name != nil) && (me.Property == nil || !slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(*me.Property))) {
		return fmt.Errorf("'name' must not be specified unless 'property' is one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	if (me.Name == nil) && (me.Property != nil && slices.Contains([]string{"ENVIRONMENT_VARIABLE", "JAVA_SYS_PROP"}, string(*me.Property))) {
		return fmt.Errorf("'name' must be specified when 'property' is one of ['ENVIRONMENT_VARIABLE', 'JAVA_SYS_PROP']; got 'property'='%v'", opt.ValOrNil(me.Property))
	}
	if (me.Property != nil) && (string(me.Type) != "EXISTING") {
		return fmt.Errorf("'property' must not be specified unless 'type' is set to 'EXISTING'; got 'type'='%v'", me.Type)
	}
	if (me.Property == nil) && (string(me.Type) == "EXISTING") {
		return fmt.Errorf("'property' must be specified when 'type' is set to 'EXISTING'; got 'type'='%v'", me.Type)
	}
	return nil
}

func (me *GroupIdSource) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"advanced_settings": &me.AdvancedSettings,
		"id":                &me.ID,
		"name":              &me.Name,
		"property":          &me.Property,
		"standalone_rule":   &me.StandaloneRule,
		"type":              &me.Type,
	})
}
