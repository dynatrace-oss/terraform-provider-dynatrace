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

package browser

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Restrictions []*Restriction

func (me *Restrictions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"restriction": {
			Type:        schema.TypeList,
			Description: "Browser exclusion rules for the browsers that are to be excluded",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(Restriction).Schema()},
		},
	}
}

func (me Restrictions) MarshalHCL(properties hcl.Properties) error {
	if len(me) > 0 {
		if err := properties.EncodeSlice("restriction", me); err != nil {
			return err
		}
	}
	return nil
}

func (me *Restrictions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("restriction", me)
}

// Restriction Browser exclusion rules for the browsers that are to be excluded
type Restriction struct {
	BrowserVersion *string     `json:"browserVersion,omitempty"` // The version of the browser that is used
	BrowserType    Type        `json:"browserType"`              // The type of the browser that is used. Possible values are `ANDROID_WEBKIT`, `BOTS_SPIDERS`, `CHROME`, `EDGE`, `FIREFOX`, `INTERNET_EXPLORER, `OPERA` and `SAFARI`
	Platform       *Platform   `json:"platform,omitempty"`       // The platform on which the browser is being used. Possible values are `ALL`, `DESKTOP` and `MOBILE`
	Comparator     *Comparator `json:"comparator,omitempty"`     // No documentation available
}

func (me *Restriction) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"browser_version": {
			Type:        schema.TypeString,
			Description: "The version of the browser that is used",
			Optional:    true,
		},
		"browser_type": {
			Type:        schema.TypeString,
			Description: "The type of the browser that is used. Possible values are `ANDROID_WEBKIT`, `BOTS_SPIDERS`, `CHROME`, `EDGE`, `FIREFOX`, `INTERNET_EXPLORER, `OPERA` and `SAFARI`",
			Required:    true,
		},
		"platform": {
			Type:        schema.TypeString,
			Description: "The platform on which the browser is being used. Possible values are `ALL`, `DESKTOP` and `MOBILE`",
			Optional:    true,
		},
		"comparator": {
			Type:        schema.TypeString,
			Description: "No documentation available. Possible values are `EQUALS`, `GREATER_THAN_OR_EQUAL` and `LOWER_THAN_OR_EQUAL`.",
			Optional:    true,
		},
	}
}

func (me *Restriction) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"browser_version": me.BrowserVersion,
		"browser_type":    me.BrowserType,
		"platform":        me.Platform,
		"comparator":      me.Comparator,
	})
}

func (me *Restriction) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"browser_version": &me.BrowserVersion,
		"browser_type":    &me.BrowserType,
		"platform":        &me.Platform,
		"comparator":      &me.Comparator,
	})
}
