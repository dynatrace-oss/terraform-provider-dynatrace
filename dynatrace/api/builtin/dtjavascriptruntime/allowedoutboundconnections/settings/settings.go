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

package allowedoutboundconnections

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AllowedOutboundConnections *AllowedHostsList `json:"allowedOutboundConnections"`
}

func (me *Settings) Name() string {
	return "limit_outbound_connections"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_outbound_connections": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(AllowedHostsList).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"allowed_outbound_connections": me.AllowedOutboundConnections,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"allowed_outbound_connections": &me.AllowedOutboundConnections,
	})
}
