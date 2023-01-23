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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/http"

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

		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the WebHook endpoint",
			Required:    true,
		},
		"insecure": {
			Type:        schema.TypeBool,
			Description: "Accept any, including self-signed and invalid, SSL certificate (`true`) or only trusted (`false`) certificates",
			Optional:    true,
		},
		"headers": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Description: "A list of the additional HTTP headers",
			Elem:        &schema.Resource{Schema: new(http.Headers).Schema()},
		},
		"payload": {
			Type:             schema.TypeString,
			Description:      "The content of the notification message. You can use the following placeholders:  * `{ImpactedEntities}`: Details about the entities impacted by the problem in form of a JSON array.  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsHTML}`: All problem event details, including root cause, as an HTML-formatted string.  * `{ProblemDetailsJSON}`: All problem event details, including root cause, as a JSON object.  * `{ProblemDetailsMarkdown}`: All problem event details, including root cause, as a [Markdown-formatted](https://www.markdownguide.org/cheat-sheet/) string.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
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
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"active":  me.Enabled,
		"profile": me.ProfileID,

		"url":      me.URL,
		"insecure": me.Insecure,
		"headers":  me.Headers,
		"payload":  me.Payload,
	})
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
