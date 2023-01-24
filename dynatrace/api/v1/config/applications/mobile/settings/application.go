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

package mobile

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings/collections"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Application represents configuration of a mobile or custom application to be created
type Application struct {
	Name                             string                         `json:"name"`                                       // The name of the application
	ApplicationType                  *ApplicationType               `json:"applicationType,omitempty"`                  // The type of the application
	ApplicationID                    *string                        `json:"applicationId,omitempty"`                    // The UUID of the application.\n\nIt is used only by OneAgent to send data to Dynatrace
	CostControlUserSessionPercentage *int32                         `json:"costControlUserSessionPercentage,omitempty"` // The percentage of user sessions to be analyzed
	ApdexSettings                    *MobileCustomApdex             `json:"apdexSettings,omitempty"`
	OptInModeEnabled                 bool                           `json:"optInModeEnabled,omitempty"`     // The opt-in mode is enabled (`true`) or disabled (`false`).\n\nThis value is only applicable to mobile and not to custom apps
	SessionReplayEnabled             bool                           `json:"sessionReplayEnabled,omitempty"` // The session replay is enabled (`true`) or disabled (`false`).\nThis value is only applicable to mobile and not to custom apps
	SessionReplayOnCrashEnabled      bool                           `json:"sessionReplayOnCrashEnabled"`    // The session replay on crash is enabled (`true`) or disabled (`false`). \n\nEnabling requires both **sessionReplayEnabled** and **optInModeEnabled** values set to `true`.\nAlso, this value is only applicable to mobile and not to custom apps
	BeaconEndpointType               BeaconEndpointType             `json:"beaconEndpointType"`             // The type of the beacon endpoint
	BeaconEndpointUrl                *string                        `json:"beaconEndpointUrl,omitempty"`    // The URL of the beacon endpoint.\n\nOnly applicable when the **beaconEndpointType** is set to `ENVIRONMENT_ACTIVE_GATE` or `INSTRUMENTED_WEB_SERVER`
	KeyUserActions                   collections.Set[string]        `json:"-"`
	Properties                       UserActionAndSessionProperties `json:"-"`
}

func (me *Application) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("name", &me.Name); err != nil {
		return err
	}
	if err := decoder.Decode("application_type", &me.ApplicationType); err != nil {
		return err
	}
	if err := decoder.Decode("application_id", &me.ApplicationID); err != nil {
		return err
	}
	if err := decoder.Decode("user_session_percentage", &me.CostControlUserSessionPercentage); err != nil {
		return err
	}
	if err := decoder.Decode("opt_in_mode", &me.OptInModeEnabled); err != nil {
		return err
	}
	if err := decoder.Decode("session_replay", &me.SessionReplayEnabled); err != nil {
		return err
	}
	if err := decoder.Decode("session_replay_on_crash", &me.SessionReplayOnCrashEnabled); err != nil {
		return err
	}
	if err := decoder.Decode("beacon_endpoint_type", &me.BeaconEndpointType); err != nil {
		return err
	}
	if err := decoder.Decode("beacon_endpoint_url", &me.BeaconEndpointUrl); err != nil {
		return err
	}
	if err := decoder.Decode("apdex", &me.ApdexSettings); err != nil {
		return err
	}
	if err := decoder.Decode("key_user_actions", &me.KeyUserActions); err != nil {
		return err
	}
	if err := decoder.Decode("properties", &me.Properties); err != nil {
		return err
	}
	return nil
}

func (me *Application) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":                    me.Name,
		"application_id":          me.ApplicationID,
		"application_type":        me.ApplicationType,
		"user_session_percentage": me.CostControlUserSessionPercentage,
		"opt_in_mode":             me.OptInModeEnabled,
		"session_replay":          me.SessionReplayEnabled,
		"session_replay_on_crash": me.SessionReplayOnCrashEnabled,
		"beacon_endpoint_type":    me.BeaconEndpointType,
		"beacon_endpoint_url":     me.BeaconEndpointUrl,
		"apdex":                   me.ApdexSettings,
		"key_user_actions":        me.KeyUserActions,
		"properties":              me.Properties,
	})
}

func (me *Application) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the application",
			Required:    true,
		},
		"application_id": {
			Type:        schema.TypeString,
			Description: "The UUID of the application.\n\nIt is used only by OneAgent to send data to Dynatrace. If not specified it will get generated.",
			Optional:    true,
		},
		"application_type": {
			Type:        schema.TypeString,
			Description: "The type of the application. Either `CUSTOM_APPLICATION` or `MOBILE_APPLICATION`.",
			Optional:    true,
		},
		"user_session_percentage": {
			Type:        schema.TypeInt,
			Description: "The percentage of user sessions to be analyzed",
			Optional:    true,
		},
		"opt_in_mode": {
			Type:        schema.TypeBool,
			Description: "The opt-in mode is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"session_replay": {
			Type:        schema.TypeBool,
			Description: "The session replay is enabled (`true`) or disabled (`false`).",
			Optional:    true,
		},
		"session_replay_on_crash": {
			Type:        schema.TypeBool,
			Description: "The session replay on crash is enabled (`true`) or disabled (`false`). \n\nEnabling requires both **sessionReplayEnabled** and **optInModeEnabled** values set to `true`.",
			Optional:    true,
		},
		"beacon_endpoint_type": {
			Type:        schema.TypeString,
			Description: "The type of the beacon endpoint. Possible values are `CLUSTER_ACTIVE_GATE`, `ENVIRONMENT_ACTIVE_GATE` and `INSTRUMENTED_WEB_SERVER`.",
			Required:    true,
		},
		"beacon_endpoint_url": {
			Type:        schema.TypeString,
			Description: "The URL of the beacon endpoint.\n\nOnly applicable when the **beacon_endpoint_type** is set to `ENVIRONMENT_ACTIVE_GATE` or `INSTRUMENTED_WEB_SERVER`",
			Optional:    true,
		},
		"key_user_actions": {
			Type:        schema.TypeSet,
			Description: "User Action names to be flagged as Key User Actions",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
		},
		"apdex": {
			Type:        schema.TypeList,
			Description: "Apdex configuration of a mobile application. A duration less than the **tolerable** threshold is considered satisfied",
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(MobileCustomApdex).Schema()},
		},
		"properties": {
			Type:        schema.TypeList,
			Description: "User Action and Session Properties",
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: new(UserActionAndSessionProperties).Schema()},
		},
	}
}

func (me *Application) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if len(me.KeyUserActions) > 0 {
		if data, err = json.Marshal(me.KeyUserActions); err != nil {
			return nil, err
		}
		m["keyUserActions"] = data
	}
	if len(me.Properties) > 0 {
		if data, err = json.Marshal(me.Properties); err != nil {
			return nil, err
		}
		m["properties"] = data
	}
	return json.MarshalIndent(m, "", "  ")
}

func (me *Application) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		KeyUserActions collections.Set[string]        `json:"keyUserActions"`
		Properties     UserActionAndSessionProperties `json:"properties"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.KeyUserActions = c.KeyUserActions
	me.Properties = c.Properties

	return nil
}
