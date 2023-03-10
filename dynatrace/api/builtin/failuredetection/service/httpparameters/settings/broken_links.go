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

type BrokenLinks struct {
	BrokenLinkDomains       []string `json:"brokenLinkDomains,omitempty"` // If your application relies on other hosts at other domains, add the associated domain names here. Once configured, Dynatrace will consider 404s thrown by hosts at these domains to be service failures related to your application.
	Http404NotFoundFailures bool     `json:"http404NotFoundFailures"`     // Consider 404 HTTP response codes as failures
}

func (me *BrokenLinks) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"broken_link_domains": {
			Type:        schema.TypeSet,
			Description: "If your application relies on other hosts at other domains, add the associated domain names here. Once configured, Dynatrace will consider 404s thrown by hosts at these domains to be service failures related to your application.",
			Optional:    true,

			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"http_404_not_found_failures": {
			Type:        schema.TypeBool,
			Description: "Consider 404 HTTP response codes as failures",
			Required:    true,
		},
	}
}

func (me *BrokenLinks) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"broken_link_domains":         me.BrokenLinkDomains,
		"http_404_not_found_failures": me.Http404NotFoundFailures,
	})
}

func (me *BrokenLinks) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"broken_link_domains":         &me.BrokenLinkDomains,
		"http_404_not_found_failures": &me.Http404NotFoundFailures,
	})
}
