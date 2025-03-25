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

package customservices_order

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Technology struct {
	Name     string
	RESTName string
}

var Technologies = map[string]Technology{
	"dotnet": {".Net", "dotNet"},
	"golang": {"Golang", "go"},
	"java":   {"Java", "java"},
	"nodejs": {"NodeJS", "nodeJS"},
	"php":    {"PHP", "php"},
}

type Settings struct {
	Name string `json:"-"`
	IDs  map[string][]string
}

func (me *Settings) Schema() map[string]*schema.Schema {
	result := map[string]*schema.Schema{}
	for techID, techName := range Technologies {
		result[techID] = &schema.Schema{
			Type:        schema.TypeList,
			Optional:    true,
			Description: fmt.Sprintf("The IDs of Custom Services for %s in the desired order", techName.Name),
			Elem:        &schema.Schema{Type: schema.TypeString},
		}
	}
	return result
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	if len(me.IDs) == 0 {
		return nil
	}
	for technology := range Technologies {
		if ids, found := me.IDs[technology]; found {
			if err := properties.EncodeSlice(technology, ids); err != nil {
				return err
			}
		}
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	me.IDs = map[string][]string{}
	for technology := range Technologies {
		ids := []string{}
		if err := decoder.DecodeSlice(technology, &ids); err != nil {
			return err
		}
		me.IDs[technology] = ids
	}
	return nil
}
