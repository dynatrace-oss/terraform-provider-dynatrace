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

// SlackConfig Configuration of the Slack notification.
type SlackConfig struct {
	BaseNotificationConfig
	Title   string  `json:"title"`         // The content of the message.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas.
	URL     *string `json:"url,omitempty"` // The URL of the Slack WebHook.  This is confidential information, therefore GET requests return this field with the `null` value, and it is optional for PUT requests.
	Channel string  `json:"channel"`       // The channel (for example, `#general`) or the user (for example, `@john.smith`) to send the message to.
}

func (me *SlackConfig) GetType() Type {
	return Types.Slack
}

func (me *SlackConfig) Schema() map[string]*schema.Schema {
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
		"title": {
			Type:        schema.TypeString,
			Description: "The content of the message.  You can use the following placeholders:  * `{ImpactedEntity}`: The entity impacted by the problem or *X* impacted entities.  * `{PID}`: The ID of the reported problem.  * `{ProblemDetailsText}`: All problem event details, including root cause, as a text-formatted string.  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.  * `{ProblemURL}`: The URL of the problem within Dynatrace.  * `{State}`: The state of the problem. Possible values are `OPEN` and `RESOLVED`.  * `{Tags}`: The list of tags that are defined for all impacted entities, separated by commas",
			Required:    true,
		},
		"url": {
			Type:        schema.TypeString,
			Description: "The URL of the Slack WebHook.  This is confidential information, therefore GET requests return this field with the `null` value, and it is optional for PUT requests",
			Optional:    true,
		},
		"channel": {
			Type:        schema.TypeString,
			Description: "The channel (for example, `#general`) or the user (for example, `@john.smith`) to send the message to",
			Required:    true,
		},
	}
}

func (me *SlackConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("slack.0.url"); ok && len(value.(string)) > 0 {
		me.URL = opt.NewString(value.(string))
	}
	return nil
}

func (me *SlackConfig) FillDemoValues() []string {
	me.URL = opt.NewString("https://www.google.at/f75e68ef-aca7-3a07-9c21-94eb00ecfc56")
	return []string{"The REST API didn't provide the credentials"}
}

func (me *SlackConfig) MarshalHCL(properties hcl.Properties) error {
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
	if err := properties.Encode("title", me.Title); err != nil {
		return err
	}
	if err := properties.Encode("url", me.URL); err != nil {
		return err
	}
	if err := properties.Encode("url", me.URL); err != nil {
		return err
	}
	if err := properties.Encode("channel", me.Channel); err != nil {
		return err
	}
	return nil
}

func (me *SlackConfig) UnmarshalHCL(decoder hcl.Decoder) error {
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
		delete(me.Unknowns, "title")
		delete(me.Unknowns, "url")
		delete(me.Unknowns, "channel")
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
	if value, ok := decoder.GetOk("title"); ok {
		me.Title = value.(string)
	}
	if value, ok := decoder.GetOk("url"); ok {
		me.URL = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("channel"); ok {
		me.Channel = value.(string)
	}
	return nil
}

func (me *SlackConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":              me.ID,
		"name":            me.Name,
		"type":            me.GetType(),
		"active":          me.Active,
		"alertingProfile": me.AlertingProfile,
		"title":           me.Title,
		"url":             me.URL,
		"channel":         me.Channel,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *SlackConfig) UnmarshalJSON(data []byte) error {
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
		"title":           &me.Title,
		"url":             &me.URL,
		"channel":         &me.Channel,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
