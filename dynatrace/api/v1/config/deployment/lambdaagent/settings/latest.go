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

package settings

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Latest struct {
	Java                string `json:"java"`                  // Latest version name of Java code module
	JavaWithCollector   string `json:"java_with_collector"`   // Latest version name of Java code module with log collector
	Python              string `json:"python"`                // Latest version name of Python code module
	PythonWithCollector string `json:"python_with_collector"` // Latest version name of Python code module with log collector
	NodeJS              string `json:"nodejs"`                // Latest version name of NodeJS code module
	NodeJSWithCollector string `json:"nodejs_with_collector"` // Latest version name of NodeJS code module with log collector
	Collector           string `json:"collector"`             // Latest version name of standalone log collector
}

func (me *Latest) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"java": {
			Type:        schema.TypeString,
			Description: "Latest version name of Java code module",
			Computed:    true,
			Optional:    true,
		},
		"java_with_collector": {
			Type:        schema.TypeString,
			Description: "Latest version name of Java code module with log collector",
			Computed:    true,
			Optional:    true,
		},
		"python": {
			Type:        schema.TypeString,
			Description: "Latest version name of Python code module",
			Computed:    true,
			Optional:    true,
		},
		"python_with_collector": {
			Type:        schema.TypeString,
			Description: "Latest version name of Python code module with log collector",
			Computed:    true,
			Optional:    true,
		},
		"nodejs": {
			Type:        schema.TypeString,
			Description: "Latest version name of NodeJS code module",
			Computed:    true,
			Optional:    true,
		},
		"nodejs_with_collector": {
			Type:        schema.TypeString,
			Description: "Latest version name of NodeJS code module with log collector",
			Computed:    true,
			Optional:    true,
		},
		"collector": {
			Type:        schema.TypeString,
			Description: "Latest version name of standalone log collector",
			Computed:    true,
			Optional:    true,
		},
	}
}

func (me *Latest) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"java":                  me.Java,
		"java_with_collector":   me.JavaWithCollector,
		"python":                me.Python,
		"python_with_collector": me.PythonWithCollector,
		"nodejs":                me.NodeJS,
		"nodejs_with_collector": me.NodeJSWithCollector,
		"collector":             me.Collector,
	})
}

func (me *Latest) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"java":                  &me.Java,
		"java_with_collector":   &me.JavaWithCollector,
		"python":                &me.Python,
		"python_with_collector": &me.PythonWithCollector,
		"nodejs":                &me.NodeJS,
		"nodejs_with_collector": &me.NodeJSWithCollector,
		"collector":             &me.Collector,
	})
}
