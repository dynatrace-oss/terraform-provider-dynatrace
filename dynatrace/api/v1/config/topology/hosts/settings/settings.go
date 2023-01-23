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

package hosts

import (
	tagapi "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/tag"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Hosts is a list of short representations of host
type Hosts []Host

func (me *Hosts) ToStubs() settings.Stubs {
	res := []*settings.Stub{}
	for _, setting := range *me {
		res = append(res, &settings.Stub{ID: setting.EntityId, Name: setting.DisplayName, Value: setting})
	}
	return res
}

type HostGroup struct {
	ID   string `json:"meId"`
	Name string `json:"name"`
}

// Host is a short representation of a host
type Host struct {
	EntityId    string       `json:"entityId"`    // The entity ID of the host
	DisplayName string       `json:"displayName"` // The name of the host as displayed in the UI
	Tags        []tagapi.Tag `json:"tags"`        // The list of entity tags
	HostGroup   *HostGroup   `json:"hostGroup"`
}

func (me *Host) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tags": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "Required tags of the host to find",
			MinItems:    1,
		},
	}
}

func (me *Host) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.DisplayName); err != nil {
		return err
	}
	tags := []string{}
	if len(me.Tags) > 0 {
		for _, tag := range me.Tags {
			if tag.Value == nil {
				tags = append(tags, tag.Key)
			} else {
				tags = append(tags, tag.Key+"="+*tag.Value)
			}
		}
	}
	if err := properties.Encode("tags", tags); err != nil {
		return err
	}

	return nil
}

func (me *Host) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("name", &me.DisplayName); err != nil {
		return err
	}
	var tagList []any
	if v, ok := decoder.GetOk("tags"); ok {
		sTags := v.(*schema.Set)
		tagList = sTags.List()
		tagapi.StringsToTags(tagList, &me.Tags)
	}
	return nil
}
