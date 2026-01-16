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

package comparisoninfo

import (
	"encoding/json"
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ComparisonInfo Type-specific comparison for attributes. The actual set of fields depends on the `type` of the comparison.
// See the [Service metrics API - JSON models](https://dt-url.net/9803svb) help topic for example models of every notification type.
type ComparisonInfo interface {
	GetType() Type
	SetNegate(bool)
	IsNegate() bool
}

// BaseComparisonInfo Type-specific comparison for attributes. The actual set of fields depends on the `type` of the comparison.
// See the [Service metrics API - JSON models](https://dt-url.net/9803svb) help topic for example models of every notification type.
type BaseComparisonInfo struct {
	Negate   bool                       `json:"negate"` // Reverse the comparison **operator**. For example, it turns **equals** into **does not equal**.
	Type     Type                       `json:"type"`   // Defines the actual set of fields depending on the value. See one of the following objects:  * `STRING` -> StringComparisonInfo  * `NUMBER` -> NumberComparisonInfo  * `BOOLEAN` -> BooleanComparisonInfo  * `HTTP_METHOD` -> HttpMethodComparisonInfo  * `STRING_REQUEST_ATTRIBUTE` -> StringRequestAttributeComparisonInfo  * `NUMBER_REQUEST_ATTRIBUTE` -> NumberRequestAttributeComparisonInfo  * `ZOS_CALL_TYPE` -> ZosComparisonInfo  * `IIB_INPUT_NODE_TYPE` -> IIBInputNodeTypeComparisonInfo  * `ESB_INPUT_NODE_TYPE` -> ESBInputNodeTypeComparisonInfo  * `FAILED_STATE` -> FailedStateComparisonInfo  * `FLAW_STATE` -> FlawStateComparisonInfo  * `FAILURE_REASON` -> FailureReasonComparisonInfo  * `HTTP_STATUS_CLASS` -> HttpStatusClassComparisonInfo  * `TAG` -> TagComparisonInfo  * `FAST_STRING` -> FastStringComparisonInfo  * `SERVICE_TYPE` -> ServiceTypeComparisonInfo
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *BaseComparisonInfo) SetNegate(negate bool) {
	me.Negate = negate
}

func (me *BaseComparisonInfo) IsNegate() bool {
	return me.Negate
}

func (me *BaseComparisonInfo) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the actual set of fields depending on the value",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *BaseComparisonInfo) GetType() Type {
	return me.Type
}

func (me *BaseComparisonInfo) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	fmt.Printf("BaseComparisonInfo.Unknowns: %v", me.Unknowns)
	return properties.EncodeAll(map[string]any{
		"type":     me.Type,
		"unknowns": me.Unknowns,
	})
}

func (me *BaseComparisonInfo) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"type":     &me.Type,
		"unknowns": &me.Unknowns,
	})
}

func (me *BaseComparisonInfo) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"negate": me.Negate,
		"type":   me.Type,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *BaseComparisonInfo) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"negate": &me.Negate,
		"type":   &me.Type,
	}); err != nil {
		return err
	}
	me.Unknowns = properties
	return nil
}
