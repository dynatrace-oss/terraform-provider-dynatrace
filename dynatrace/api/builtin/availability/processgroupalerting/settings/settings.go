/**
* @license
* Copyright 2025 Dynatrace LLC
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

package processgroupalerting

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled                  bool          `json:"enabled"`                            // Enable (`true`) or disable (`false`) process group availability monitoring
	AlertingMode             *AlertingMode `json:"alertingMode,omitempty"`             // Possible Values: `ON_INSTANCE_COUNT_VIOLATION`, `ON_PGI_UNAVAILABILITY`
	MinimumInstanceThreshold *int          `json:"minimumInstanceThreshold,omitempty"` // Open a new problem if the number of active process instances in the group is fewer than X
	ProcessGroupId           string        `json:"-"`
}

func (me *Settings) Name() string {
	return me.ProcessGroupId
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"process_group": {
			Type:        schema.TypeString,
			Description: "The process group ID for availability monitoring",
			Required:    true,
			ForceNew:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Enable (`true`) or disable (`false`) process group availability monitoring",
			Required:    true,
		},
		"alerting_mode": {
			Type:        schema.TypeString,
			Description: "Possible Values: `ON_INSTANCE_COUNT_VIOLATION`, `ON_PGI_UNAVAILABILITY`",
			Optional:    true,
		},
		"minimum_instance_threshold": {
			Type:        schema.TypeInt,
			Description: "Open a new problem if the number of active process instances in the group is fewer than X",
			Optional:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"process_group":              me.ProcessGroupId,
		"enabled":                    me.Enabled,
		"alerting_mode":              me.AlertingMode,
		"minimum_instance_threshold": me.MinimumInstanceThreshold,
	}); err != nil {
		return err
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"process_group":              &me.ProcessGroupId,
		"enabled":                    &me.Enabled,
		"alerting_mode":              &me.AlertingMode,
		"minimum_instance_threshold": &me.MinimumInstanceThreshold,
	})
}

func (me *Settings) SetScope(scope string) {
	me.ProcessGroupId = scope
}

func (me *Settings) GetScope() string {
	return me.ProcessGroupId
}

func (me *Settings) Store() ([]byte, error) {
	var data []byte
	var err error
	if data, err = json.Marshal(me); err != nil {
		return nil, err
	}
	m := map[string]json.RawMessage{}
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	if data, err = json.Marshal(me.ProcessGroupId); err != nil {
		return nil, err
	}
	m["process_group"] = data
	return json.MarshalIndent(m, "", "  ")
}

func (me *Settings) Load(data []byte) error {
	if err := json.Unmarshal(data, &me); err != nil {
		return err
	}

	c := struct {
		ProcessGroupId string `json:"process_group"`
	}{}
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	me.ProcessGroupId = c.ProcessGroupId

	return nil
}
