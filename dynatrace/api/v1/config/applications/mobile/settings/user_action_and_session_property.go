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
	"sort"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UserActionAndSessionProperty represents configuration of the mobile session or user action property
type UserActionAndSessionProperty struct {
	Key                        string       `json:"key"`                                  // The unique key of the mobile session or user action property
	DisplayName                *string      `json:"displayName,omitempty"`                // The display name of the property
	Type                       PropertyType `json:"type"`                                 // The data type of the property
	Origin                     Origin       `json:"origin"`                               // The origin of the property
	Aggregation                *Aggregation `json:"aggregation,omitempty"`                // The aggregation type of the property. It defines how multiple values of the property are aggregated
	StoreAsUserActionProperty  bool         `json:"storeAsUserActionProperty"`            // If `true`, the property is stored as a user action property
	StoreAsSessionProperty     bool         `json:"storeAsSessionProperty"`               // If `true`, the property is stored as a session property
	CleanupRule                *string      `json:"cleanupRule,omitempty"`                // The cleanup rule of the property. Defines how to extract the data you need from a string value. Specify the [regular expression](https://dt-url.net/k9e0iaq) for the data you need there
	ServerSideRequestAttribute *string      `json:"serverSideRequestAttribute,omitempty"` // The ID of the request attribute. Only applicable when the **origin** is set to `SERVER_SIDE_REQUEST_ATTRIBUTE`
	Name                       *string      `json:"name,omitempty"`                       // The name of the reported value. Only applicable when the **origin** is set to `API`
}

type UserActionAndSessionProperties []*UserActionAndSessionProperty

func (me UserActionAndSessionProperties) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"request_attribute": {
			Type:        schema.TypeList,
			Description: "A User Action / Session Property based on a Server Side Request Attribute",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(ServerSideRequestAttribute).Schema()},
		},
		"api_value": {
			Type:        schema.TypeList,
			Description: "A User Action / Session Property based on a value reported by the API",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(APIValue).Schema()},
		},
	}
}

func (me UserActionAndSessionProperties) MarshalHCL(properties hcl.Properties) error {
	var err error
	api_values := []*APIValue{}
	for _, property := range me {
		if property.Origin == Origins.API {
			api_values = append(api_values, &APIValue{
				Key:                        property.Key,
				DisplayName:                property.DisplayName,
				Type:                       property.Type,
				Origin:                     Origins.API,
				Aggregation:                property.Aggregation,
				StoreAsUserActionProperty:  property.StoreAsUserActionProperty,
				StoreAsSessionProperty:     property.StoreAsSessionProperty,
				CleanupRule:                property.CleanupRule,
				ServerSideRequestAttribute: nil,
				Name:                       property.Name,
			})
		}
	}
	sort.SliceStable(api_values, func(i, j int) bool {
		return api_values[i].Key < api_values[j].Key
	})

	if len(api_values) > 0 {
		if err = properties.EncodeSlice("api_value", api_values); err != nil {
			return err
		}
	}

	request_attributes := []*ServerSideRequestAttribute{}
	for _, property := range me {
		if property.Origin == Origins.ServerSideRequestAttribute {
			request_attributes = append(request_attributes, &ServerSideRequestAttribute{
				Key:                        property.Key,
				DisplayName:                property.DisplayName,
				Type:                       property.Type,
				Origin:                     Origins.ServerSideRequestAttribute,
				Aggregation:                property.Aggregation,
				StoreAsUserActionProperty:  property.StoreAsUserActionProperty,
				StoreAsSessionProperty:     property.StoreAsSessionProperty,
				CleanupRule:                property.CleanupRule,
				ServerSideRequestAttribute: property.ServerSideRequestAttribute,
				Name:                       nil,
			})
		}
	}

	sort.SliceStable(request_attributes, func(i, j int) bool {
		return request_attributes[i].Key < request_attributes[j].Key
	})
	if len(request_attributes) > 0 {
		if err = properties.EncodeSlice("request_attribute", request_attributes); err != nil {
			return err
		}
	}

	return nil
}

func (me *UserActionAndSessionProperties) UnmarshalHCL(decoder hcl.Decoder) error {
	if result, ok := decoder.GetOk("api_value.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(APIValue)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "api_value", idx)); err != nil {
				return err
			}
			*me = append(*me, &UserActionAndSessionProperty{
				Key:                        entry.Key,
				DisplayName:                entry.DisplayName,
				Type:                       entry.Type,
				Origin:                     Origins.API,
				Aggregation:                entry.Aggregation,
				StoreAsUserActionProperty:  entry.StoreAsUserActionProperty,
				StoreAsSessionProperty:     entry.StoreAsSessionProperty,
				CleanupRule:                entry.CleanupRule,
				ServerSideRequestAttribute: nil,
				Name:                       entry.Name,
			})
		}
	}
	if result, ok := decoder.GetOk("request_attribute.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(ServerSideRequestAttribute)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "request_attribute", idx)); err != nil {
				return err
			}
			*me = append(*me, &UserActionAndSessionProperty{
				Key:                        entry.Key,
				DisplayName:                entry.DisplayName,
				Type:                       entry.Type,
				Origin:                     Origins.ServerSideRequestAttribute,
				Aggregation:                entry.Aggregation,
				StoreAsUserActionProperty:  entry.StoreAsUserActionProperty,
				StoreAsSessionProperty:     entry.StoreAsSessionProperty,
				CleanupRule:                entry.CleanupRule,
				ServerSideRequestAttribute: entry.ServerSideRequestAttribute,
				Name:                       nil,
			})
		}
	}
	return nil
}
