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

package azure

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AzureCredentials Configuration of Azure app-level credentials.
type AzureCredentials struct {
	Label                        string                     `json:"label"`                        // The unique name of the Azure credentials configuration.  Allowed characters are letters, numbers, and spaces. Also the special characters `.+-_` are allowed.
	DirectoryID                  string                     `json:"directoryId"`                  // The Directory ID (also referred to as Tenant ID)  The combination of Application ID and Directory ID must be unique.
	AutoTagging                  *bool                      `json:"autoTagging"`                  // The automatic capture of Azure tags is on (`true`) or off (`false`).
	MonitorOnlyTagPairs          []*CloudTag                `json:"monitorOnlyTagPairs"`          // A list of Azure tags to be monitored.  You can specify up to 10 tags. A resource tagged with *any* of the specified tags is monitored.  Only applicable when the **monitorOnlyTaggedEntities** parameter is set to `true`.
	MonitorOnlyExcludingTagPairs []*CloudTag                `json:"monitorOnlyExcludingTagPairs"` // A list of Azure tags to be excluded from monitoring.  You can specify up to 10 tags. A resource tagged with *any* of the specified tags is monitored.  Only applicable when the **monitorOnlyTaggedEntities** parameter is set to `true`.
	Active                       *bool                      `json:"active,omitempty"`             // The monitoring is enabled (`true`) or disabled (`false`).  If not set on creation, the `true` value is used.  If the field is omitted during an update, the old value remains unaffected.
	AppID                        string                     `json:"appId"`                        // The Application ID (also referred to as Client ID)  The combination of Application ID and Directory ID must be unique.
	Key                          *string                    `json:"key,omitempty"`                // The secret key associated with the Application ID.  For security reasons, GET requests return this field as `null`.   Submit your key on creation or update of the configuration. If the field is omitted during an update, the old value remains unaffected.
	MonitorOnlyTaggedEntities    *bool                      `json:"monitorOnlyTaggedEntities"`    // Monitor only resources that have specified Azure tags (`true`) or all resources (`false`).
	SupportingServices           []*AzureSupportingService  `json:"supportingServices,omitempty"` // A list of Azure supporting services to be monitored. For each service there's a sublist of its metrics and the metrics' dimensions that should be monitored. All of these elements (services, metrics, dimensions) must have corresponding static definitions on the server.
	Unknowns                     map[string]json.RawMessage `json:"-"`
}

func (ac *AzureCredentials) Name() string {
	return ac.Label
}

