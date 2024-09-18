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

package awsconnection

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name        string       `json:"name"` // Name
	Type        Type         `json:"type"` // Possible Values: `WebIdentity`
	WebIdentity *WebIdentity `json:"webIdentity,omitempty"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `WebIdentity`",
			Required:    true,
			ForceNew:    true,
		},
		"web_identity": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(WebIdentity).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.Name,
		"type":         me.Type,
		"web_identity": me.WebIdentity,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.WebIdentity == nil) && (string(me.Type) == "webIdentity") {
		return fmt.Errorf("'web_identity' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.WebIdentity != nil) && (string(me.Type) != "webIdentity") {
		return fmt.Errorf("'web_identity' must not be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"type":         &me.Type,
		"web_identity": &me.WebIdentity,
	})
}
