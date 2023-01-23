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

type Trello struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	ApplicationKey     string `json:"applicationKey"`     // The application key for the Trello account.\n\nYou must be logged into Trello to have Trello automatically generate an application key for you. [Get application key](https://trello.com/app-key)
	AuthorizationToken string `json:"authorizationToken"` // The authorization token for the Trello account
	BoardID            string `json:"boardId"`            // Trello board ID problem cards should be assigned to
	ListID             string `json:"listId"`             // Trello list ID new problem cards should be assigned to
	ResolvedListID     string `json:"resolvedListId"`     // Trello list ID resolved problem cards should be assigned to
	Text               string `json:"text"`               // The card text and problem placeholders to appear on new problem cards. Type '{' for placeholder suggestions
	Description        string `json:"description"`        // The description of the Trello card. Type '{' for placeholder suggestions
}

func (me *Trello) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if authorizationToken, ok := decoder.GetOk("authorization_token"); ok && len(authorizationToken.(string)) > 0 {
		me.AuthorizationToken = authorizationToken.(string)
	}
	return nil
}

func (me *Trello) FillDemoValues() []string {
	me.AuthorizationToken = "#######"
	return []string{"Please fill in the Authorization Token"}
}

func (me *Trello) Schema() map[string]*schema.Schema {
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

		"resolved_list_id": {
			Type:        schema.TypeString,
			Description: "The Trello list to which the card of the resolved problem should be assigned",
			Required:    true,
		},
		"text": {
			Type:        schema.TypeString,
			Description: "The text of the generated Trello card.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"application_key": {
			Type:        schema.TypeString,
			Description: "The application key for the Trello account",
			Required:    true,
		},
		"authorization_token": {
			Type:        schema.TypeString,
			Description: "The application token for the Trello account",
			Optional:    true,
		},
		"board_id": {
			Type:        schema.TypeString,
			Description: "The Trello board to which the card should be assigned",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the Trello card.   You can use same placeholders as in card text",
			Required:    true,
		},
		"list_id": {
			Type:        schema.TypeString,
			Description: "The Trello list to which the card should be assigned",
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

func (me *Trello) MarshalHCL(properties hcl.Properties) error { // The authorization_token field MUST NOT get serialized into HCL here
	// The Dynatrace Settings 2.0 API delivers a scrambled version of any previously stored authorization_token here
	// Evaluation at this point would lead to that scrambled version to make it into the Terraform State
	// As a result any plans would be non-empty
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"active":  me.Enabled,
		"profile": me.ProfileID,

		"resolved_list_id":    me.ResolvedListID,
		"text":                me.Text,
		"application_key":     me.ApplicationKey,
		"board_id":            me.BoardID,
		"description":         me.Description,
		"list_id":             me.ListID,
		"authorization_token": me.AuthorizationToken,
	})
}

func (me *Trello) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"resolved_list_id":    &me.ResolvedListID,
		"text":                &me.Text,
		"application_key":     &me.ApplicationKey,
		"board_id":            &me.BoardID,
		"description":         &me.Description,
		"list_id":             &me.ListID,
		"authorization_token": &me.AuthorizationToken,
	})
}
