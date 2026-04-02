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

package extension_config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name  string         `json:"-"`
	Scope string         `json:"scope"`
	Value map[string]any `json:"value"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The fully qualified name of the extension, such as `com.dynatrace.extension.jmx-liberty-cp`.",
			ForceNew:    true,
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope this monitoring configuration will be defined for",
			Required:    true,
			ForceNew:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The JSON encoded value for this monitoring configuration",
			Required:    true,
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, ok := i.(string)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
					return warnings, errs
				}
				m := map[string]any{}
				if err := json.Unmarshal([]byte(v), &m); err != nil {
					errs = append(errs, fmt.Errorf("%s is not valid JSON", k))
					return warnings, errs
				}
				if _, found := m["description"]; !found {
					errs = append(errs, fmt.Errorf("%s doesn't contain a description", k))
					return warnings, errs
				}
				description := ""
				if description, ok = m["description"].(string); !ok {
					errs = append(errs, fmt.Errorf("expected type of %s.description to be a string", k))
					return warnings, errs
				}
				if len(strings.TrimSpace(description)) == 0 {
					errs = append(errs, fmt.Errorf("%s.description must not be empty", k))
					return warnings, errs
				}

				if _, found := m["version"]; !found {
					errs = append(errs, fmt.Errorf("%s doesn't contain a version number", k))
					return warnings, errs
				}
				var version string
				if version, ok = m["version"].(string); !ok {
					errs = append(errs, fmt.Errorf("expected type of %s.version to be a string", k))
					return warnings, errs
				}
				if len(strings.TrimSpace(version)) == 0 {
					errs = append(errs, fmt.Errorf("%s.version must not be empty", k))
					return warnings, errs
				}

				return warnings, errs
			},
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	valueJSON, err := json.Marshal(me.Value)
	if err != nil {
		return err
	}

	return properties.EncodeAll(map[string]any{
		"name":  me.Name,
		"value": string(valueJSON),
		"scope": me.Scope,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	var valueString string
	err := decoder.DecodeAll(map[string]any{
		"name":  &me.Name,
		"value": &valueString,
		"scope": &me.Scope,
	})
	if err != nil {
		return err
	}

	if valueString != "" {
		return json.Unmarshal([]byte(valueString), &me.Value)
	}
	return nil
}
