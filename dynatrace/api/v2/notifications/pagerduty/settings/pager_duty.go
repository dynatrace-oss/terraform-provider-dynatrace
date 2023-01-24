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

package notifications

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PagerDuty struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	Account     string `json:"account"`       // The name of the PagerDuty account
	ServiceName string `json:"serviceName"`   // The name of the service
	APIKey      string `json:"serviceApiKey"` // The API key to access PagerDuty
}

func (me *PagerDuty) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if apiKey, ok := decoder.GetOk("api_key"); ok && len(apiKey.(string)) > 0 {
		me.APIKey = apiKey.(string)
	}
	return nil
}

func (me *PagerDuty) FillDemoValues() []string {
	me.APIKey = "#######"
	return []string{"Please fill in the API Key"}
}

func (me *PagerDuty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "The configuration is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},

		"account": {
			Type:        schema.TypeString,
			Description: "The name of the PagerDuty account",
			Required:    true,
		},
		"api_key": {
			Type:        schema.TypeString,
			Sensitive:   true,
			Description: "The API key to access PagerDuty",
			Optional:    true,
		},
		"service": {
			Type:        schema.TypeString,
			Description: "The name of the PagerDuty Service",
			Required:    true,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of these settings when referred to from resources requiring the REST API V1 keys",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *PagerDuty) MarshalHCL(properties hcl.Properties) error { // The api_key field MUST NOT get serialized into HCL here
	// The Dynatrace Settings 2.0 API delivers a scrambled version of any previously stored api_key here
	// Evaluation at this point would lead to that scrambled version to make it into the Terraform State
	// As a result any plans would be non-empty
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"active":  me.Enabled,
		"profile": me.ProfileID,

		"account": me.Account,
		"service": me.ServiceName,
		"api_key": me.APIKey,
	})
}

func (me *PagerDuty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"account": &me.Account,
		"service": &me.ServiceName,
		"api_key": &me.APIKey,
	})
}
