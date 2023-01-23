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

package event

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Javascript struct {
	EventBase
	Javascript string         `json:"javaScript"`       // The JavaScript code to be executed in this event
	Wait       *WaitCondition `json:"wait,omitempty"`   // The wait condition for the event—defines how long Dynatrace should wait before the next action is executed
	Target     *Target        `json:"target,omitempty"` // The tab on which the page should open
}

func (me *Javascript) GetType() Type {
	return Types.Javascript
}

func (me *Javascript) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code": {
			Type:             schema.TypeString,
			Description:      "The JavaScript code to be executed in this event",
			Required:         true,
			DiffSuppressFunc: SuppressEquivalent,
		},
		"wait": {
			Type:        schema.TypeList,
			Description: "The wait condition for the event—defines how long Dynatrace should wait before the next action is executed",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(WaitCondition).Schema()},
		},
		"target": {
			Type:        schema.TypeList,
			Description: "The tab on which the page should open",
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(Target).Schema()},
		},
	}
}

func JSONStringsEqual(s1, s2 string) bool {
	b1 := bytes.NewBufferString("")
	if err := json.Compact(b1, []byte(s1)); err != nil {
		return false
	}

	b2 := bytes.NewBufferString("")
	if err := json.Compact(b2, []byte(s2)); err != nil {
		return false
	}

	return JSONBytesEqual(b1.Bytes(), b2.Bytes())
}

func JSONBytesEqual(b1, b2 []byte) bool {
	var o1 interface{}
	if err := json.Unmarshal(b1, &o1); err != nil {
		return false
	}

	var o2 interface{}
	if err := json.Unmarshal(b2, &o2); err != nil {
		return false
	}

	return reflect.DeepEqual(o1, o2)
}

func equalLineByLine(s1, s2 string) bool {
	parts1 := strings.Split(strings.TrimSpace(s1), "\n")
	parts2 := strings.Split(strings.TrimSpace(s2), "\n")
	if len(parts1) != len(parts2) {
		return false
	}
	for idx := range parts1 {
		part1 := strings.TrimSpace(parts1[idx])
		part2 := strings.TrimSpace(parts2[idx])
		if part1 != part2 {
			return false
		}
	}
	return true
}

func SuppressEquivalent(k, old, new string, d *schema.ResourceData) bool {
	if JSONStringsEqual(old, new) {
		return true
	}
	return equalLineByLine(old, new)
}

func (me *Javascript) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("code", me.Javascript); err != nil {
		return err
	}
	if err := properties.Encode("wait", me.Wait); err != nil {
		return err
	}
	if err := properties.Encode("target", me.Target); err != nil {
		return err
	}
	return nil
}

func (me *Javascript) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Type = Types.Tap
	if err := decoder.Decode("code", &me.Javascript); err != nil {
		return err
	}
	if err := decoder.Decode("wait", &me.Wait); err != nil {
		return err
	}
	if err := decoder.Decode("target", &me.Target); err != nil {
		return err
	}
	return nil
}
