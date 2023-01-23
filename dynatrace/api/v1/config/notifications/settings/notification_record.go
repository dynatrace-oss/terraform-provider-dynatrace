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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NotificationRecord struct {
	NotificationConfig NotificationConfig `json:"-"`
}

func (me *NotificationRecord) Name() string {
	return me.NotificationConfig.GetName()
}

func (me *NotificationRecord) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ansible_tower": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for Ansible Tower Notification",
			Elem:        &schema.Resource{Schema: new(AnsibleTowerConfig).Schema()},
		},
		"email": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for Email Notification",
			Elem:        &schema.Resource{Schema: new(EmailConfig).Schema()},
		},
		"jira": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for Jira Notification",
			Elem:        &schema.Resource{Schema: new(JiraConfig).Schema()},
		},
		"ops_genie": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for OpsGenie Notification",
			Elem:        &schema.Resource{Schema: new(OpsGenieConfig).Schema()},
		},
		"pager_duty": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for PagerDuty Notification",
			Elem:        &schema.Resource{Schema: new(PagerDutyConfig).Schema()},
		},
		"service_now": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for ServiceNow Notification",
			Elem:        &schema.Resource{Schema: new(ServiceNowConfig).Schema()},
		},
		"slack": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for Slack Notification",
			Elem:        &schema.Resource{Schema: new(SlackConfig).Schema()},
		},
		"trello": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for Trello Notification",
			Elem:        &schema.Resource{Schema: new(TrelloConfig).Schema()},
		},
		"victor_ops": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for VictorOps Notification",
			Elem:        &schema.Resource{Schema: new(VictorOpsConfig).Schema()},
		},
		"web_hook": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for WebHook Notification",
			Elem:        &schema.Resource{Schema: new(WebHookConfig).Schema()},
		},
		"xmatters": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for XMatters Notification",
			Elem:        &schema.Resource{Schema: new(XMattersConfig).Schema()},
		},
		"config": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration for Generic Notification",
			Elem:        &schema.Resource{Schema: new(BaseNotificationConfig).Schema()},
		},
	}
}

func (me *NotificationRecord) MarshalHCL(properties hcl.Properties) error {
	if me.NotificationConfig != nil {
		switch config := me.NotificationConfig.(type) {
		case *AnsibleTowerConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["ansible_tower"] = []any{marshalled}
			} else {
				return nil
			}
		case *EmailConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["email"] = []any{marshalled}
			} else {
				return nil
			}
		case *JiraConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["jira"] = []any{marshalled}
			} else {
				return nil
			}
		case *OpsGenieConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["ops_genie"] = []any{marshalled}
			} else {
				return nil
			}
		case *PagerDutyConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["pager_duty"] = []any{marshalled}
			} else {
				return nil
			}
		case *ServiceNowConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["service_now"] = []any{marshalled}
			} else {
				return nil
			}
		case *SlackConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["slack"] = []any{marshalled}
			} else {
				return nil
			}
		case *TrelloConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["trello"] = []any{marshalled}
			} else {
				return nil
			}
		case *VictorOpsConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["victor_ops"] = []any{marshalled}
			} else {
				return nil
			}
		case *WebHookConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["web_hook"] = []any{marshalled}
			} else {
				return nil
			}
		case *XMattersConfig:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["xmatters"] = []any{marshalled}
			} else {
				return nil
			}
		default:
			marshalled := hcl.Properties{}

			if err := config.MarshalHCL(marshalled); err == nil {
				properties["config"] = []any{marshalled}
			} else {
				return nil
			}
		}
	}
	return nil
}

type MarshalPreparer interface {
	PrepareMarshalHCL(hcl.Decoder) error
}

func (me *NotificationRecord) PrepareMarshalHCL(d hcl.Decoder) error {
	if me.NotificationConfig != nil {
		if preparer, ok := me.NotificationConfig.(MarshalPreparer); ok {
			return preparer.PrepareMarshalHCL(d)
		}
	}
	return nil
}

func (me *NotificationRecord) FillDemoValues() []string {
	if me.NotificationConfig != nil {
		if demo, ok := me.NotificationConfig.(settings.DemoSettings); ok {
			return demo.FillDemoValues()
		}
	}
	return []string{}
}

func (me *NotificationRecord) UnmarshalHCL(decoder hcl.Decoder) error {
	if _, ok := decoder.GetOk("ansible_tower.#"); ok {
		cfg := new(AnsibleTowerConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "ansible_tower", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("email.#"); ok {
		cfg := new(EmailConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "email", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("jira.#"); ok {
		cfg := new(JiraConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "jira", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("ops_genie.#"); ok {
		cfg := new(OpsGenieConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "ops_genie", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("pager_duty.#"); ok {
		cfg := new(PagerDutyConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "pager_duty", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("service_now.#"); ok {
		cfg := new(ServiceNowConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "service_now", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("slack.#"); ok {
		cfg := new(SlackConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "slack", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("trello.#"); ok {
		cfg := new(TrelloConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "trello", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("victor_ops.#"); ok {
		cfg := new(VictorOpsConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "victor_ops", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("web_hook.#"); ok {
		cfg := new(WebHookConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "web_hook", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("xmatters.#"); ok {
		cfg := new(XMattersConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "xmatters", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	if _, ok := decoder.GetOk("config.#"); ok {
		cfg := new(BaseNotificationConfig)
		if err := cfg.UnmarshalHCL(hcl.NewDecoder(decoder, "config", 0)); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	}
	return nil
}

func (me *NotificationRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.NotificationConfig)
}

func (me *NotificationRecord) UnmarshalJSON(data []byte) error {
	config := new(BaseNotificationConfig)
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}
	switch config.Type {
	case Types.Ansibletower:
		cfg := new(AnsibleTowerConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Email:
		cfg := new(EmailConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Jira:
		cfg := new(JiraConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.OpsGenie:
		cfg := new(OpsGenieConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.PagerDuty:
		cfg := new(PagerDutyConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.ServiceNow:
		cfg := new(ServiceNowConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Slack:
		cfg := new(SlackConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Trello:
		cfg := new(TrelloConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Victorops:
		cfg := new(VictorOpsConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Webhook:
		cfg := new(WebHookConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	case Types.Xmatters:
		cfg := new(XMattersConfig)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		me.NotificationConfig = cfg
	default:
		me.NotificationConfig = config
	}
	return nil
}
