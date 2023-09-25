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

package notificationintegration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WebhookConfiguration struct {
	AcceptAnyCertificate bool                        `json:"acceptAnyCertificate"` // Accept any SSL certificate (including self-signed and invalid certificates)
	Headers              WebhookConfigurationHeaders `json:"headers,omitempty"`    // Use additional HTTP headers to attach any additional information, for example, configuration, authorization, or metadata.  \n  \nNote that JSON-based webhook endpoints require the addition of the **Content-Type: application/json** header to enable escaping of special characters and to avoid malformed JSON content.
	Url                  string                      `json:"url"`                  // Webhook endpoint URL
}

func (me *WebhookConfiguration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"accept_any_certificate": {
			Type:        schema.TypeBool,
			Description: "Accept any SSL certificate (including self-signed and invalid certificates)",
			Required:    true,
		},
		"headers": {
			Type:        schema.TypeList,
			Description: "Use additional HTTP headers to attach any additional information, for example, configuration, authorization, or metadata.  \n  \nNote that JSON-based webhook endpoints require the addition of the **Content-Type: application/json** header to enable escaping of special characters and to avoid malformed JSON content.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(WebhookConfigurationHeaders).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "Webhook endpoint URL",
			Required:    true,
		},
	}
}

func (me *WebhookConfiguration) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"accept_any_certificate": me.AcceptAnyCertificate,
		"headers":                me.Headers,
		"url":                    me.Url,
	})
}

func (me *WebhookConfiguration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"accept_any_certificate": &me.AcceptAnyCertificate,
		"headers":                &me.Headers,
		"url":                    &me.Url,
	})
}
