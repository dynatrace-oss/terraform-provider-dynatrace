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
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name            string `json:"-"`
	Value           string `json:"-"`
	Host            string `json:"-"`
	HostGroup       string `json:"-"`
	ManagementZone  string `json:"-"`
	ActiveGateGroup string `json:"-"`
}

var reg = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The fully qualified name of the extension, such as `com.dynatrace.extension.jmx-liberty-cp`. You can query for these names using the data source `dynatrace_hub_items`",
			ForceNew:    true,
			Required:    true,
		},
		"host": {
			Type:          schema.TypeString,
			Description:   "The ID of the host this monitoring configuration will be defined for",
			ForceNew:      true,
			Optional:      true,
			ConflictsWith: []string{"active_gate_group", "management_zone", "host_group"},
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, ok := i.(string)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
					return warnings, errs
				}
				if !strings.HasPrefix(v, "HOST-") {
					errs = append(errs, fmt.Errorf("value '%s' for %s is not the ID of a host ('HOST-#####')", v, k))
					return warnings, errs
				}
				return warnings, errs
			},
		},
		"host_group": {
			Type:          schema.TypeString,
			Description:   "The ID of the host group this monitoring configuration will be defined for",
			ForceNew:      true,
			Optional:      true,
			ConflictsWith: []string{"active_gate_group", "management_zone", "host"},
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, ok := i.(string)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
					return warnings, errs
				}
				if !strings.HasPrefix(v, "HOST_GROUP-") {
					errs = append(errs, fmt.Errorf("value '%s' for %s is not the ID of a host ('HOST_GROUP-#####')", v, k))
					return warnings, errs
				}
				return warnings, errs
			},
		},
		"management_zone": {
			Type:          schema.TypeString,
			Description:   "The name of the Management Zone this monitoring configuration will be defined for",
			ForceNew:      true,
			Optional:      true,
			ConflictsWith: []string{"active_gate_group", "host_group", "host"},
		},
		"active_gate_group": {
			Type:          schema.TypeString,
			Description:   "The name of the Active Gate Group this monitoring configuration will be defined for",
			ForceNew:      true,
			Optional:      true,
			ConflictsWith: []string{"management_zone", "host_group", "host"},
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
				version := ""
				if version, ok = m["version"].(string); !ok {
					errs = append(errs, fmt.Errorf("expected type of %s.version to be a string", k))
					return warnings, errs
				}
				if len(strings.TrimSpace(description)) == 0 {
					errs = append(errs, fmt.Errorf("%s.version must not be empty", k))
					return warnings, errs
				}
				if !reg.MatchString(version) {
					// this is perhaps too strict
					// errs = append(errs, fmt.Errorf("'%s' found in %s.version is not a valid version number. Expected format: `MAJOR.MINOR.REVISION` (e.g. `1.0.0`)", version, k))
					return warnings, errs
				}

				return warnings, errs
			},
		},
	}
}

func (me *Settings) CustomizeDiff(ctx context.Context, rd *schema.ResourceDiff, i any) error {
	getVersion := func(tv any) string {
		v := ""
		var err error
		var ok bool
		if v, ok = tv.(string); !ok {
			return ""
		}
		m := map[string]any{}
		if err = json.Unmarshal([]byte(v), &m); err != nil {
			return ""
		}
		if version, ok := m["version"]; ok {
			if sVersion, ok := version.(string); ok {
				return sVersion
			}
			return ""
		}
		return ""
	}
	prev, next := rd.GetChange("value")
	prevVersion := getVersion(prev)
	nextVersion := getVersion(next)
	if len(prevVersion) > 0 && prevVersion != nextVersion {
		rd.ForceNew("value")
	}
	return nil
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":              me.Name,
		"value":             me.Value,
		"host":              me.Host,
		"host_group":        me.HostGroup,
		"active_gate_group": me.ActiveGateGroup,
		"management_zone":   me.ManagementZone,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":              &me.Name,
		"value":             &me.Value,
		"host":              &me.Host,
		"host_group":        &me.HostGroup,
		"active_gate_group": &me.ActiveGateGroup,
		"management_zone":   &me.ManagementZone,
	})
}
