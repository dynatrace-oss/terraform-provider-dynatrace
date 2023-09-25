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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WebhookConfigurationHeaders []*WebhookConfigurationHeader

func (me *WebhookConfigurationHeaders) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"header": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(WebhookConfigurationHeader).Schema()},
		},
	}
}

func (me WebhookConfigurationHeaders) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("header", me)
}

func (me *WebhookConfigurationHeaders) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("header", me)
}

type WebhookConfigurationHeader struct {
	Name        string  `json:"name"`
	Secret      bool    `json:"secret"`                // Secret HTTP header value
	SecretValue *string `json:"secretValue,omitempty"` // The secret value of the HTTP header. May contain an empty value.
	Value       *string `json:"value,omitempty"`       // The value of the HTTP header. May contain an empty value.
}

func (me *WebhookConfigurationHeader) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
		"secret": {
			Type:        schema.TypeBool,
			Description: "Secret HTTP header value",
			Required:    true,
		},
		"secret_value": {
			Type:        schema.TypeString,
			Description: "The secret value of the HTTP header. May contain an empty value.",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the HTTP header. May contain an empty value.",
			Optional:    true, // precondition
		},
	}
}

func (me *WebhookConfigurationHeader) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":         me.Name,
		"secret":       me.Secret,
		"secret_value": me.SecretValue,
		"value":        me.Value,
	})
}

func (me *WebhookConfigurationHeader) HandlePreconditions() error {
	if (me.SecretValue == nil) && (me.Secret) {
		me.SecretValue = opt.NewString("")
	}
	if (me.Value == nil) && (!me.Secret) {
		me.Value = opt.NewString("")
	}
	return nil
}

func (me *WebhookConfigurationHeader) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":         &me.Name,
		"secret":       &me.Secret,
		"secret_value": &me.SecretValue,
		"value":        &me.Value,
	})
}
