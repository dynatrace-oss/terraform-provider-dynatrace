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

package alerting

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SeverityRules []*SeverityRule

func (me *SeverityRules) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rule": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "A conditions for the metric usage",
			Elem:        &schema.Resource{Schema: new(SeverityRule).Schema()},
		},
	}
}

func (me SeverityRules) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("rule", me)
}

func (me *SeverityRules) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("rule", me)
}

type SeverityRule struct {
	SeverityLevel        SeverityLevel        `json:"severityLevel"`        // Problem severity level
	DelayInMinutes       int32                `json:"delayInMinutes"`       // Send a notification if a problem remains open longer than X minutes. Must be between 0 and 10000.
	TagFilterIncludeMode TagFilterIncludeMode `json:"tagFilterIncludeMode"` // Possible values are `NONE`, `INCLUDE_ANY` and `INCLUDE_ALL`
	Tags                 []string             `json:"tagFilter"`            // SET / no documentation available
}

func (me *SeverityRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"severity_level": {
			Type:        schema.TypeString,
			Description: "The severity level to trigger the alert. Possible values are `AVAILABILITY`,	`CUSTOM_ALERT`,	`ERRORS`,`MONITORING_UNAVAILABLE`,`PERFORMANCE` and `RESOURCE_CONTENTION`.",
			Required:    true,
		},
		"delay_in_minutes": {
			Type:        schema.TypeInt,
			Description: "Send a notification if a problem remains open longer than *X* minutes",
			Required:    true,
		},
		"include_mode": {
			Type:        schema.TypeString,
			Description: "The filtering mode:  * `INCLUDE_ANY`: The rule applies to monitored entities that have at least one of the specified tags. You can specify up to 100 tags.  * `INCLUDE_ALL`: The rule applies to monitored entities that have **all** of the specified tags. You can specify up to 10 tags.  * `NONE`: The rule applies to all monitored entities",
			Required:    true,
		},
		"tags": {
			Type:        schema.TypeSet,
			Description: "A set of tags you want to filter by. You can also specify a tag value alongside the tag name using the syntax `name:value`.",
			MinItems:    1,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (me *SeverityRule) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("delay_in_minutes", int(me.DelayInMinutes)); err != nil {
		return err
	}
	if err := properties.Encode("severity_level", string(me.SeverityLevel)); err != nil {
		return err
	}
	if err := properties.Encode("include_mode", string(me.TagFilterIncludeMode)); err != nil {
		return err
	}
	if err := properties.Encode("tags", me.Tags); err != nil {
		return err
	}
	if tags, ok := properties["tags"]; ok && tags == nil {
		properties["tags"] = []string{}
	}
	return nil
}

func (me *SeverityRule) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("severity_level"); ok {
		me.SeverityLevel = SeverityLevel(value.(string))
	}
	if value, ok := decoder.GetOk("include_mode"); ok {
		me.TagFilterIncludeMode = TagFilterIncludeMode(value.(string))
	}
	if value, ok := decoder.GetOk("delay_in_minutes"); ok {
		me.DelayInMinutes = int32(value.(int))
	}
	if err := decoder.Decode("tags", &me.Tags); err != nil {
		return err
	}
	return nil
}
