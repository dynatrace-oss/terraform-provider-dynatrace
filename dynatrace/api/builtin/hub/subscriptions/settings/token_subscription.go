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

package subscriptions

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TokenSubscriptions []*TokenSubscription

func (me *TokenSubscriptions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"token_subscription": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(TokenSubscription).Schema()},
		},
	}
}

func (me TokenSubscriptions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("token_subscription", me)
}

func (me *TokenSubscriptions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("token_subscription", me)
}

type TokenSubscription struct {
	Description *string `json:"description,omitempty"`
	Enabled     bool    `json:"enabled"` // This setting is enabled (`true`) or disabled (`false`)
	Name        string  `json:"name"`    // Name of subscription
	Token       string  `json:"token"`   // Subscription token
}

func (me *TokenSubscription) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Optional:    true, // nullable
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name of subscription",
			Required:    true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "Subscription token",
			Required:    true,
		},
	}
}

func (me *TokenSubscription) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"description": me.Description,
		"enabled":     me.Enabled,
		"name":        me.Name,
		"token":       me.Token,
	})
}

func (me *TokenSubscription) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"description": &me.Description,
		"enabled":     &me.Enabled,
		"name":        &me.Name,
		"token":       &me.Token,
	})
}
