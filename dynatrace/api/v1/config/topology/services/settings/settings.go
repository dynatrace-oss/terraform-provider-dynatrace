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

package service

import (
	tagapi "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/topology/tag"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Services is a list of short representations of services
type Services []*Settings

func (me *Services) ToStubs() settings.Stubs {
	res := []*settings.Stub{}
	for _, setting := range *me {
		res = append(res, &settings.Stub{ID: setting.EntityId, Name: setting.DisplayName, Value: setting})
	}
	return res
}

// Service is a short representation of a service
type Settings struct {
	EntityId    string       `json:"entityId"`    // The entity ID of the service
	DisplayName string       `json:"displayName"` // The name of the service as displayed in the UI
	Tags        []tagapi.Tag `json:"tags"`        // The list of entity tags
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tags": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Description: "Required tags of the service to find",
			MinItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
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

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
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
