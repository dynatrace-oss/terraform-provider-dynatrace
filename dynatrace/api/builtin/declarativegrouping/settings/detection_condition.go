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

type DetectionConditions []*DetectionCondition

func (me *DetectionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DetectionCondition).Schema()},
		},
	}
}

func (me DetectionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *DetectionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type DetectionCondition struct {
	Condition string      `json:"condition"` // - $contains(svc) – Matches if svc appears anywhere in the process property value.\n- $eq(svc.exe) – Matches if svc.exe matches the process property value exactly.\n- $prefix(svc) – Matches if app matches the prefix of the process property value.\n- $suffix(svc.py) – Matches if svc.py matches the suffix of the process property value.\n\nFor example, $suffix(svc.py) would detect processes named loyaltysvc.py and paymentssvc.py.\n\nFor more details, see [Declarative process grouping](https://dt-url.net/j142w57).
	Property  ProcessItem `json:"property"`  // Possible Values: `Executable`, `ExecutablePath`, `CommandLine`
}

func (me *DetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"condition": {
			Type:        schema.TypeString,
			Description: "- $contains(svc) – Matches if svc appears anywhere in the process property value.\n- $eq(svc.exe) – Matches if svc.exe matches the process property value exactly.\n- $prefix(svc) – Matches if app matches the prefix of the process property value.\n- $suffix(svc.py) – Matches if svc.py matches the suffix of the process property value.\n\nFor example, $suffix(svc.py) would detect processes named loyaltysvc.py and paymentssvc.py.\n\nFor more details, see [Declarative process grouping](https://dt-url.net/j142w57).",
			Required:    true,
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Possible Values: `Executable`, `ExecutablePath`, `CommandLine`",
			Required:    true,
		},
	}
}

func (me *DetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"condition": me.Condition,
		"property":  me.Property,
	})
}

func (me *DetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"condition": &me.Condition,
		"property":  &me.Property,
	})
}
