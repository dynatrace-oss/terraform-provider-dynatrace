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
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/http"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WebHook struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	URL                      string             `json:"url"`                         // The URL of the webhook endpoint
	Insecure                 bool               `json:"acceptAnyCertificate"`        // Accept any SSL certificate (including self-signed and invalid certificates)
	NotifyEventMergesEnabled bool               `json:"notifyEventMergesEnabled"`    // Call webhook if new events merge into existing problems
	NotifyClosedProblems     bool               `json:"notifyClosedProblems"`        // Call webhook if problem is closed
	Headers                  http.Headers       `json:"headers,omitempty"`           // Additional HTTP headers
	Payload                  string             `json:"payload"`                     // The content of the notification message. Type '{' for placeholder suggestions
	UseOAuth2                *bool              `json:"useOAuth2,omitempty"`         // Use OAuth 2.0 for authentication
	OAuth2Credentials        *OAuth2Credentials `json:"oAuth2Credentials,omitempty"` // To authenticate your integration, the OAuth 2.0 *Client Credentials* Flow (Grant Type) is used. For details see [Client Credentials Flow](https://dt-url.net/ym22wsm)).\n\nThe obtained Access Token is subsequently provided in the *Authorization* header of the request carrying the notification payload.
	SecretUrl                *string            `json:"secretUrl,omitempty"`         // The secret URL of the webhook endpoint.
	UrlContainsSecret        *bool              `json:"urlContainsSecret,omitempty"` // Secret webhook URL
}

func (me *WebHook) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if url, ok := decoder.GetOk("secret_url"); ok && len(url.(string)) > 0 {
		me.SecretUrl = opt.NewString(url.(string))
	}
	return nil
}

func (me *WebHook) FillDemoValues() []string {
	result := []string{}
	if len(me.Headers) > 0 {
		filled := false
		for _, header := range me.Headers {
			if len(header.FillDemoValues()) > 0 {
				filled = true
			}
		}
		if filled {
			result = append(result, "One or more secret HTTP headers need to get filled in")
		}
	}
	if me.UseOAuth2 != nil && *me.UseOAuth2 {
		// me.OAuth2Credentials.ClientSecret = "#######"
		result = append(result, me.OAuth2Credentials.FillDemoValues()...)
	}

	return result
}

