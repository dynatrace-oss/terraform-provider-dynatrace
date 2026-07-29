/**
* @license
* Copyright 2026 Dynatrace LLC
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

package settings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CloudWatchLogsConfig struct {
	Enabled bool
	Regions []string
}

func (me *CloudWatchLogsConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "Whether to ingest CloudWatch logs in addition to metrics.",
			Optional:    true,
			Default:     true,
		},
		"regions": {
			Type:        schema.TypeList,
			Description: "Regions to ingest CloudWatch logs from. When empty, the configuration's `regions` are used.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *CloudWatchLogsConfig) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled": me.Enabled,
		"regions": me.Regions,
	})
}

func (me *CloudWatchLogsConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled": &me.Enabled,
		"regions": &me.Regions,
	})
}
