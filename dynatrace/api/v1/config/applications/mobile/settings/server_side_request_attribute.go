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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ServerSideRequestAttribute UserActionAndSessionProperty

func (me *ServerSideRequestAttribute) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The unique key of the mobile session or user action property",
			Required:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Description: "The ID of the request attribute",
			Required:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The data type of the property. Possible values are `DOUBLE`, `LONG` and `STRING`. The value MUST match the data type of the Request Attribute.",
			Required:    true,
		},
		"display_name": {
			Type:        schema.TypeString,
			Description: "The display name of the property",
			Optional:    true,
		},
		"store_as_user_action_property": {
			Type:        schema.TypeBool,
			Description: "If `true`, the property is stored as a user action property",
			Optional:    true,
		},
		"store_as_session_property": {
			Type:        schema.TypeBool,
			Description: "If `true`, the property is stored as a session property",
			Optional:    true,
		},
		"cleanup_rule": {
			Type:        schema.TypeString,
			Description: "The cleanup rule of the property. Defines how to extract the data you need from a string value. Specify the [regular expression](https://dt-url.net/k9e0iaq) for the data you need there",
			Optional:    true,
		},
		"aggregation": {
			Type:        schema.TypeString,
			Description: "The aggregation type of the property. It defines how multiple values of the property are aggregated. Possible values are `SUM`, `MIN`, `MAX`, `FIRST` and `LAST`",
			Optional:    true,
		},
	}
}

func (me *ServerSideRequestAttribute) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("key", &me.Key); err != nil {
		return err
	}
	if err := decoder.Decode("id", &me.ServerSideRequestAttribute); err != nil {
		return err
	}
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("display_name", &me.DisplayName); err != nil {
		return err
	}
	if err := decoder.Decode("store_as_user_action_property", &me.StoreAsUserActionProperty); err != nil {
		return err
	}
	if err := decoder.Decode("store_as_session_property", &me.StoreAsSessionProperty); err != nil {
		return err
	}
	if err := decoder.Decode("cleanup_rule", &me.CleanupRule); err != nil {
		return err
	}
	if err := decoder.Decode("aggregation", &me.Aggregation); err != nil {
		return err
	}
	me.Origin = Origins.ServerSideRequestAttribute
	me.Name = nil
	return nil
}

func (me *ServerSideRequestAttribute) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"key":                           me.Key,
		"id":                            me.ServerSideRequestAttribute,
		"type":                          me.Type,
		"display_name":                  me.DisplayName,
		"store_as_session_property":     me.StoreAsSessionProperty,
		"store_as_user_action_property": me.StoreAsUserActionProperty,
		"cleanup_rule":                  me.CleanupRule,
		"aggregation":                   me.Aggregation,
	})
}
