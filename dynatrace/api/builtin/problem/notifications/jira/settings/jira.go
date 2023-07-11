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

type Jira struct {
	Enabled   bool   `json:"-"`
	Name      string `json:"-"`
	ProfileID string `json:"-"`

	URL         string `json:"url"`         // Jira endpoint URL
	Username    string `json:"username"`    // The username of the Jira profile
	APIToken    string `json:"apiToken"`    // The API token for the Jira profile. Using password authentication [was deprecated by Jira](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/)
	ProjectKey  string `json:"projectKey"`  // The project key of the Jira issue to be created by this notification
	IssueType   string `json:"issueType"`   // The type of the Jira issue to be created by this notification.\n\nTo find all available issue types, or to create your own issue type, within JIRA go to Options > Issues
	Summary     string `json:"summary"`     // The summary of the Jira issue to be created by this notification. Type '{' for placeholder suggestions
	Description string `json:"description"` // The description of the Jira issue to be created by this notification. Type '{' for placeholder suggestions
}

func (me *Jira) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if apiToken, ok := decoder.GetOk("api_token"); ok && len(apiToken.(string)) > 0 {
		me.APIToken = apiToken.(string)
	}
	return nil
}

func (me *Jira) FillDemoValues() []string {
	me.APIToken = "#######"
	return []string{"Please fill in the API Token"}
}

func (me *Jira) Schema() map[string]*schema.Schema {
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
			Description: "The URL of the Jira API endpoint",
			Required:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username of the Jira profile",
			Required:    true,
		},
		"api_token": {
			Type:        schema.TypeString,
			Sensitive:   true,
			Description: "The API token for the Jira profile. Using password authentication [was deprecated by Jira](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/)",
			Optional:    true,
		},
		"project_key": {
			Type:        schema.TypeString,
			Description: "The project key of the Jira issue to be created by this notification",
			Required:    true,
		},
		"issue_type": {
			Type:        schema.TypeString,
			Description: "The type of the Jira issue to be created by this notification",
			Required:    true,
		},
		"summary": {
			Type:        schema.TypeString,
			Description: "The summary of the Jira issue to be created by this notification.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the Jira issue to be created by this notification.   You can use same placeholders as in issue summary",
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

func (me *Jira) MarshalHCL(properties hcl.Properties) error { // The api_token field MUST NOT get serialized into HCL here
	// The Dynatrace Settings 2.0 API delivers a scrambled version of any previously stored api_token here
	// Evaluation at this point would lead to that scrambled version to make it into the Terraform State
	// As a result any plans would be non-empty
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"active":  me.Enabled,
		"profile": me.ProfileID,

		"url":         me.URL,
		"username":    me.Username,
		"project_key": me.ProjectKey,
		"issue_type":  me.IssueType,
		"summary":     me.Summary,
		"description": me.Description,
		"api_token":   me.APIToken,
	})
}

func (me *Jira) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"active":  &me.Enabled,
		"profile": &me.ProfileID,

		"url":         &me.URL,
		"username":    &me.Username,
		"api_token":   &me.APIToken,
		"project_key": &me.ProjectKey,
		"issue_type":  &me.IssueType,
		"summary":     &me.Summary,
		"description": &me.Description,
	})
}
