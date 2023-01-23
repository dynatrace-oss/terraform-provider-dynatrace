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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NotificationConfig Configuration of a notification. The actual set of fields depends on the `type` of the notification.
// See the [Notifications API - JSON models](https://www.dynatrace.com/support/help/shortlink/api-config-notifications-models) help topic for example models of every notification type.
type NotificationConfig interface {
	GetType() Type
	GetID() *string
	SetID(*string)
	GetName() string
	MarshalHCL(hcl.Properties) error
	UnmarshalHCL(decoder hcl.Decoder) error
}

// BaseNotificationConfig Configuration of a notification. The actual set of fields depends on the `type` of the notification.
// See the [Notifications API - JSON models](https://www.dynatrace.com/support/help/shortlink/api-config-notifications-models) help topic for example models of every notification type.
type BaseNotificationConfig struct {
	ID              *string                    `json:"id,omitempty" hcl:"-"`                   // The ID of the notification configuration.
	Name            string                     `json:"name" hcl:"name"`                        // The name of the notification configuration.
	Type            Type                       `json:"type" hcl:"type"`                        // Defines the actual set of fields depending on the value. See one of the following objects:  * `EMAIL` -> EmailNotificationConfig  * `PAGER_DUTY` -> PagerDutyNotificationConfig  * `WEBHOOK` -> WebHookNotificationConfig  * `SLACK` -> SlackNotificationConfig  * `HIPCHAT` -> HipChatNotificationConfig  * `VICTOROPS` -> VictorOpsNotificationConfig  * `SERVICE_NOW` -> ServiceNowNotificationConfig  * `XMATTERS` -> XMattersNotificationConfig  * `ANSIBLETOWER` -> AnsibleTowerNotificationConfig  * `OPS_GENIE` -> OpsGenieNotificationConfig  * `JIRA` -> JiraNotificationConfig  * `TRELLO` -> TrelloNotificationConfig
	Active          bool                       `json:"active" hcl:"active"`                    // The configuration is enabled (`true`) or disabled (`false`).
	AlertingProfile string                     `json:"alertingProfile" hcl:"alerting_profile"` // The ID of the associated alerting profile.
	Unknowns        map[string]json.RawMessage `json:"-"`
}

func (me *BaseNotificationConfig) GetType() Type {
	return me.Type
}

func (me *BaseNotificationConfig) GetID() *string {
	return me.ID
}

func (me *BaseNotificationConfig) GetName() string {
	return me.Name
}

func (me *BaseNotificationConfig) SetID(id *string) {
	me.ID = id
}

func (me *BaseNotificationConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the notification configuration",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Defines the actual set of fields depending on the value. See one of the following objects:  * `EMAIL` -> EmailNotificationConfig  * `PAGER_DUTY` -> PagerDutyNotificationConfig  * `WEBHOOK` -> WebHookNotificationConfig  * `SLACK` -> SlackNotificationConfig  * `HIPCHAT` -> HipChatNotificationConfig  * `VICTOROPS` -> VictorOpsNotificationConfig  * `SERVICE_NOW` -> ServiceNowNotificationConfig  * `XMATTERS` -> XMattersNotificationConfig  * `ANSIBLETOWER` -> AnsibleTowerNotificationConfig  * `OPS_GENIE` -> OpsGenieNotificationConfig  * `JIRA` -> JiraNotificationConfig  * `TRELLO` -> TrelloNotificationConfig",
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
	}
}

func (me *BaseNotificationConfig) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("active", me.Active); err != nil {
		return err
	}
	if err := properties.Encode("type", me.Type); err != nil {
		return err
	}
	if err := properties.Encode("alerting_profile", me.AlertingProfile); err != nil {
		return err
	}

	return nil
}

func (me *BaseNotificationConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.Type = Type(value.(string))
	}

	if value, ok := decoder.GetOk("active"); ok {
		me.Active = value.(bool)
	}
	if value, ok := decoder.GetOk("alerting_profile"); ok {
		me.AlertingProfile = value.(string)
	}

	return nil
}

func (me *BaseNotificationConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":              me.ID,
		"name":            me.Name,
		"type":            me.Type,
		"active":          me.Active,
		"alertingProfile": me.AlertingProfile,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *BaseNotificationConfig) UnmarshalJSON(data []byte) error {
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
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}

// Type Defines the actual set of fields depending on the value. See one of the following objects:
// * `EMAIL` -> EmailNotificationConfig
// * `PAGER_DUTY` -> PagerDutyNotificationConfig
// * `WEBHOOK` -> WebHookNotificationConfig
// * `SLACK` -> SlackNotificationConfig
// * `HIPCHAT` -> HipChatNotificationConfig
// * `VICTOROPS` -> VictorOpsNotificationConfig
// * `SERVICE_NOW` -> ServiceNowNotificationConfig
// * `XMATTERS` -> XMattersNotificationConfig
// * `ANSIBLETOWER` -> AnsibleTowerNotificationConfig
// * `OPS_GENIE` -> OpsGenieNotificationConfig
// * `JIRA` -> JiraNotificationConfig
// * `TRELLO` -> TrelloNotificationConfig
type Type string

// Types offers the known enum values
var Types = struct {
	Ansibletower Type
	Email        Type
	Hipchat      Type
	Jira         Type
	OpsGenie     Type
	PagerDuty    Type
	ServiceNow   Type
	Slack        Type
	Trello       Type
	Victorops    Type
	Webhook      Type
	Xmatters     Type
}{
	"ANSIBLETOWER",
	"EMAIL",
	"HIPCHAT",
	"JIRA",
	"OPS_GENIE",
	"PAGER_DUTY",
	"SERVICE_NOW",
	"SLACK",
	"TRELLO",
	"VICTOROPS",
	"WEBHOOK",
	"XMATTERS",
}
