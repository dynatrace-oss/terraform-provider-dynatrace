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

package declarativegrouping

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProcessDefinitions []*ProcessDefinition

func (me *ProcessDefinitions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"process_definition": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(ProcessDefinition).Schema()},
		},
	}
}

func (me ProcessDefinitions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("process_definition", me)
}

func (me *ProcessDefinitions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("process_definition", me)
}

type ProcessDefinition struct {
	ID               string              `json:"id"`               // Process group identifier
	ProcessGroupName string              `json:"processGroupName"` // This identifier is used by Dynatrace to recognize this process group.
	Report           ReportItem          `json:"report"`           // Possible Values: `never`, `always`, `highResourceUsage`
	Rules            DetectionConditions `json:"rules,omitempty"`  // Define process detection rules by selecting a process property and a condition. Each process group can have multiple detection rules associated with it.
}

func (me *ProcessDefinition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "Process group identifier",
			Required:    true,
		},
		"process_group_name": {
			Type:        schema.TypeString,
			Description: "This identifier is used by Dynatrace to recognize this process group.",
			Required:    true,
		},
		"report": {
			Type:        schema.TypeString,
			Description: "Possible Values: `never`, `always`, `highResourceUsage`",
			Required:    true,
		},
		"rules": {
			Type:        schema.TypeList,
			Description: "Define process detection rules by selecting a process property and a condition. Each process group can have multiple detection rules associated with it.",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(DetectionConditions).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *ProcessDefinition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":                 me.ID,
		"process_group_name": me.ProcessGroupName,
		"report":             me.Report,
		"rules":              me.Rules,
	})
}

func (me *ProcessDefinition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":                 &me.ID,
		"process_group_name": &me.ProcessGroupName,
		"report":             &me.Report,
		"rules":              &me.Rules,
	})
}
