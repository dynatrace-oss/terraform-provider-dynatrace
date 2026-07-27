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

package slackconnection

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ExternalApproval *bool   `json:"externalApproval,omitempty"` // Accept external approvals can enable Slack users to directly respond to approval request.
	Name             string  `json:"name"`                       // Provide a unique and clearly identifiable connection name to your Slack App.
	SigningSecret    *string `json:"signingSecret,omitempty"`    // The signing secret obtained from the Slack App Management UI.
	Token            string  `json:"token"`                      // The bot token obtained from the Slack App Management UI.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"external_approval": {
			Type:        schema.TypeBool,
			Description: "Accept external approvals can enable Slack users to directly respond to approval request.",
			Optional:    true, // nullable
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Provide a unique and clearly identifiable connection name to your Slack App.",
			Required:    true,
		},
		"signing_secret": {
			Type:        schema.TypeString,
			Description: "The signing secret obtained from the Slack App Management UI.",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "The bot token obtained from the Slack App Management UI.",
			Required:    true,
			Sensitive:   true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"external_approval": me.ExternalApproval,
		"name":              me.Name,
		"signing_secret":    "${state.secret_value}",
		"token":             "${state.secret_value}",
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.SigningSecret != nil) && (me.ExternalApproval == nil || !*me.ExternalApproval) {
		return fmt.Errorf("'signing_secret' must not be specified unless 'external_approval' is set to 'true'; got 'external_approval'='%v'", opt.ValOrNil(me.ExternalApproval))
	}
	if (me.SigningSecret == nil) && (me.ExternalApproval != nil && *me.ExternalApproval) {
		return fmt.Errorf("'signing_secret' must be specified when 'external_approval' is set to 'true'; got 'external_approval'='%v'", opt.ValOrNil(me.ExternalApproval))
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"external_approval": &me.ExternalApproval,
		"name":              &me.Name,
		"signing_secret":    &me.SigningSecret,
		"token":             &me.Token,
	})
}

func (me *Settings) FillDemoValues() []string {
	me.Token = "#######"
	me.SigningSecret = new("#######")
	return []string{"REST API didn't provide token data"}
}
