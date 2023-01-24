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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// MethodRule TODO: documentation
type MethodRule struct {
	MethodName    string                     `json:"methodName"`           // The method to instrument
	ArgumentTypes []string                   `json:"argumentTypes"`        // Fully qualified types of argument the method expects
	ReturnType    string                     `json:"returnType"`           // Fully qualified type the method returns
	Visibility    *Visibility                `json:"visibility,omitempty"` // The visibility of the method rule
	Modifiers     []Modifier                 `json:"modifiers,omitempty"`  // The modifiers of the method rule
	Unknowns      map[string]json.RawMessage `json:"-"`
}

func (me *MethodRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of the method rule",
			Computed:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The method to instrument",
			Required:    true,
		},
		"returns": {
			Type:        schema.TypeString,
			Description: "Fully qualified type the method returns",
			Optional:    true,
		},
		"arguments": {
			Type:        schema.TypeList,
			Description: "Fully qualified types of argument the method expects",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"modifiers": {
			Type:        schema.TypeList,
			Description: "The modifiers of the method rule. Possible values are `ABSTRACT`, `EXTERN`, `FINAL`, `NATIVE` and `STATIC`",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"visibility": {
			Type:        schema.TypeString,
			Description: "The visibility of the method rule. Possible values are `INTERNAL`, `PACKAGE_PROTECTED`, `PRIVATE`, `PROTECTED` and `PUBLIC`",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *MethodRule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	// if me.ID != nil {
	// 	if err := properties.Encode("id", me.ID); err != nil { return err }
	// }
	if err := properties.Encode("name", me.MethodName); err != nil {
		return err
	}
	if err := properties.Encode("returns", me.ReturnType); err != nil {
		return err
	}
	if err := properties.Encode("arguments", me.ArgumentTypes); err != nil {
		return err
	}
	if err := properties.Encode("modifiers", me.Modifiers); err != nil {
		return err
	}
	if err := properties.Encode("visibility", me.Visibility); err != nil {
		return err
	}
	return nil
}

func (me *MethodRule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "id")
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "returns")
		delete(me.Unknowns, "arguments")
		delete(me.Unknowns, "modifiers")
		delete(me.Unknowns, "visibility")
		delete(me.Unknowns, "returnType")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.MethodName = value.(string)
	}
	if value, ok := decoder.GetOk("arguments"); ok {
		me.ArgumentTypes = []string{}
		for _, v := range value.([]any) {
			me.ArgumentTypes = append(me.ArgumentTypes, v.(string))
		}
	}
	if value, ok := decoder.GetOk("modifiers"); ok {
		me.Modifiers = []Modifier{}
		for _, v := range value.([]any) {
			me.Modifiers = append(me.Modifiers, Modifier(v.(string)))
		}
		if len(me.Modifiers) == 0 {
			me.Modifiers = nil
		}
	}
	if value, ok := decoder.GetOk("visibility"); ok && len(value.(string)) > 0 {
		me.Visibility = Visibility(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("returns"); ok && len(value.(string)) > 0 {
		me.ReturnType = value.(string)
	}

	return nil
}

func (me *MethodRule) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	delete(m, "id")
	if err := m.Marshal("methodName", me.MethodName); err != nil {
		return nil, err
	}
	if err := m.Marshal("argumentTypes", me.ArgumentTypes); err != nil {
		return nil, err
	}
	if err := m.Marshal("returnType", me.ReturnType); err != nil {
		return nil, err
	}
	if err := m.Marshal("modifiers", me.Modifiers); err != nil {
		return nil, err
	}
	if me.Visibility != nil && len(*me.Visibility) > 0 {
		if err := m.Marshal("visibility", me.Visibility); err != nil {
			return nil, err
		}
	}
	return json.Marshal(m)
}

func (me *MethodRule) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "id")
	if err := m.Unmarshal("methodName", &me.MethodName); err != nil {
		return err
	}
	if err := m.Unmarshal("returnType", &me.ReturnType); err != nil {
		return err
	}
	if err := m.Unmarshal("argumentTypes", &me.ArgumentTypes); err != nil {
		return err
	}
	if err := m.Unmarshal("modifiers", &me.Modifiers); err != nil {
		return err
	}
	if err := m.Unmarshal("visibility", &me.Visibility); err != nil {
		return err
	}
	delete(m, "id")
	delete(m, "methodName")
	delete(m, "returnType")
	delete(m, "argumentTypes")
	delete(m, "modifiers")
	delete(m, "visibility")
	delete(m, "returnType")

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
