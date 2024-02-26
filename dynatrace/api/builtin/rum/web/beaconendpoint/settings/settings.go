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

package beaconendpoint

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID string                `json:"-" scope:"applicationId"` // The scope of this setting
	Type          WebBeaconEndpointType `json:"type"`                    // Possible Values: `ACTIVEGATE`, `DEFAULT_CONFIG`, `ONEAGENT`
	Url           *string               `json:"url,omitempty"`           // You can specify either path segments or an absolute URL.
	UseCors       *bool                 `json:"useCors,omitempty"`       // Learn more about [sending beacon data via CORS](https://dt-url.net/r7038sa)
}

func (me *Settings) Name() string {
	return me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this setting",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ACTIVEGATE`, `DEFAULT_CONFIG`, `ONEAGENT`",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "You can specify either path segments or an absolute URL.",
			Optional:    true, // precondition
		},
		"use_cors": {
			Type:        schema.TypeBool,
			Description: "Learn more about [sending beacon data via CORS](https://dt-url.net/r7038sa)",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": me.ApplicationID,
		"type":           me.Type,
		"url":            me.Url,
		"use_cors":       me.UseCors,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Url == nil) && (string(me.Type) == "ONEAGENT") {
		me.Url = opt.NewString("")
	}
	if (me.UseCors == nil) && (string(me.Type) == "ONEAGENT") {
		me.UseCors = opt.NewBool(false)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"type":           &me.Type,
		"url":            &me.Url,
		"use_cors":       &me.UseCors,
	})
}
