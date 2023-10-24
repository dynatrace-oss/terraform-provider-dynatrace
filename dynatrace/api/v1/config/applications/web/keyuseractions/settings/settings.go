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

package keyuseractions

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID string  `json:"-"`                // The ID of the WebApplication
	Type          Type    `json:"actionType"`       // The type of the action. Possible values are `Custom`, `Load` and `Xhr`
	Domain        *string `json:"domain,omitempty"` // The domain where the action is performed
	Name          string  `json:"name"`             // The name of the action
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The ID of the WebApplication",
			Required:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the action. Possible values are `Custom`, `Load` and `Xhr`",
			Required:    true,
			ForceNew:    true,
		},
		"domain": {
			Type:        schema.TypeString,
			Description: "The domain where the action is performed",
			Optional:    true,
			ForceNew:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the action",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": me.ApplicationID,
		"type":           me.Type,
		"domain":         me.Domain,
		"name":           me.Name,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"type":           &me.Type,
		"domain":         &me.Domain,
		"name":           &me.Name,
	})
}
