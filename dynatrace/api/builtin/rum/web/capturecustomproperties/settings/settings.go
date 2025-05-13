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

package capturecustomproperties

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID                    string           `json:"-" scope:"applicationId"`                    // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	CustomEventPropertiesAllowList   CustomProperties `json:"customEventPropertiesAllowList,omitempty"`   // List of allowed custom event properties
	CustomSessionPropertiesAllowList CustomProperties `json:"customSessionPropertiesAllowList,omitempty"` // List of allowed custom session properties
}

func (me *Settings) Name() string {
	return me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
		},
		"custom_event_properties_allow_list": {
			Type:        schema.TypeList,
			Description: "List of allowed custom event properties",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(CustomProperties).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"custom_session_properties_allow_list": {
			Type:        schema.TypeList,
			Description: "List of allowed custom session properties",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(CustomProperties).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id":                       me.ApplicationID,
		"custom_event_properties_allow_list":   me.CustomEventPropertiesAllowList,
		"custom_session_properties_allow_list": me.CustomSessionPropertiesAllowList,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id":                       &me.ApplicationID,
		"custom_event_properties_allow_list":   &me.CustomEventPropertiesAllowList,
		"custom_session_properties_allow_list": &me.CustomSessionPropertiesAllowList,
	})
}
