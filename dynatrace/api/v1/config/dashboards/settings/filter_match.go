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

package dashboards

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FilterMatch struct {
	Key    string
	Values []string
}

func (me *FilterMatch) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Description: "The entity type (e.g. HOST, SERVICE, ...)",
			Required:    true,
		},
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "the tiles this Dashboard consist of",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *FilterMatch) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("key", me.Key); err != nil {
		return err
	}
	if err := properties.Encode("values", me.Values); err != nil {
		return err
	}

	return nil
}

func (me *FilterMatch) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("key"); ok {
		me.Key = value.(string)
	}
	if err := decoder.Decode("values", &me.Values); err != nil {
		return err
	}
	return nil
}
