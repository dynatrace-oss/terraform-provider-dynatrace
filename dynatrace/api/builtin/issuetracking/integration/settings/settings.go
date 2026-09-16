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

package integration

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled            bool               `json:"enabled"`            // This setting is enabled (`true`) or disabled (`false`)
	InsertAfter        *string            `json:"-"`                  // Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched
	Issuelabel         string             `json:"issuelabel"`         // Set a label to identify these issues, for example, `release_blocker` or `non-critical`
	Issuequery         string             `json:"issuequery"`         // You can use the following placeholders to automatically insert values from the **Release monitoring** page in your query: `{NAME}`, `{VERSION}`, `{STAGE}`, `{PRODUCT}`.
	Issuetheme         IssueTheme         `json:"issuetheme"`         // Select the issue type to be displayed. Possible values: `ERROR`, `INFO`, `RESOLVED`
	Issuetrackersystem IssueTrackerSystem `json:"issuetrackersystem"` // Select the issue-tracking system you want to query. Possible values: `GITHUB`, `GITLAB`, `JIRA`, `JIRA_CLOUD`, `JIRA_ON_PREMISE`, `SERVICENOW`
	Password           *string            `json:"password,omitempty"` // Password
	Token              *string            `json:"token,omitempty"`    // Token
	Url                string             `json:"url"`                // For Jira, use the base URL (for example, https://jira.yourcompany.com); for GitHub, use the repository URL (for example, https://github.com/org/repo); for GitLab, use the specific project API for a single project (for example, https://gitlab.com/api/v4/projects/:projectId), and the specific group API for a multiple projects (for example, https://gitlab.com/api/v4/groups/:groupId); for ServiceNow, use your company instance URL (for example, https://yourinstance.service-now.com/)
	Username           string             `json:"username"`           // Username
}

func (me *Settings) Name() string {
	return me.Issuelabel
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Computed:    true,
			Optional:    true,
		},
		"issuelabel": {
			Type:        schema.TypeString,
			Description: "Set a label to identify these issues, for example, `release_blocker` or `non-critical`",
			Required:    true,
		},
		"issuequery": {
			Type:             schema.TypeString,
			Description:      "You can use the following placeholders to automatically insert values from the **Release monitoring** page in your query: `{NAME}`, `{VERSION}`, `{STAGE}`, `{PRODUCT}`.",
			Required:         true,
			DiffSuppressFunc: hcl.SuppressEOT,
		},
		"issuetheme": {
			Type:        schema.TypeString,
			Description: "Select the issue type to be displayed. Possible values: `ERROR`, `INFO`, `RESOLVED`",
			Required:    true,
		},
		"issuetrackersystem": {
			Type:        schema.TypeString,
			Description: "Select the issue-tracking system you want to query. Possible values: `GITHUB`, `GITLAB`, `JIRA`, `JIRA_CLOUD`, `JIRA_ON_PREMISE`, `SERVICENOW`",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Password",
			Optional:    true, // nullable & precondition
			Sensitive:   true,
		},
		"token": {
			Type:        schema.TypeString,
			Description: "Token",
			Optional:    true, // nullable & precondition
			Sensitive:   true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "For Jira, use the base URL (for example, https://jira.yourcompany.com); for GitHub, use the repository URL (for example, https://github.com/org/repo); for GitLab, use the specific project API for a single project (for example, https://gitlab.com/api/v4/projects/:projectId), and the specific group API for a multiple projects (for example, https://gitlab.com/api/v4/groups/:groupId); for ServiceNow, use your company instance URL (for example, https://yourinstance.service-now.com/)",
			Required:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "Username",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(), map[string]any{
			"enabled":            me.Enabled,
			"insert_after":       me.InsertAfter,
			"issuelabel":         me.Issuelabel,
			"issuequery":         me.Issuequery,
			"issuetheme":         me.Issuetheme,
			"issuetrackersystem": me.Issuetrackersystem,
			"password":           sensitive.SecretValue,
			"token":              sensitive.SecretValue,
			"url":                me.Url,
			"username":           me.Username,
		}))
}

func (me *Settings) HandlePreconditions() error {
	if (me.Password != nil) && (!slices.Contains([]string{"JIRA", "JIRA_ON_PREMISE", "SERVICENOW"}, string(me.Issuetrackersystem))) {
		return fmt.Errorf("'password' must not be specified unless 'issuetrackersystem' is one of ['JIRA', 'JIRA_ON_PREMISE', 'SERVICENOW']; got 'issuetrackersystem'='%v'", me.Issuetrackersystem)
	}
	if (me.Token != nil) && (!slices.Contains([]string{"JIRA", "GITHUB", "GITLAB", "JIRA_CLOUD"}, string(me.Issuetrackersystem))) {
		return fmt.Errorf("'token' must not be specified unless 'issuetrackersystem' is one of ['JIRA', 'GITHUB', 'GITLAB', 'JIRA_CLOUD']; got 'issuetrackersystem'='%v'", me.Issuetrackersystem)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":            &me.Enabled,
		"insert_after":       &me.InsertAfter,
		"issuelabel":         &me.Issuelabel,
		"issuequery":         &me.Issuequery,
		"issuetheme":         &me.Issuetheme,
		"issuetrackersystem": &me.Issuetrackersystem,
		"password":           &me.Password,
		"token":              &me.Token,
		"url":                &me.Url,
		"username":           &me.Username,
	})
}

const credsNotProvided = "REST API didn't provide token data"

func (me *Settings) FillDemoValues() []string {
	if slices.Contains([]string{"JIRA", "GITHUB", "GITLAB", "JIRA_CLOUD"}, string(me.Issuetrackersystem)) {
		me.Token = new("################")
		return []string{credsNotProvided}
	}
	if slices.Contains([]string{"JIRA", "JIRA_ON_PREMISE", "SERVICENOW"}, string(me.Issuetrackersystem)) {
		me.Password = new("################")
		return []string{credsNotProvided}
	}
	return nil
}
