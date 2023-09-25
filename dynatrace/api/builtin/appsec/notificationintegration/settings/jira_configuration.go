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

package notificationintegration

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JiraConfiguration struct {
	ApiToken   string `json:"apiToken"`   // The API token for the Jira profile. Using password authentication [was deprecated by Jira](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/)
	IssueType  string `json:"issueType"`  // The type of the Jira issue to be created by this notification.\n\nTo find all available issue types or create your own, in Jira, go to Project settings > Issue types.
	ProjectKey string `json:"projectKey"` // The project key of the Jira issue to be created by this notification.
	Url        string `json:"url"`        // The URL of the Jira API endpoint.
	Username   string `json:"username"`   // The username of the Jira profile.
}

func (me *JiraConfiguration) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_token": {
			Type:        schema.TypeString,
			Description: "The API token for the Jira profile. Using password authentication [was deprecated by Jira](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/)",
			Required:    true,
			Sensitive:   true,
		},
		"issue_type": {
			Type:        schema.TypeString,
			Description: "The type of the Jira issue to be created by this notification.\n\nTo find all available issue types or create your own, in Jira, go to Project settings > Issue types.",
			Required:    true,
		},
		"project_key": {
			Type:        schema.TypeString,
			Description: "The project key of the Jira issue to be created by this notification.",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the Jira API endpoint.",
			Required:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username of the Jira profile.",
			Required:    true,
		},
	}
}

func (me *JiraConfiguration) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"api_token":   me.ApiToken,
		"issue_type":  me.IssueType,
		"project_key": me.ProjectKey,
		"url":         me.Url,
		"username":    me.Username,
	})
}

func (me *JiraConfiguration) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"api_token":   &me.ApiToken,
		"issue_type":  &me.IssueType,
		"project_key": &me.ProjectKey,
		"url":         &me.Url,
		"username":    &me.Username,
	})
}
