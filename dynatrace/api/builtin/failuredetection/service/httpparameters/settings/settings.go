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

package httpparameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	BrokenLinks       *BrokenLinks       `json:"brokenLinks,omitempty"`       // HTTP 404 response codes are thrown when a web server can't find a certain page. 404s are classified as broken links on the client side and therefore aren't considered to be service failures. By enabling this setting, you can have 404s treated as server-side service failures.
	Enabled           bool               `json:"enabled"`                     // This setting is enabled (`true`) or disabled (`false`)
	HttpResponseCodes *HttpResponseCodes `json:"httpResponseCodes,omitempty"` // HTTP response codes
	ServiceID         string             `json:"-" scope:"serviceId"`         // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Name() string {
	return me.ServiceID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"broken_links": {
			Type:        schema.TypeList,
			Description: "HTTP 404 response codes are thrown when a web server can't find a certain page. 404s are classified as broken links on the client side and therefore aren't considered to be service failures. By enabling this setting, you can have 404s treated as server-side service failures.",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(BrokenLinks).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"http_response_codes": {
			Type:        schema.TypeList,
			Description: "HTTP response codes",
			Optional:    true,

			Elem:     &schema.Resource{Schema: new(HttpResponseCodes).Schema()},
			MinItems: 1,
			MaxItems: 1,
		},
		"service_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"broken_links":        me.BrokenLinks,
		"enabled":             me.Enabled,
		"http_response_codes": me.HttpResponseCodes,
		"service_id":          me.ServiceID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"broken_links":        &me.BrokenLinks,
		"enabled":             &me.Enabled,
		"http_response_codes": &me.HttpResponseCodes,
		"service_id":          &me.ServiceID,
	})
}
