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

package contextpropagation

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PropagationSetting Context propagation enables you to connect PurePaths through OpenTelemetry/OpenTracing. Define rules to enable context propagation for certain spans within OneAgent
type PropagationSetting struct {
	PropagationRule *PropagationRule `json:"contextPropagationRule"`
	InsertAfter     string           `json:"-"`
}

func (me *PropagationSetting) Name() string {
	return me.PropagationRule.Name
}

func (me *PropagationSetting) Schema() map[string]*schema.Schema {
	var sch = new(PropagationRule).Schema()
	sch["insert_after"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
		Optional:    true,
		Computed:    true,
	}
	return sch
}

func (me *PropagationSetting) MarshalHCL(properties hcl.Properties) error {
	if err := me.PropagationRule.MarshalHCL(properties); err != nil {
		return err
	}
	if len(me.InsertAfter) > 0 {
		properties["insert_after"] = me.InsertAfter
	}
	return nil
}

func (me *PropagationSetting) UnmarshalHCL(decoder hcl.Decoder) error {
	me.PropagationRule = new(PropagationRule)
	if err := me.PropagationRule.UnmarshalHCL(decoder); err != nil {
		return err
	}
	if v, ok := decoder.GetOk("insert_after"); ok {
		if sv, ok := v.(string); ok {
			me.InsertAfter = sv
		}
	}
	return nil
}
