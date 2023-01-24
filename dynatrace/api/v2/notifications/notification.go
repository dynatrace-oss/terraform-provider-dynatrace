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
	"errors"
	"fmt"

	ansible "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/ansible/settings"
	email "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/email/settings"
	jira "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/jira/settings"
	opsgenie "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/opsgenie/settings"
	pagerduty "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/pagerduty/settings"
	servicenow "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/servicenow/settings"
	slack "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/slack/settings"
	trello "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/trello/settings"
	victorops "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/victorops/settings"
	webhook "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/webhook/settings"
	xmatters "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v2/notifications/xmatters/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Notification struct {
	Enabled   bool    `json:"enabled"`
	Name      string  `json:"displayName"`
	ProfileID string  `json:"alertingProfile"`
	Type      Type    `json:"type"`
	LegacyID  *string `json:"-"`

	Email        *email.Email           `json:"emailNotification,omitempty"`
	Slack        *slack.Slack           `json:"slackNotification,omitempty"`
	Jira         *jira.Jira             `json:"jiraNotification,omitempty"`
	AnsibleTower *ansible.AnsibleTower  `json:"ansibleTowerNotification,omitempty"`
	OpsGenie     *opsgenie.OpsGenie     `json:"opsGenieNotification,omitempty"`
	PagerDuty    *pagerduty.PagerDuty   `json:"pagerDutyNotification,omitempty"`
	VictorOps    *victorops.VictorOps   `json:"victorOpsNotification,omitempty"`
	WebHook      *webhook.WebHook       `json:"webHookNotification,omitempty"`
	XMatters     *xmatters.XMatters     `json:"xMattersNotification,omitempty"`
	Trello       *trello.Trello         `json:"trelloNotification,omitempty"`
	ServiceNow   *servicenow.ServiceNow `json:"serviceNowNotification,omitempty"`
}

type demofiller interface {
	FillDemoValues() []string
}

type preparer interface {
	PrepareMarshalHCL(decoder hcl.Decoder) error
}

func (me *Notification) PrepareMarshalHCL(decoder hcl.Decoder) error {
	config := me.resolveConfig()
	if preparer, ok := config.(preparer); ok {
		return preparer.PrepareMarshalHCL(decoder)
	}
	return nil
}

func (me *Notification) FillDemoValues() []string {
	config := me.resolveConfig()
	if demofiller, ok := config.(demofiller); ok {
		return demofiller.FillDemoValues()
	}
	return nil
}

func (me *Notification) Schema() map[string]*schema.Schema {
	switch me.Type {
	case Types.AnsibleTower:
		return new(ansible.AnsibleTower).Schema()
	case Types.Email:
		return new(email.Email).Schema()
	case Types.Slack:
		return new(slack.Slack).Schema()
	case Types.Jira:
		return new(jira.Jira).Schema()
	case Types.OpsGenie:
		return new(opsgenie.OpsGenie).Schema()
	case Types.PagerDuty:
		return new(pagerduty.PagerDuty).Schema()
	case Types.VictorOps:
		return new(victorops.VictorOps).Schema()
	case Types.WebHook:
		return new(webhook.WebHook).Schema()
	case Types.XMatters:
		return new(xmatters.XMatters).Schema()
	case Types.Trello:
		return new(trello.Trello).Schema()
	case Types.ServiceNow:
		return new(servicenow.ServiceNow).Schema()
	}
	panic("unsupported type " + string(me.Type))
}

func (me *Notification) resolveConfig() hcl.Unmarshaler {
	switch me.Type {
	case Types.AnsibleTower:
		return me.AnsibleTower
	case Types.Email:
		return me.Email
	case Types.Slack:
		return me.Slack
	case Types.Jira:
		return me.Jira
	case Types.OpsGenie:
		return me.OpsGenie
	case Types.PagerDuty:
		return me.PagerDuty
	case Types.VictorOps:
		return me.VictorOps
	case Types.WebHook:
		return me.WebHook
	case Types.XMatters:
		return me.XMatters
	case Types.Trello:
		return me.Trello
	case Types.ServiceNow:
		return me.ServiceNow
	}
	return nil
}

func (me *Notification) resolveAndCreateConfig() hcl.Unmarshaler {
	var config hcl.Unmarshaler
	switch me.Type {
	case Types.AnsibleTower:
		me.AnsibleTower = new(ansible.AnsibleTower)
		config = me.AnsibleTower
	case Types.Email:
		me.Email = new(email.Email)
		config = me.Email
	case Types.Slack:
		me.Slack = new(slack.Slack)
		config = me.Slack
	case Types.Jira:
		me.Jira = new(jira.Jira)
		config = me.Jira
	case Types.OpsGenie:
		me.OpsGenie = new(opsgenie.OpsGenie)
		config = me.OpsGenie
	case Types.PagerDuty:
		me.PagerDuty = new(pagerduty.PagerDuty)
		config = me.PagerDuty
	case Types.VictorOps:
		me.VictorOps = new(victorops.VictorOps)
		config = me.VictorOps
	case Types.WebHook:
		me.WebHook = new(webhook.WebHook)
		config = me.WebHook
	case Types.XMatters:
		me.XMatters = new(xmatters.XMatters)
		config = me.XMatters
	case Types.Trello:
		me.Trello = new(trello.Trello)
		config = me.Trello
	case Types.ServiceNow:
		me.ServiceNow = new(servicenow.ServiceNow)
		config = me.ServiceNow
	}
	return config
}

func (me *Notification) UnmarshalHCL(decoder hcl.Decoder) error {
	config := me.resolveAndCreateConfig()

	if err := config.UnmarshalHCL(decoder); err != nil {
		return err
	}
	return decoder.DecodeAll(map[string]any{
		"name":      &me.Name,
		"active":    &me.Enabled,
		"profile":   &me.ProfileID,
		"legacy_id": &me.LegacyID,
	})
}

func (me *Notification) MarshalHCL(properties hcl.Properties) error {
	m := map[Type]hcl.Marshaler{
		Types.AnsibleTower: me.AnsibleTower,
		Types.Email:        me.Email,
		Types.Slack:        me.Slack,
		Types.Jira:         me.Jira,
		Types.OpsGenie:     me.OpsGenie,
		Types.PagerDuty:    me.PagerDuty,
		Types.VictorOps:    me.VictorOps,
		Types.WebHook:      me.WebHook,
		Types.XMatters:     me.XMatters,
		Types.Trello:       me.Trello,
		Types.ServiceNow:   me.ServiceNow,
	}
	if len(me.Type) == 0 {
		return errors.New("no notification type set")
	}
	if config, ok := m[me.Type]; ok {
		if config == nil {
			return fmt.Errorf("notification type is `%v` but the corresponding configuration is missing", me.Type)
		}
		if err := config.MarshalHCL(properties); err != nil {
			return err
		}
		// if me.Enabled {
		// 	if err := properties.Encode("enabled", me.Enabled); err != nil { return err }
		// }
		if err := properties.Encode("name", me.Name); err != nil {
			return err
		}
		if err := properties.Encode("profile", me.ProfileID); err != nil {
			return err
		}
		if err := properties.Encode("active", me.Enabled); err != nil {
			return err
		}
		if err := properties.Encode("legacy_id", me.LegacyID); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("notification type `%v` not supported", me.Type)
}
