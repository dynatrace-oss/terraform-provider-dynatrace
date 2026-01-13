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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/http"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type XMatters struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	URL      string       `json:"url"`                  // The URL of the xMatters webhook
	Insecure bool         `json:"acceptAnyCertificate"` // Accept any SSL certificate (including self-signed and invalid certificates)
	Headers  http.Headers `json:"headers,omitempty"`    // Additional HTTP headers
	Payload  string       `json:"payload"`              // The content of the notification message. Type '{' for placeholder suggestions
}

func (me *XMatters) FillDemoValues() []string {
	if len(me.Headers) == 0 {
		return []string{}
	}
	filled := false
	for _, header := range me.Headers {
		if len(header.FillDemoValues()) > 0 {
			filled = true
		}
	}
	if filled {
		return []string{"One or more secret HTTP headers need to get filled in"}
	}

	return []string{}
}

func (me *XMatters) Schema() map[string]*schema.Schema {
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

		"insecure": {
			Type:        schema.TypeBool,
			Description: "Accept any SSL certificate (including self-signed and invalid certificates)",
			Optional:    true,
		},
		"headers": {
			Type:        schema.TypeList,
			Description: "A list of the additional HTTP headers.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(http.Headers).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"payload": {
			Type:             schema.TypeString,
			Description:      "The content of the notification message. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntities}**: Details about the entities impacted by the problem in form of a json array.\n\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemDetailsHTML}**: All problem event details including root cause as an HTML-formatted string.\n\n**{ProblemDetailsJSONv2}**: Problem as json object following the structure from the [Dynatrace Problems V2 API](https://dt-url.net/7a03ti2). The optional fields evidenceDetails and impactAnalysis are included, but recentComments is not.\n\n**{ProblemDetailsJSON}**: Problem as json object following the structure from the [Dynatrace Problems V1 API](https://dt-url.net/qn23tk2).\n\n**{ProblemDetailsMarkdown}**: All problem event details including root cause as a Markdown-formatted string.\n\n**{ProblemDetailsText}**: All problem event details including root cause as a text-formatted string.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist.",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the xMatters webhook.",
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

func (me *XMatters) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMapPlus(
		me.Schema(),
		map[string]any{
			"name":    me.Name,
			"active":  me.Enabled,
			"profile": me.ProfileID,

			"url":      me.URL,
			"insecure": me.Insecure,
			"headers":  me.Headers,
			"payload":  me.Payload,
		},
		me.Headers.GenIgnoreChanges("headers"),
	))
}

func (me *XMatters) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"url":      &me.URL,
		"insecure": &me.Insecure,
		"headers":  &me.Headers,
		"payload":  &me.Payload,
	})
}
