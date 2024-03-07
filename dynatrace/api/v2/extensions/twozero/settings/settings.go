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

package twozero

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name    string `json:"-"`
	Version string `json:"-"`
}

var reg = regexp.MustCompile(`^\d+\.\d+\.\d+$`)

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The fully qualified name of the extension, such as `com.dynatrace.extension.jmx-liberty-cp`",
			ForceNew:    true,
			Required:    true,
		},
		"version": {
			Type:        schema.TypeString,
			Description: "The fully qualified name of the extension, such as `com.dynatrace.extension.jmx-liberty-cp`",
			ValidateFunc: func(i any, k string) (warnings []string, errs []error) {
				v, ok := i.(string)
				if !ok {
					errs = append(errs, fmt.Errorf("expected type of %s to be string", k))
					return warnings, errs
				}
				if !reg.MatchString(v) {
					errs = append(errs, errors.New("expected format of %s to be `MAJOR.MINOR.REVISION` (e.g. `1.0.0`)"))
					return warnings, errs
				}
				return warnings, errs
			},
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":    me.Name,
		"version": me.Version,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":    &me.Name,
		"version": &me.Version,
	})
}
