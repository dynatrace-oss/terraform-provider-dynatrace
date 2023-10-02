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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Host    *Host    `json:"host"`            // Hosts
	Network *Network `json:"network"`         // Network
	Scope   string   `json:"-" scope:"scope"` // The scope of this setting (HOST HOST_GROUP environment)
}

func (me *Settings) Name() string {
	return me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeList,
			Description: "Hosts",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Host).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"network": {
			Type:        schema.TypeList,
			Description: "Network",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(Network).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST HOST_GROUP environment)",
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"host":    me.Host,
		"network": me.Network,
		"scope":   me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"host":    &me.Host,
		"network": &me.Network,
		"scope":   &me.Scope,
	})
}
