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

package ctxprop

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// PropagationSetting Context propagation enables you to connect PurePaths through OpenTelemetry/OpenTracing. Define rules to enable context propagation for certain spans within OneAgent
type PropagationSetting struct {
	PropagationRule *PropagationRule `json:"contextPropagationRule"`
}

func (me *PropagationSetting) Name() string {
	return me.PropagationRule.Name
}

func (me *PropagationSetting) Schema() map[string]*schema.Schema {
	return new(PropagationRule).Schema()
}

func (me *PropagationSetting) MarshalHCL(properties hcl.Properties) error {
	return me.PropagationRule.MarshalHCL(properties)
}

func (me *PropagationSetting) UnmarshalHCL(decoder hcl.Decoder) error {
	me.PropagationRule = new(PropagationRule)
	return me.PropagationRule.UnmarshalHCL(decoder)
}
