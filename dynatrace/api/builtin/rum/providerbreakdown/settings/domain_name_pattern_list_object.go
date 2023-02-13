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

package providerbreakdown

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DomainNamePatternListObjects []*DomainNamePatternListObject

func (me *DomainNamePatternListObjects) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"domain_name_pattern": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DomainNamePatternListObject).Schema()},
		},
	}
}

func (me DomainNamePatternListObjects) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("domain_name_pattern", me)
}

func (me *DomainNamePatternListObjects) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("domain_name_pattern", me)
}

type DomainNamePatternListObject struct {
	DomainNamePattern string `json:"domainNamePattern"` // Please type at least part of this content provider's URL. Asterisks (*) can be used as wildcard characters.
}

func (me *DomainNamePatternListObject) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pattern": {
			Type:        schema.TypeString,
			Description: "Please type at least part of this content provider's URL. Asterisks (*) can be used as wildcard characters.",
			Required:    true,
		},
	}
}

func (me *DomainNamePatternListObject) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"pattern": me.DomainNamePattern,
	})
}

func (me *DomainNamePatternListObject) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"pattern": &me.DomainNamePattern,
	})
}
