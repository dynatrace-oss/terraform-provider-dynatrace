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

package appdetection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ApplicationID string  `json:"applicationId"`         // Select an existing application or create a new one.
	Description   *string `json:"description,omitempty"` // (v1.274) Add a description for your rule
	Matcher       Matcher `json:"matcher"`               // Possible Values: `DOMAIN_CONTAINS`, `DOMAIN_ENDS_WITH`, `DOMAIN_EQUALS`, `DOMAIN_MATCHES`, `DOMAIN_STARTS_WITH`, `URL_CONTAINS`, `URL_ENDS_WITH`, `URL_EQUALS`, `URL_STARTS_WITH`
	Pattern       string  `json:"pattern"`               // Pattern
	InsertAfter   string  `json:"-"`
}

func (me *Settings) Name() string {
	return me.ApplicationID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_id": {
			Type:        schema.TypeString,
			Description: "Select an existing application or create a new one.",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "(v1.274) Add a description for your rule",
			Optional:    true, // nullable
		},
		"matcher": {
			Type:        schema.TypeString,
			Description: "Possible Values: `DOMAIN_CONTAINS`, `DOMAIN_ENDS_WITH`, `DOMAIN_EQUALS`, `DOMAIN_MATCHES`, `DOMAIN_STARTS_WITH`, `URL_CONTAINS`, `URL_ENDS_WITH`, `URL_EQUALS`, `URL_STARTS_WITH`",
			Required:    true,
		},
		"pattern": {
			Type:        schema.TypeString,
			Description: "Pattern",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"application_id": me.ApplicationID,
		"description":    me.Description,
		"matcher":        me.Matcher,
		"pattern":        me.Pattern,
		"insert_after":   me.InsertAfter,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"application_id": &me.ApplicationID,
		"description":    &me.Description,
		"matcher":        &me.Matcher,
		"pattern":        &me.Pattern,
		"insert_after":   &me.InsertAfter,
	})
}
