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

package automaticinjection

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type MonitoringCodeSource struct {
	CodeSource         string  `json:"codeSource"`                   // Real User Monitoring code source
	MonitoringCodePath *string `json:"monitoringCodePath,omitempty"` // Specify the source path for placement of your application's custom JavaScript library file. By default, this path is set to the root directory of your web server. A custom source path may be necessary if your server operates behind a firewall.
}

func (me *MonitoringCodeSource) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"code_source": {
			Type:        schema.TypeString,
			Description: "Real User Monitoring code source",
			Required:    true,
		},
		"monitoring_code_path": {
			Type:        schema.TypeString,
			Description: "Specify the source path for placement of your application's custom JavaScript library file. By default, this path is set to the root directory of your web server. A custom source path may be necessary if your server operates behind a firewall.",
			Optional:    true, // nullable & precondition
		},
	}
}

func (me *MonitoringCodeSource) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"code_source":          me.CodeSource,
		"monitoring_code_path": me.MonitoringCodePath,
	})
}

func (me *MonitoringCodeSource) HandlePreconditions() error {
	if (me.MonitoringCodePath != nil) && (string(me.CodeSource) != "OneAgent") {
		return fmt.Errorf("'monitoring_code_path' must not be specified if 'code_source' is set to '%v'", me.CodeSource)
	}
	return nil
}

func (me *MonitoringCodeSource) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"code_source":          &me.CodeSource,
		"monitoring_code_path": &me.MonitoringCodePath,
	})
}