func (ac *AzureCredentials) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"label": {
			Type:        schema.TypeString,
			Description: "The unique name of the Azure credentials configuration.  Allowed characters are letters, numbers, and spaces. Also the special characters `.+-_` are allowed",
			Optional:    true,
		},
		"directory_id": {
			Type:        schema.TypeString,
			Description: "The Directory ID (also referred to as Tenant ID)  The combination of Application ID and Directory ID must be unique",
			Optional:    true,
		},
		"app_id": {
			Type:        schema.TypeString,
			Description: "The Application ID (also referred to as Client ID)  The combination of Application ID and Directory ID must be unique",
			Optional:    true,
		},
		"auto_tagging": {
			Type:        schema.TypeBool,
			Description: "The automatic capture of Azure tags is on (`true`) or off (`false`)",
			Optional:    true,
		},
		"monitor_only_tag_pairs": {
			Type:        schema.TypeList,
			Description: "A list of Azure tags to be monitored.  You can specify up to 20 tags. A resource tagged with *any* of the specified tags is monitored.  Only applicable when the **monitorOnlyTaggedEntities** parameter is set to `true`",
			Optional:    true,
			MaxItems:    20,
			Elem: &schema.Resource{
				Schema: new(CloudTag).Schema(),
			},
		},
		"monitor_only_excluding_tag_pairs": {
			Type:        schema.TypeList,
			Description: "A list of Azure tags to be excluded from monitoring.  You can specify up to 20 tags. A resource tagged with *any* of the specified tags is monitored.  Only applicable when the **monitorOnlyTaggedEntities** parameter is set to `true`.",
			Optional:    true,
			MaxItems:    20,
			Elem: &schema.Resource{
				Schema: new(CloudTag).Schema(),
			},
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "The monitoring is enabled (`true`) or disabled (`false`).  If not set on creation, the `true` value is used.  If the field is omitted during an update, the old value remains unaffected",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The secret key associated with the Application ID.  For security reasons, GET requests return this field as `null`. Submit your key on creation or update of the configuration. If the field is omitted during an update, the old value remains unaffected.",
			Optional:    true,
			Sensitive:   true,
		},
		"monitor_only_tagged_entities": {
			Type:        schema.TypeBool,
			Description: "Monitor only resources that have specified Azure tags (`true`) or all resources (`false`).",
			Required:    true,
		},
		"supporting_services": {
			Type:        schema.TypeList,
			Description: "A list of Azure supporting services to be monitored. For each service there's a sublist of its metrics and the metrics' dimensions that should be monitored. All of these elements (services, metrics, dimensions) must have corresponding static definitions on the server.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(AzureSupportingService).Schema(),
			},
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (ac *AzureCredentials) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(ac.Unknowns) > 0 {
		delete(ac.Unknowns, "id")
		delete(ac.Unknowns, "metadata")
		for k, v := range ac.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(ac.Label)
		if err != nil {
			return nil, err
		}
		m["label"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(ac.DirectoryID)
		if err != nil {
			return nil, err
		}
		m["directoryId"] = rawMessage
	}
	if rawMessage, err := json.Marshal(opt.Bool(ac.AutoTagging)); err == nil {
		m["autoTagging"] = rawMessage
	} else {
		return nil, err
	}
	if ac.MonitorOnlyTagPairs != nil {
		rawMessage, err := json.Marshal(ac.MonitorOnlyTagPairs)
		if err != nil {
			return nil, err
		}
		m["monitorOnlyTagPairs"] = rawMessage
	} else {
		rawMessage, err := json.Marshal([]*CloudTag{})
		if err != nil {
			return nil, err
		}
		m["monitorOnlyTagPairs"] = rawMessage
	}

	if ac.MonitorOnlyExcludingTagPairs != nil {
		rawMessage, err := json.Marshal(ac.MonitorOnlyExcludingTagPairs)
		if err != nil {
			return nil, err
		}
		m["monitorOnlyExcludingTagPairs"] = rawMessage
	} else {
		rawMessage, err := json.Marshal([]*CloudTag{})
		if err != nil {
			return nil, err
		}
		m["monitorOnlyExcludingTagPairs"] = rawMessage
	}

	if rawMessage, err := json.Marshal(opt.Bool(ac.Active)); err == nil {
		m["active"] = rawMessage
	} else {
		return nil, err
	}

	if rawMessage, err := json.Marshal(ac.AppID); err == nil {
		m["appId"] = rawMessage
	} else {
		return nil, err
	}

	if ac.Key != nil {
		if rawMessage, err := json.Marshal(ac.Key); err == nil {
			m["key"] = rawMessage
		} else {
			return nil, err
		}
	}
	if rawMessage, err := json.Marshal(opt.Bool(ac.MonitorOnlyTaggedEntities)); err == nil {
		m["monitorOnlyTaggedEntities"] = rawMessage
	} else {
		return nil, err
	}
	if ac.SupportingServices != nil {
		rawMessage, err := json.Marshal(ac.SupportingServices)
		if err != nil {
			return nil, err
		}
		m["supportingServices"] = rawMessage
	}
	return json.Marshal(m)
}

func (ac *AzureCredentials) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(ac.Unknowns, "id")
	delete(ac.Unknowns, "metadata")
	if v, found := m["label"]; found {
		if err := json.Unmarshal(v, &ac.Label); err != nil {
			return err
		}
	}
	if v, found := m["directoryId"]; found {
		if err := json.Unmarshal(v, &ac.DirectoryID); err != nil {
			return err
		}
	}
	if v, found := m["autoTagging"]; found {
		if err := json.Unmarshal(v, &ac.AutoTagging); err != nil {
			return err
		}
	}
	if v, found := m["monitorOnlyTagPairs"]; found {
		if err := json.Unmarshal(v, &ac.MonitorOnlyTagPairs); err != nil {
			return err
		}
	}
	if v, found := m["monitorOnlyExcludingTagPairs"]; found {
		if err := json.Unmarshal(v, &ac.MonitorOnlyExcludingTagPairs); err != nil {
			return err
		}
	}

	if v, found := m["active"]; found {
		if err := json.Unmarshal(v, &ac.Active); err != nil {
			return err
		}
	}
	if v, found := m["appId"]; found {
		if err := json.Unmarshal(v, &ac.AppID); err != nil {
			return err
		}
	}
	if v, found := m["key"]; found {
		if err := json.Unmarshal(v, &ac.Key); err != nil {
			return err
		}
	}
	if v, found := m["monitorOnlyTaggedEntities"]; found {
		if err := json.Unmarshal(v, &ac.MonitorOnlyTaggedEntities); err != nil {
			return err
		}
	}
	if v, found := m["supportingServices"]; found {
		if err := json.Unmarshal(v, &ac.SupportingServices); err != nil {
			return err
		}
	}
	delete(m, "id")
	delete(m, "label")
	delete(m, "directoryId")
	delete(m, "autoTagging")
	delete(m, "metadata")
	delete(m, "monitorOnlyTagPairs")
	delete(m, "monitorOnlyExcludingTagPairs")
	delete(m, "active")
	delete(m, "appId")
	delete(m, "key")
	delete(m, "monitorOnlyTaggedEntities")
	delete(m, "supportingServices")

	if len(m) > 0 {
		ac.Unknowns = m
	}
	return nil
}

