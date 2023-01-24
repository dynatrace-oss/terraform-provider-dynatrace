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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// JiraConfig Configuration of the Jira notification.
type JiraConfig struct {
	BaseNotificationConfig
	IssueType   string  `json:"issueType"`          // The type of the Jira issue to be created by this notification.
	Password    *string `json:"password,omitempty"` // The password for the Jira profile.
	ProjectKey  string  `json:"projectKey"`         // The project key of the Jira issue to be created by this notification.
	Summary     string  `json:"summary"`            // The summary of the Jira issue to be created by this notification.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas.
	URL         string  `json:"url"`                // The URL of the Jira API endpoint.
	Username    string  `json:"username"`           // The username of the Jira profile.
	Description string  `json:"description"`        // The description of the Jira issue to be created by this notification.   You can use same placeholders as in issue summary.
}

func (me *JiraConfig) GetType() Type {
	return Types.Jira
}

func (me *JiraConfig) Schema() map[string]*schema.Schema {
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
		"alerting_profile": {
			Type:        schema.TypeString,
			Description: "The ID of the associated alerting profile",
			Required:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
		"issue_type": {
			Type:        schema.TypeString,
			Description: "The type of the Jira issue to be created by this notification",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password for the Jira profile",
			Optional:    true,
		},
		"project_key": {
			Type:        schema.TypeString,
			Description: "The project key of the Jira issue to be created by this notification",
			Required:    true,
		},
		"summary": {
			Type:        schema.TypeString,
			Description: "The summary of the Jira issue to be created by this notification.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
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
		"description": {
			Type:        schema.TypeString,
			Description: "The description of the Jira issue to be created by this notification.   You can use same placeholders as in issue summary",
			Required:    true,
		},
	}
}

func (me *JiraConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("jira.0.password"); ok && len(value.(string)) > 0 {
		me.Password = opt.NewString(value.(string))
	}
	return nil
}

func (me *JiraConfig) FillDemoValues() []string {
	me.Password = opt.NewString("#######")
	return []string{"The REST API didn't provide the credentials"}
}

func (me *JiraConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("active", me.Active); err != nil {
		return err
	}
	if err := properties.Encode("alerting_profile", me.AlertingProfile); err != nil {
		return err
	}
	if err := properties.Encode("issue_type", me.IssueType); err != nil {
		return err
	}
	if err := properties.Encode("password", me.Password); err != nil {
		return err
	}
	if err := properties.Encode("project_key", me.ProjectKey); err != nil {
		return err
	}
	if err := properties.Encode("summary", me.Summary); err != nil {
		return err
	}
	if err := properties.Encode("url", me.URL); err != nil {
		return err
	}
	if err := properties.Encode("username", me.Username); err != nil {
		return err
	}
	if err := properties.Encode("description", me.Description); err != nil {
		return err
	}
	return nil
}

func (me *JiraConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "name")
		delete(me.Unknowns, "active")
		delete(me.Unknowns, "alerting_profile")
		delete(me.Unknowns, "issue_type")
		delete(me.Unknowns, "password")
		delete(me.Unknowns, "project_key")
		delete(me.Unknowns, "summary")
		delete(me.Unknowns, "url")
		delete(me.Unknowns, "username")
		delete(me.Unknowns, "description")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("active"); ok {
		me.Active = value.(bool)
	}
	if value, ok := decoder.GetOk("alerting_profile"); ok {
		me.AlertingProfile = value.(string)
	}
	if value, ok := decoder.GetOk("issue_type"); ok {
		me.IssueType = value.(string)
	}
	if value, ok := decoder.GetOk("password"); ok {
		me.Password = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("project_key"); ok {
		me.ProjectKey = value.(string)
	}
	if value, ok := decoder.GetOk("summary"); ok {
		me.Summary = value.(string)
	}
	if value, ok := decoder.GetOk("url"); ok {
		me.URL = value.(string)
	}
	if value, ok := decoder.GetOk("username"); ok {
		me.Username = value.(string)
	}
	if value, ok := decoder.GetOk("description"); ok {
		me.Description = value.(string)
	}
	return nil
}

func (me *JiraConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":              me.ID,
		"name":            me.Name,
		"type":            me.GetType(),
		"active":          me.Active,
		"alertingProfile": me.AlertingProfile,
		"issueType":       me.IssueType,
		"password":        me.Password,
		"projectKey":      me.ProjectKey,
		"summary":         me.Summary,
		"url":             me.URL,
		"username":        me.Username,
		"description":     me.Description,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *JiraConfig) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"id":              &me.ID,
		"name":            &me.Name,
		"type":            &me.Type,
		"active":          &me.Active,
		"alertingProfile": &me.AlertingProfile,
		"issueType":       &me.IssueType,
		"password":        &me.Password,
		"projectKey":      &me.ProjectKey,
		"summary":         &me.Summary,
		"url":             &me.URL,
		"username":        &me.Username,
		"description":     &me.Description,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
