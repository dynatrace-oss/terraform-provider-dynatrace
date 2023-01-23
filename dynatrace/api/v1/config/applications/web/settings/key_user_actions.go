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

package web

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KeyUserActionList struct {
	KeyUserActions KeyUserActions `json:"keyUserActionList,omitempty"` // The list of key user actions in the web application
}

type KeyUserActions []*KeyUserAction

func (me KeyUserActions) Equals(other any) bool {
	if other == nil {
		return false
	}
	o, ok := other.(KeyUserActions)
	if !ok {
		return false
	}
	if len(me) != len(o) {
		return false
	}
	for _, kua := range me {
		found := false
		for _, kuo := range o {
			if kua.Equals(kuo) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (me *KeyUserActions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:        schema.TypeList,
			Description: "Configuration of the key user action",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(KeyUserAction).Schema()},
		},
	}
}

func (me KeyUserActions) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("action", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *KeyUserActions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("action", me)
}

// KeyUserAction represents configuration of the key user action
type KeyUserAction struct {
	Name   string            `json:"name"`             // The name of the action
	Type   KeyUserActionType `json:"actionType"`       // The type of the action. Possible values are `Custom`, `Load` and `Xhr`.
	Domain *string           `json:"domain,omitempty"` // The domain where the action is performed
}

func (me *KeyUserAction) Equals(other any) bool {
	if other == nil {
		return false
	}
	ot, ok := other.(*KeyUserAction)
	if !ok {
		return false
	}
	if me.Name != ot.Name {
		return false
	}
	if string(me.Type) != string(ot.Type) {
		return false
	}
	if me.Domain == nil && ot.Domain != nil {
		return false
	}
	if me.Domain != nil && ot.Domain == nil {
		return false
	}
	if me.Domain != nil && ot.Domain != nil && *me.Domain != *ot.Domain {
		return false
	}
	return true
}

func (me *KeyUserAction) String() string {
	tmp := struct {
		Name   string            `json:"name"`
		Type   KeyUserActionType `json:"actionType"`
		Domain *string           `json:"domain,omitempty"`
	}{
		me.Name,
		me.Type,
		me.Domain,
	}
	data, _ := json.Marshal(tmp)
	return string(data)
}

func (me *KeyUserAction) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the action",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the action. Possible values are `Custom`, `Load` and `Xhr`.",
			Required:    true,
		},
		"domain": {
			Type:        schema.TypeString,
			Description: "The domain where the action is performed.",
			Optional:    true,
		},
	}
}

func (me *KeyUserAction) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":   me.Name,
		"type":   me.Type,
		"domain": me.Domain,
	})
}

func (me *KeyUserAction) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":   &me.Name,
		"type":   &me.Type,
		"domain": &me.Domain,
	})
}
