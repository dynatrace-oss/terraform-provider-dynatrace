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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
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
			Description: "The name of the notification configuration.",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},

		"application_key": {
			Type:        schema.TypeString,
			Description: "The application key for the Trello account.\n\nYou must be logged into Trello to have Trello automatically generate an application key for you. [Get application key](https://trello.com/app-key)",
			Required:    true,
		},
		"authorization_token": {
			Type:        schema.TypeString,
			Description: "The authorization token for the Trello account.",
			Required:    true,
			Sensitive:   true,
		},
		"board_id": {
			Type:        schema.TypeString,
			Description: "Trello board ID problem cards should be assigned to",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the Trello card. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemDetailsMarkdown}**: All problem event details including root cause as a Markdown-formatted string.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist.",
			Required:    true,
		},
		"list_id": {
			Type:        schema.TypeString,
			Description: "Trello list ID new problem cards should be assigned to",
			Required:    true,
		},
		"resolved_list_id": {
			Type:        schema.TypeString,
			Description: "Trello list ID resolved problem cards should be assigned to",
			Required:    true,
		},
		"text": {
			Type:        schema.TypeString,
			Description: "The card text and problem placeholders to appear on new problem cards. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist.",
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
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMapPlus(
		me.Schema(),
		map[string]any{
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
		},
		[]string{"authorization_token"},
	))
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
