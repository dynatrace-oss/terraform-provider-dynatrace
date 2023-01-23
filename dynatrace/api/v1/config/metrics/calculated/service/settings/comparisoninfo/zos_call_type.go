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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Zos Comparison for `ZOS_CALL_TYPE` attributes.
type ZOSCallType struct {
	BaseComparisonInfo
	Comparison ZOSCallTypeComparison `json:"comparison"`       // Operator of the comparision. You can reverse it by setting **negate** to `true`.
	Value      *ZOSCallTypeValue     `json:"value,omitempty"`  // The value to compare to.
	Values     []ZOSCallTypeValue    `json:"values,omitempty"` // The values to compare to.
}

func (me *ZOSCallType) GetType() Type {
	return Types.ZosCallType
}

func (me *ZOSCallType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The values to compare to. Possible values are `CTG`, `DPL`, `EXPLICIT_ADK`, `IMS_CONNECT`, `IMS_CONNECT_API`, `IMS_ITRA`, `IMS_MSC`, `IMS_PGM_SWITCH`, `IMS_SHARED_QUEUES`, `IMS_TRANS_EXEC`, `MQ`, `SOAP`, `START`, `TX` and `UNKNOWN`",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare to. Possible values are `CTG`, `DPL`, `EXPLICIT_ADK`, `IMS_CONNECT`, `IMS_CONNECT_API`, `IMS_ITRA`, `IMS_MSC`, `IMS_PGM_SWITCH`, `IMS_SHARED_QUEUES`, `IMS_TRANS_EXEC`, `MQ`, `SOAP`, `START`, `TX` and `UNKNOWN`",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Operator of the comparison. You can reverse it by setting `negate` to `true`. Possible values are `EQUALS`, `EQUALS_ANY_OF` and `EXISTS`",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *ZOSCallType) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"values":   me.Values,
		"value":    me.Value,
		"operator": me.Comparison,
		"unknowns": me.Unknowns,
	})
}

func (me *ZOSCallType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":   &me.Values,
		"value":    &me.Value,
		"operator": &me.Comparison,
		"unknowns": &me.Unknowns,
	})
}

func (me *ZOSCallType) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"type":       me.GetType(),
		"negate":     me.Negate,
		"values":     me.Values,
		"value":      me.Value,
		"comparison": me.Comparison,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *ZOSCallType) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	return properties.UnmarshalAll(map[string]any{
		"negate":     &me.Negate,
		"values":     &me.Values,
		"value":      &me.Value,
		"comparison": &me.Comparison,
	})
}

// ZOSCallTypeComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type ZOSCallTypeComparison string

// ZOSCallTypeComparisons offers the known enum values
var ZOSCallTypeComparisons = struct {
	Equals      ZOSCallTypeComparison
	EqualsAnyOf ZOSCallTypeComparison
	Exists      ZOSCallTypeComparison
}{
	"EQUALS",
	"EQUALS_ANY_OF",
	"EXISTS",
}

// ZOSCallTypeValue The value to compare to.
type ZOSCallTypeValue string

// ZOSCallTypeValues offers the known enum values
var ZOSCallTypeValues = struct {
	CTG             ZOSCallTypeValue
	Dpl             ZOSCallTypeValue
	ExplicitAdk     ZOSCallTypeValue
	IMSConnect      ZOSCallTypeValue
	IMSConnectApi   ZOSCallTypeValue
	IMSItra         ZOSCallTypeValue
	IMSMsc          ZOSCallTypeValue
	IMSPgmSwitch    ZOSCallTypeValue
	IMSSharedQueues ZOSCallTypeValue
	IMSTransExec    ZOSCallTypeValue
	Mq              ZOSCallTypeValue
	Soap            ZOSCallTypeValue
	Start           ZOSCallTypeValue
	Tx              ZOSCallTypeValue
	Unknown         ZOSCallTypeValue
}{
	"CTG",
	"DPL",
	"EXPLICIT_ADK",
	"IMS_CONNECT",
	"IMS_CONNECT_API",
	"IMS_ITRA",
	"IMS_MSC",
	"IMS_PGM_SWITCH",
	"IMS_SHARED_QUEUES",
	"IMS_TRANS_EXEC",
	"MQ",
	"SOAP",
	"START",
	"TX",
	"UNKNOWN",
}
