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

// PagerDutyConfig Configuration of the PagerDuty notification.
type PagerDutyConfig struct {
	BaseNotificationConfig
	Account       string  `json:"account"`                 // The name of the PagerDuty account.
	ServiceAPIKey *string `json:"serviceApiKey,omitempty"` // The API key to access PagerDuty.
	ServiceName   string  `json:"serviceName"`             // The name of the service.
}

func (me *PagerDutyConfig) GetType() Type {
	return Types.PagerDuty
}

func (me *PagerDutyConfig) Schema() map[string]*schema.Schema {
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
		"account": {
			Type:        schema.TypeString,
			Description: "The name of the PagerDuty account",
			Required:    true,
		},
		"service_api_key": {
			Type:        schema.TypeString,
			Description: "The API key to access PagerDuty",
			Optional:    true,
		},
		"service_name": {
			Type:        schema.TypeString,
			Description: "The name of the service",
			Required:    true,
		},
	}
}

func (me *PagerDutyConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("pager_duty.0.service_api_key"); ok && len(value.(string)) > 0 {
		me.ServiceAPIKey = opt.NewString(value.(string))
	}
	return nil
}

func (me *PagerDutyConfig) FillDemoValues() []string {
	me.ServiceAPIKey = opt.NewString("#######")
	return []string{"The REST API didn't provide the credentials"}
}

func (me *PagerDutyConfig) MarshalHCL(properties hcl.Properties) error {
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
	if err := properties.Encode("account", me.Account); err != nil {
		return err
	}
	if err := properties.Encode("service_api_key", me.ServiceAPIKey); err != nil {
		return err
	}

	if err := properties.Encode("service_name", me.ServiceName); err != nil {
		return err
	}

	return nil
}

func (me *PagerDutyConfig) UnmarshalHCL(decoder hcl.Decoder) error {
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
		delete(me.Unknowns, "account")
		delete(me.Unknowns, "service_api_key")
		delete(me.Unknowns, "service_name")
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
	if value, ok := decoder.GetOk("account"); ok {
		me.Account = value.(string)
	}
	if value, ok := decoder.GetOk("service_api_key"); ok {
		me.ServiceAPIKey = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("service_name"); ok {
		me.ServiceName = value.(string)
	}
	return nil
}

func (me *PagerDutyConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":              me.ID,
		"name":            me.Name,
		"type":            me.GetType(),
		"active":          me.Active,
		"alertingProfile": me.AlertingProfile,
		"account":         me.Account,
		"serviceApiKey":   me.ServiceAPIKey,
		"serviceName":     me.ServiceName,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *PagerDutyConfig) UnmarshalJSON(data []byte) error {
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
		"account":         &me.Account,
		"serviceApiKey":   &me.ServiceAPIKey,
		"serviceName":     &me.ServiceName,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
