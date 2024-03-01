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

package script

import (
	http "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/http/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SyntheticMonitor HTTP synthetic monitor update. Some fields are inherited from base `SyntheticMonitorUpdate` model
type Settings struct {
	HttpId string       `json:"-"`
	Script *http.Script `json:"script,omitempty"`
}

func (me *Settings) Name() string {
	return me.HttpId
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"http_id": {
			Type:        schema.TypeString,
			Description: "The ID of the HTTP monitor",
			Required:    true,
		},
		"script": {
			Type:        schema.TypeList,
			Description: "The HTTP Script",
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: new(http.Script).Schema(),
			},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"http_id": me.HttpId,
		"script":  me.Script,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"http_id": &me.HttpId,
		"script":  &me.Script,
	})
}
