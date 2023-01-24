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

type UserActionAndSessionProperties []*UserActionAndSessionProperty

func (me *UserActionAndSessionProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"property": {
			Type:        schema.TypeList,
			Description: "User action and session properties settings",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(UserActionAndSessionProperty).Schema()},
		},
	}
}

func (me UserActionAndSessionProperties) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("property", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *UserActionAndSessionProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("property", me)
}

type UserActionAndSessionProperty struct {
	DisplayName                *string        `json:"displayName,omitempty"`                // The display name of the property
	Type                       PropertyType   `json:"type"`                                 // The data type of the property. Possible values are `DATE`, `DOUBLE`, `LONG`, `LONG_STRING` and `STRING`.
	Origin                     PropertyOrigin `json:"origin"`                               // The origin of the property. Possible values are `JAVASCRIPT_API`, `META_DATA` and `SERVER_SIDE_REQUEST_ATTRIBUTE`.
	Aggregation                *Aggregation   `json:"aggregation,omitempty"`                // The aggregation type of the property. \n\n  It defines how multiple values of the property are aggregated. Possible values are `AVERAGE`, `FIRST`, `LAST`, `MAXIMUM`, `MINIMUM` and `SUM`.
	StoreAsUserActionProperty  bool           `json:"storeAsUserActionProperty"`            // If `true`, the property is stored as a user action property
	StoreAsSessionProperty     bool           `json:"storeAsSessionProperty"`               // If `true`, the property is stored as a session property
	CleanupRule                *string        `json:"cleanupRule,omitempty"`                // The cleanup rule of the property. \n\nDefines how to extract the data you need from a string value. Specify the [regular expression](https://dt-url.net/k9e0iaq) for the data you need there
	ServerSideRequestAttribute *string        `json:"serverSideRequestAttribute,omitempty"` // The ID of the request attribute. \n\nOnly applicable when the **origin** is set to `SERVER_SIDE_REQUEST_ATTRIBUTE`
	UniqueID                   int32          `json:"uniqueId"`                             // Unique id among all userTags and properties of this application
	Key                        string         `json:"key"`                                  // Key of the property
	MetaDataID                 *int32         `json:"metadataId,omitempty"`                 // If the origin is `META_DATA`, metaData id of the property
	IgnoreCase                 bool           `json:"ignoreCase,omitempty"`                 // If `true`, the value of this property will always be stored in lower case. Defaults to `false`.
	LongStringLength           *int32         `json:"longStringLength,omitempty"`           // If the `type` is `LONG_STRING`, the max length for this property. Must be a multiple of `100`. Defaults to `200`. Maximum is `1000`.
}

func (me *UserActionAndSessionProperty) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Type:        schema.TypeString,
			Description: "The display name of the property",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The data type of the property. Possible values are `DATE`, `DOUBLE`, `LONG`, `LONG_STRING` and `STRING`.",
			Required:    true,
		},
		"origin": {
			Type:        schema.TypeString,
			Description: "The origin of the property. Possible values are `JAVASCRIPT_API`, `META_DATA` and `SERVER_SIDE_REQUEST_ATTRIBUTE`.",
			Required:    true,
		},
		"aggregation": {
			Type:        schema.TypeString,
			Description: "The aggregation type of the property. \n\n  It defines how multiple values of the property are aggregated. Possible values are `AVERAGE`, `FIRST`, `LAST`, `MAXIMUM`, `MINIMUM` and `SUM`.",
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
			Description: "The cleanup rule of the property. \n\nDefines how to extract the data you need from a string value. Specify the [regular expression](https://dt-url.net/k9e0iaq) for the data you need there",
			Optional:    true,
		},
		"server_side_request_attribute": {
			Type:        schema.TypeString,
			Description: "The ID of the request attribute. \n\nOnly applicable when the **origin** is set to `SERVER_SIDE_REQUEST_ATTRIBUTE`",
			Optional:    true,
		},
		"id": {
			Type:        schema.TypeInt,
			Description: "Unique id among all userTags and properties of this application",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Key of the property",
			Required:    true,
		},
		"metadata_id": {
			Type:        schema.TypeInt,
			Description: "If the origin is `META_DATA`, metaData id of the property",
			Optional:    true,
		},
		"ignore_case": {
			Type:        schema.TypeBool,
			Description: "If `true`, the value of this property will always be stored in lower case. Defaults to `false`.",
			Optional:    true,
		},
		"long_string_length": {
			Type:        schema.TypeInt,
			Description: "If the `type` is `LONG_STRING`, the max length for this property. Must be a multiple of `100`. Defaults to `200`. Maximum is `1000`.",
			Optional:    true,
		},
	}
}

func (me *UserActionAndSessionProperty) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"display_name":                  me.DisplayName,
		"type":                          me.Type,
		"origin":                        me.Origin,
		"aggregation":                   me.Aggregation,
		"store_as_user_action_property": me.StoreAsUserActionProperty,
		"store_as_session_property":     me.StoreAsSessionProperty,
		"cleanup_rule":                  me.CleanupRule,
		"server_side_request_attribute": me.ServerSideRequestAttribute,
		"id":                            me.UniqueID,
		"key":                           me.Key,
		"metadata_id":                   me.MetaDataID,
		"ignore_case":                   me.IgnoreCase,
		"long_string_length":            me.LongStringLength,
	})
}

func (me *UserActionAndSessionProperty) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"display_name":                  &me.DisplayName,
		"type":                          &me.Type,
		"origin":                        &me.Origin,
		"aggregation":                   &me.Aggregation,
		"store_as_user_action_property": &me.StoreAsUserActionProperty,
		"store_as_session_property":     &me.StoreAsSessionProperty,
		"cleanup_rule":                  &me.CleanupRule,
		"server_side_request_attribute": &me.ServerSideRequestAttribute,
		"id":                            &me.UniqueID,
		"key":                           &me.Key,
		"metadata_id":                   &me.MetaDataID,
		"ignore_case":                   &me.IgnoreCase,
		"long_string_length":            &me.LongStringLength,
	})
}
