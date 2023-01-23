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

package event

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Target struct {
	Window   *string  `json:"window,omitempty"`   // The tab of the target
	Locators Locators `json:"locators,omitempty"` // The list of locators identifying the desired element
}

func (me *Target) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"window": {
			Type:        schema.TypeString,
			Description: "The tab of the target",
			Optional:    true,
		},
		"locators": {
			Type:        schema.TypeList,
			Description: "The list of locators identifying the desired element",
			Optional:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Locators).Schema()},
		},
	}
}

func (me *Target) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("window", me.Window); err != nil {
		return err
	}
	if err := properties.Encode("locators", me.Locators); err != nil {
		return err
	}
	return nil
}

func (me *Target) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("window", &me.Window); err != nil {
		return err
	}
	if err := decoder.Decode("locators", &me.Locators); err != nil {
		return err
	}
	return nil
}