func (me *WebHook) Schema() map[string]*schema.Schema {
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
			ForceNew:    http.ForceNewOnHeaders,
		},
		"notify_closed_problems": {
			Type:        schema.TypeBool,
			Description: "Call webhook if problem is closed",
			Optional:    true,
		},
		"notify_event_merges": {
			Type:        schema.TypeBool,
			Description: "Call webhook if new events merge into existing problems",
			Optional:    true,
		},
		"oauth_2_credentials": {
			Type:        schema.TypeList,
			Description: "To authenticate your integration, the OAuth 2.0 *Client Credentials* Flow (Grant Type) is used. For details see [Client Credentials Flow](https://dt-url.net/ym22wsm)).\n\nThe obtained Access Token is subsequently provided in the *Authorization* header of the request carrying the notification payload.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(OAuth2Credentials).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"payload": {
			Type:             schema.TypeString,
			Description:      "The content of the notification message. Type '{' for placeholder suggestions.. #### Available placeholders\n**{ImpactedEntities}**: Details about the entities impacted by the problem in form of a json array.\n\n**{ImpactedEntity}**: A short description of the problem and impacted entity (or multiple impacted entities).\n\n**{ImpactedEntityNames}**: The entity impacted by the problem.\n\n**{NamesOfImpactedEntities}**: The names of all entities that are impacted by the problem.\n\n**{PID}**: Unique system identifier of the reported problem.\n\n**{ProblemDetailsHTML}**: All problem event details including root cause as an HTML-formatted string.\n\n**{ProblemDetailsJSONv2}**: Problem as json object following the structure from the [Dynatrace Problems V2 API](https://dt-url.net/7a03ti2). The optional fields evidenceDetails and impactAnalysis are included, but recentComments is not.\n\n**{ProblemDetailsJSON}**: Problem as json object following the structure from the [Dynatrace Problems V1 API](https://dt-url.net/qn23tk2).\n\n**{ProblemDetailsMarkdown}**: All problem event details including root cause as a Markdown-formatted string.\n\n**{ProblemDetailsText}**: All problem event details including root cause as a text-formatted string.\n\n**{ProblemID}**: Display number of the reported problem.\n\n**{ProblemImpact}**: Impact level of the problem. Possible values are APPLICATION, SERVICE, or INFRASTRUCTURE.\n\n**{ProblemSeverity}**: Severity level of the problem. Possible values are AVAILABILITY, ERROR, PERFORMANCE, RESOURCE_CONTENTION, or CUSTOM_ALERT.\n\n**{ProblemTitle}**: Short description of the problem.\n\n**{ProblemURL}**: URL of the problem within Dynatrace.\n\n**{State}**: Problem state. Possible values are OPEN or RESOLVED.\n\n**{Tags}**: Comma separated list of tags that are defined for all impacted entities. To refer to the value of a specific tag, specify the tag's key in square brackets: **{Tags[key]}**. If the tag does not have any assigned value, the placeholder will be replaced by an empty string. The placeholder will not be replaced if the tag key does not exist.",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressJSONorEOT,
		},
		"secret_url": {
			Type:        schema.TypeString,
			Description: "The secret URL of the webhook endpoint.",
			Optional:    true, // precondition
			Sensitive:   true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the webhook endpoint.",
			Optional:    true, // precondition
		},
		"url_contains_secret": {
			Type:        schema.TypeBool,
			Description: "Secret webhook URL",
			Optional:    true, // nullable
		},
		"use_oauth_2": {
			Type:        schema.TypeBool,
			Description: "Use OAuth 2.0 for authentication",
			Optional:    true, // nullable
		},
		"legacy_id": {
			Type:        schema.TypeString,
			Description: "The ID of these settings when referred to from resources requiring the REST API V1 keys",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *WebHook) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMapPlus(
		me.Schema(),
		map[string]any{
			"name":    me.Name,
			"active":  me.Enabled,
			"profile": me.ProfileID,

			"notify_event_merges":    me.NotifyEventMergesEnabled,
			"insecure":               me.Insecure,
			"headers":                me.Headers,
			"payload":                me.Payload,
			"url":                    me.URL,
			"notify_closed_problems": me.NotifyClosedProblems,
			"use_oauth_2":            me.UseOAuth2,
			"oauth_2_credentials":    me.OAuth2Credentials,
			"secret_url":             me.SecretUrl,
			"url_contains_secret":    me.UrlContainsSecret,
		},
		append(
			me.Headers.GenIgnoreChanges("headers"),
			"oauth_2_credentials",
		),
	))
}

func (me *WebHook) HandlePreconditions() error {
	if (me.OAuth2Credentials == nil) && (me.UseOAuth2 != nil && *me.UseOAuth2) {
		return fmt.Errorf("'oauth_2_credentials' must be specified if 'use_oauth_2' is set to '%v'", me.UseOAuth2)
	}
	if (me.OAuth2Credentials != nil) && (me.UseOAuth2 != nil && !*me.UseOAuth2) {
		return fmt.Errorf("'oauth_2_credentials' must not be specified if 'use_oauth_2' is set to '%v'", me.UseOAuth2)
	}
	if (me.SecretUrl == nil) && (me.UrlContainsSecret != nil && *me.UrlContainsSecret) {
		return fmt.Errorf("'secret_url' must be specified if 'url_contains_secret' is set to '%v'", me.UrlContainsSecret)
	}
	return nil
}

func (me *WebHook) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"notify_event_merges":    &me.NotifyEventMergesEnabled,
		"insecure":               &me.Insecure,
		"headers":                &me.Headers,
		"payload":                &me.Payload,
		"url":                    &me.URL,
		"notify_closed_problems": &me.NotifyClosedProblems,
		"use_oauth_2":            &me.UseOAuth2,
		"oauth_2_credentials":    &me.OAuth2Credentials,
		"secret_url":             &me.SecretUrl,
		"url_contains_secret":    &me.UrlContainsSecret,
	})
}
