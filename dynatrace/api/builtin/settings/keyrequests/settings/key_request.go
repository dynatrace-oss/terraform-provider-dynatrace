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

package keyrequests

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// KeyRequest has no documentation
type KeyRequest struct {
	Names         []string          `json:"keyRequestNames,omitempty"`
	ServiceID     string            `json:"-" scope:"serviceId"`
	KeyRequestIDs map[string]string `json:"-"`
}

func (me *KeyRequest) Name() string {
	return "Key Requests for " + me.ServiceID
}

func (me *KeyRequest) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"service": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ID of Dynatrace Service, eg. SERVICE-123ABC45678EFGH",
			ForceNew:    true,
		},
		"names": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The names of the key requests",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"key_request_ids": {
			Type:        schema.TypeMap,
			Description: "The ids of the key requests",
			Computed:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *KeyRequest) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"names":           me.Names,
		"service":         me.ServiceID,
		"key_request_ids": me.KeyRequestIDs,
	})
}

func (me *KeyRequest) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"names":           &me.Names,
		"service":         &me.ServiceID,
		"key_request_ids": &me.KeyRequestIDs,
	})
}
