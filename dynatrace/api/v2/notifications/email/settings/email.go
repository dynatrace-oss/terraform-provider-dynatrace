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

type Email struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	Subject              string   `json:"subject"`              // #### Available placeholders\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem (or multiple impacted entities).\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist
	Recipients           []string `json:"recipients"`           // pattern: ^[\\.a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+?@[a-zA-Z0-9](?:[a-zA-Z0-9-]*?[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*?[a-zA-Z0-9])?)*$
	CCRecipients         []string `json:"ccRecipients"`         // pattern: ^[\\.a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+?@[a-zA-Z0-9](?:[a-zA-Z0-9-]*?[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*?[a-zA-Z0-9])?)*$
	BCCRecipients        []string `json:"bccRecipients"`        // pattern: ^[\\.a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+?@[a-zA-Z0-9](?:[a-zA-Z0-9-]*?[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*?[a-zA-Z0-9])?)*$
	NotifyClosedProblems bool     `json:"notifyClosedProblems"` // Send email if problem is closed
	Body                 string   `json:"body"`                 // The template of the email notifications. Type '{' for placeholder suggestions
}

func (me *Email) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration",
			Required:    true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "The configuration is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},
		"bcc": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The list of the email BCC-recipients",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"cc": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The list of the email CC-recipients",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"to": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "The list of the email recipients",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"body": {
			Type:        schema.TypeString,
			Description: "The template of the email notification.  You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"subject": {
			Type:        schema.TypeString,
			Description: "The subject of the email notifications",
			Required:    true,
		},
		"notify_closed_problems": {
			Type:        schema.TypeBool,
			Description: "Send email if problem is closed",
			Optional:    true,
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of these settings when referred to from resources requiring the REST API V1 keys",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *Email) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"active":  me.Enabled,
		"profile": me.ProfileID,

		"bcc":                    me.BCCRecipients,
		"cc":                     me.CCRecipients,
		"to":                     me.Recipients,
		"body":                   me.Body,
		"subject":                me.Subject,
		"notify_closed_problems": me.NotifyClosedProblems,
	})
}

func (me *Email) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"bcc":                    &me.BCCRecipients,
		"cc":                     &me.CCRecipients,
		"to":                     &me.Recipients,
		"body":                   &me.Body,
		"subject":                &me.Subject,
		"notify_closed_problems": &me.NotifyClosedProblems,
	})
}
