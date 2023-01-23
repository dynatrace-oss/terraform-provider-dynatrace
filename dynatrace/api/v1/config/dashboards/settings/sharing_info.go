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

package dashboards

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SharingInfo represents sharing configuration of a dashboard
type SharingInfo struct {
	LinkShared *bool                      `json:"linkShared,omitempty"` // If `true`, the dashboard is shared via link and authenticated users with the link can view
	Published  *bool                      `json:"published,omitempty"`  // If `true`, the dashboard is published to anyone on this environment
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (me *SharingInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"link_shared": {
			Type:        schema.TypeBool,
			Description: "If `true`, the dashboard is shared via link and authenticated users with the link can view",
			Optional:    true,
		},
		"published": {
			Type:        schema.TypeBool,
			Description: "If `true`, the dashboard is published to anyone on this environment",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *SharingInfo) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	if err := properties.Encode("link_shared", opt.Bool(me.LinkShared)); err != nil {
		return err
	}
	if err := properties.Encode("published", opt.Bool(me.Published)); err != nil {
		return err
	}
	return nil
}

func (me *SharingInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "link_shared")
		delete(me.Unknowns, "published")
		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("link_shared"); ok {
		me.LinkShared = opt.NewBool(value.(bool))
	}
	if value, ok := decoder.GetOk("published"); ok {
		me.Published = opt.NewBool(value.(bool))
	}
	return nil
}

func (me *SharingInfo) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(me.Unknowns) > 0 {
		for k, v := range me.Unknowns {
			m[k] = v
		}
	}
	if me.Published != nil {
		rawMessage, err := json.Marshal(me.Published)
		if err != nil {
			return nil, err
		}
		m["published"] = rawMessage
	}

	if me.LinkShared != nil {
		rawMessage, err := json.Marshal(me.LinkShared)
		if err != nil {
			return nil, err
		}
		m["linkShared"] = rawMessage
	}
	return json.Marshal(m)
}

func (me *SharingInfo) UnmarshalJSON(data []byte) error {
	m := xjson.Properties{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := m.Unmarshal("linkShared", &me.LinkShared); err != nil {
		return err
	}
	if err := m.Unmarshal("published", &me.Published); err != nil {
		return err
	}

	if len(m) > 0 {
		me.Unknowns = m
	}
	return nil
}
