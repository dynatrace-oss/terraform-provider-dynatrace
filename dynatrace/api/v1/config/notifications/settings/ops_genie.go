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

// OpsGenieConfig Configuration of the OpsGenie notification.
type OpsGenieConfig struct {
	BaseNotificationConfig
	APIKey  *string `json:"apiKey,omitempty"` // The API key to access OpsGenie.
	Domain  string  `json:"domain"`           // The region domain of the OpsGenie.
	Message string  `json:"message"`          // The content of the message.  You can use the following placeholders:  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem.
}

func (me *OpsGenieConfig) GetType() Type {
	return Types.OpsGenie
}

func (me *OpsGenieConfig) Schema() map[string]*schema.Schema {
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
		"api_key": {
			Type:        schema.TypeString,
			Description: "The API key to access OpsGenie",
			Optional:    true,
		},
		"domain": {
			Type:        schema.TypeString,
			Description: "The region domain of the OpsGenie",
			Required:    true,
		},
		"message": {
			Type:        schema.TypeString,
			Description: "The content of the message.  You can use the following placeholders:  * `{ProblemID}`: The display number of the reported problem.  * `{ProblemImpact}`: The [impact level](https://www.dynatrace.com/support/help/shortlink/impact-analysis) of the problem. Possible values are `APPLICATION`, `SERVICE`, and `INFRASTRUCTURE`.  * `{ProblemSeverity}`: The [severity level](https://www.dynatrace.com/support/help/shortlink/event-types) of the problem. Possible values are `AVAILABILITY`, `ERROR`, `PERFORMANCE`, `RESOURCE_CONTENTION`, and `CUSTOM_ALERT`.  * `{ProblemTitle}`: A short description of the problem",
			Required:    true,
		},
	}
}

func (me *OpsGenieConfig) PrepareMarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("ops_genie.0.api_key"); ok && len(value.(string)) > 0 {
		me.APIKey = opt.NewString(value.(string))
	}
	return nil
}

func (me *OpsGenieConfig) FillDemoValues() []string {
	me.APIKey = opt.NewString("#######")
	return []string{"The REST API didn't provide the credentials"}
}

func (me *OpsGenieConfig) MarshalHCL(properties hcl.Properties) error {
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
	if err := properties.Encode("api_key", me.APIKey); err != nil {
		return err
	}

	if err := properties.Encode("domain", me.Domain); err != nil {
		return err
	}
	if err := properties.Encode("message", me.Message); err != nil {
		return err
	}
	return nil
}

func (me *OpsGenieConfig) UnmarshalHCL(decoder hcl.Decoder) error {
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
		delete(me.Unknowns, "api_key")
		delete(me.Unknowns, "domain")
		delete(me.Unknowns, "message")
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
	if value, ok := decoder.GetOk("api_key"); ok {
		me.APIKey = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("domain"); ok {
		me.Domain = value.(string)
	}
	if value, ok := decoder.GetOk("message"); ok {
		me.Message = value.(string)
	}

	return nil
}

func (me *OpsGenieConfig) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"id":              me.ID,
		"name":            me.Name,
		"type":            me.GetType(),
		"active":          me.Active,
		"alertingProfile": me.AlertingProfile,
		"apiKey":          me.APIKey,
		"domain":          me.Domain,
		"message":         me.Message,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *OpsGenieConfig) UnmarshalJSON(data []byte) error {
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
		"apiKey":          &me.APIKey,
		"domain":          &me.Domain,
		"message":         &me.Message,
	}); err != nil {
		return err
	}
	if len(properties) > 0 {
		me.Unknowns = properties
	}
	return nil
}
