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

package http

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Script struct {
	Version  string   `json:"version"`  // Script version—use the `1.0` value here
	Requests Requests `json:"requests"` // A list of HTTP requests to be performed by the monitor.\n\nThe requests are executed in the order in which they appear in the script
}

func (me *Script) GetVersion() string {
	return "1.0"
}

func (me *Script) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"request": {
			Type:        schema.TypeList,
			Description: "A HTTP request to be performed by the monitor.",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Request).Schema()},
		},
	}
}

func (me *Script) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("request", me.Requests); err != nil {
		return err
	}
	return nil
}

func (me *Script) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Version = me.GetVersion()
	if result, ok := decoder.GetOk("request.#"); ok {
		me.Requests = Requests{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Request)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "request", idx)); err != nil {
				return err
			}
			me.Requests = append(me.Requests, entry)
		}
	}
	return nil
}
