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

package endpointmetrics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	EnableEndpointMetrics bool    `json:"enableEndpointMetrics"` // Should metrics be written for endpoints? Please be aware that this setting has billing implications. Check out this [documentation](https://dt-url.net/td23cgh) for further details.
	ServiceID             *string `json:"-" scope:"serviceId"`   // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_endpoint_metrics": {
			Type:        schema.TypeBool,
			Description: "Should metrics be written for endpoints? Please be aware that this setting has billing implications. Check out this [documentation](https://dt-url.net/td23cgh) for further details.",
			Required:    true,
		},
		"service_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) Name() string {
	return *me.ServiceID
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enable_endpoint_metrics": me.EnableEndpointMetrics,
		"service_id":              me.ServiceID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enable_endpoint_metrics": &me.EnableEndpointMetrics,
		"service_id":              &me.ServiceID,
	})
}
