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

package web

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UserTags []*UserTag

func (me *UserTags) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tag": {
			Type:        schema.TypeList,
			Description: "User tag settings",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(UserTag).Schema()},
		},
	}
}

func (me UserTags) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("tag", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *UserTags) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("tag", me)
}

type UserTag struct {
	UniqueID                   int32   `json:"uniqueId"`                             // A unique ID among all userTags and properties of this application. Minimum value is 1.
	MetaDataID                 *int32  `json:"metadataId,omitempty"`                 // If it's of type metaData, metaData id of the userTag
	CleanUpRule                *string `json:"cleanupRule,omitempty"`                // Cleanup rule expression of the userTag
	ServerSideRequestAttribute *string `json:"serverSideRequestAttribute,omitempty"` // The ID of the RrequestAttribute for the userTag
	IgnoreCase                 bool    `json:"ignoreCase,omitempty"`                 // If `true`, the value of this tag will always be stored in lower case. Defaults to `false`.
}

func (me *UserTag) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Description: "A unique ID among all userTags and properties of this application. Minimum value is 1.",
			Required:    true,
		},
		"metadata_id": {
			Type:        schema.TypeInt,
			Description: "If it's of type metaData, metaData id of the userTag",
			Optional:    true,
		},
		"cleanup_rule": {
			Type:        schema.TypeString,
			Description: "Cleanup rule expression of the userTag",
			Optional:    true,
		},
		"server_side_request_attribute": {
			Type:        schema.TypeString,
			Description: "The ID of the RrequestAttribute for the userTag",
			Optional:    true,
		},
		"ignore_case": {
			Type:        schema.TypeBool,
			Description: "If `true`, the value of this tag will always be stored in lower case. Defaults to `false`.",
			Optional:    true,
		},
	}
}

func (me *UserTag) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"id":                            me.UniqueID,
		"metadata_id":                   me.MetaDataID,
		"cleanup_rule":                  me.CleanUpRule,
		"server_side_request_attribute": me.ServerSideRequestAttribute,
		"ignore_case":                   me.IgnoreCase,
	})
}

func (me *UserTag) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"id":                            &me.UniqueID,
		"metadata_id":                   &me.MetaDataID,
		"cleanup_rule":                  &me.CleanUpRule,
		"server_side_request_attribute": &me.ServerSideRequestAttribute,
		"ignore_case":                   &me.IgnoreCase,
	})
}
