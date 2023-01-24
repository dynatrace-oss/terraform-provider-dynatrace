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

package requestattributes

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ScopeConditions Conditions for data capturing.
type ScopeConditions struct {
	HostGroup         *string                    `json:"hostGroup,omitempty"`         // Only applies to this host group.
	ProcessGroup      *string                    `json:"processGroup,omitempty"`      // Only applies to this process group. Note that this can't be transferred between different clusters or environments.
	ServiceTechnology *ServiceTechnology         `json:"serviceTechnology,omitempty"` // Only applies to this service technology.
	TagOfProcessGroup *string                    `json:"tagOfProcessGroup,omitempty"` // Only apply to process groups matching this tag.
	Unknowns          map[string]json.RawMessage `json:"-"`
}

func (me *ScopeConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host_group": {
			Type:        schema.TypeString,
			Description: "Only applies to this host group",
			Optional:    true,
		},
		"process_group": {
			Type:        schema.TypeString,
			Description: "Only applies to this process group. Note that this can't be transferred between different clusters or environments",
			Optional:    true,
		},
		"service_technology": {
			Type:        schema.TypeString,
			Description: "Only applies to this service technology",
			Optional:    true,
		},
		"tag_of_process_group": {
			Type:        schema.TypeString,
			Description: "Only apply to process groups matching this tag",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ScopeConditions) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("host_group", me.HostGroup); err != nil {
		return err
	}
	if err := properties.Encode("process_group", me.ProcessGroup); err != nil {
		return err
	}
	if err := properties.Encode("service_technology", me.ServiceTechnology); err != nil {
		return err
	}
	if err := properties.Encode("tag_of_process_group", me.TagOfProcessGroup); err != nil {
		return err
	}
	return nil
}

func (me *ScopeConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "host_group")
		delete(me.Unknowns, "process_group")
		delete(me.Unknowns, "service_technology")
		delete(me.Unknowns, "tag_of_process_group")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("host_group"); ok {
		me.HostGroup = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("process_group"); ok {
		me.ProcessGroup = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("service_technology"); ok {
		me.ServiceTechnology = ServiceTechnology(value.(string)).Ref()
	}
	if value, ok := decoder.GetOk("tag_of_process_group"); ok {
		me.TagOfProcessGroup = opt.NewString(value.(string))
	}
	return nil
}

func (me *ScopeConditions) MarshalJSON() ([]byte, error) {
	m := xjson.NewProperties(me.Unknowns)
	if err := m.Marshal("hostGroup", me.HostGroup); err != nil {
		return nil, err
	}
	if err := m.Marshal("processGroup", me.ProcessGroup); err != nil {
		return nil, err
	}
	if err := m.Marshal("serviceTechnology", me.ServiceTechnology); err != nil {
		return nil, err
	}
	if err := m.Marshal("tagOfProcessGroup", me.TagOfProcessGroup); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (me *ScopeConditions) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("hostGroup", &me.HostGroup); err != nil {
		return err
	}
	if err := m.Unmarshal("processGroup", &me.ProcessGroup); err != nil {
		return err
	}
	if err := m.Unmarshal("serviceTechnology", &me.ServiceTechnology); err != nil {
		return err
	}
	if err := m.Unmarshal("tagOfProcessGroup", &me.TagOfProcessGroup); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