func (ac *AzureCredentials) MarshalHCL(properties hcl.Properties) error {
	if len(ac.Unknowns) > 0 {
		delete(ac.Unknowns, "id")
		delete(ac.Unknowns, "metadata")
		data, err := json.Marshal(ac.Unknowns)
		if err != nil {
			return err
		}
		if err := properties.Encode("unknowns", string(data)); err != nil {
			return err
		}
	}

	if err := properties.Encode("label", ac.Label); err != nil {
		return err
	}
	if err := properties.Encode("directory_id", ac.DirectoryID); err != nil {
		return err
	}
	if err := properties.Encode("auto_tagging", opt.Bool(ac.AutoTagging)); err != nil {
		return err
	}
	if err := properties.Encode("monitor_only_tag_pairs", ac.MonitorOnlyTagPairs); err != nil {
		return err
	}
	if err := properties.Encode("monitor_only_excluding_tag_pairs", ac.MonitorOnlyExcludingTagPairs); err != nil {
		return err
	}
	if err := properties.Encode("active", opt.Bool(ac.Active)); err != nil {
		return err
	}
	if err := properties.Encode("app_id", ac.AppID); err != nil {
		return err
	}
	if err := properties.Encode("key", "${state.secret_value}"); err != nil {
		return err
	}
	if err := properties.Encode("monitor_only_tagged_entities", opt.Bool(ac.MonitorOnlyTaggedEntities)); err != nil {
		return err
	}
	if err := properties.Encode("supporting_services", ac.SupportingServices); err != nil {
		return err
	}
	return nil
}

func (ac *AzureCredentials) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), ac); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &ac.Unknowns); err != nil {
			return err
		}
		delete(ac.Unknowns, "label")
		delete(ac.Unknowns, "directory_id")
		delete(ac.Unknowns, "auto_tagging")
		delete(ac.Unknowns, "monitor_only_tag_pairs")
		delete(ac.Unknowns, "monitorOnlyExcludingTagPairs")
		delete(ac.Unknowns, "active")
		delete(ac.Unknowns, "app_id")
		delete(ac.Unknowns, "key")
		delete(ac.Unknowns, "monitor_only_tagged_entities")
		delete(ac.Unknowns, "supporting_services")
		delete(ac.Unknowns, "id")
		delete(ac.Unknowns, "metadata")
		if len(ac.Unknowns) == 0 {
			ac.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("label"); ok {
		ac.Label = value.(string)
	}
	if value, ok := decoder.GetOk("directory_id"); ok {
		ac.DirectoryID = value.(string)
	}
	if value, ok := decoder.GetOk("auto_tagging"); ok {
		ac.AutoTagging = opt.NewBool(value.(bool))
	}
	if result, ok := decoder.GetOk("monitor_only_tag_pairs.#"); ok {
		ac.MonitorOnlyTagPairs = []*CloudTag{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(CloudTag)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "monitor_only_tag_pairs", idx)); err != nil {
				return err
			}
			ac.MonitorOnlyTagPairs = append(ac.MonitorOnlyTagPairs, entry)
		}
	}
	if result, ok := decoder.GetOk("monitor_only_excluding_tag_pairs.#"); ok {
		ac.MonitorOnlyExcludingTagPairs = []*CloudTag{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(CloudTag)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "monitor_only_excluding_tag_pairs", idx)); err != nil {
				return err
			}
			ac.MonitorOnlyExcludingTagPairs = append(ac.MonitorOnlyExcludingTagPairs, entry)
		}
	}
	if value, ok := decoder.GetOk("active"); ok {
		ac.Active = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("app_id"); ok {
		ac.AppID = value.(string)
	}
	if value, ok := decoder.GetOk("key"); ok {
		ac.Key = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("monitor_only_tagged_entities"); ok {
		ac.MonitorOnlyTaggedEntities = opt.NewBool(value.(bool))
	}
	if result, ok := decoder.GetOk("supporting_services.#"); ok {
		ac.SupportingServices = []*AzureSupportingService{}
		for idx := 0; idx < result.(int); idx++ {
			entry := new(AzureSupportingService)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "supporting_services", idx)); err != nil {
				return err
			}
			ac.SupportingServices = append(ac.SupportingServices, entry)
		}
	}
	return nil
}

const credsNotProvided = "REST API didn't provide credential data"

func (me *AzureCredentials) FillDemoValues() []string {
	me.Key = opt.NewString("################")
	return []string{credsNotProvided}
}
