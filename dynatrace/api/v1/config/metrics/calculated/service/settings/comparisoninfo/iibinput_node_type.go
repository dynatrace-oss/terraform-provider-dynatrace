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

// IIBInputNodeType Comparison for `IIB_INPUT_NODE_TYPE` attributes.
type IIBInputNodeType struct {
	BaseComparisonInfo
	Value      *IIBInputNodeTypeValue     `json:"value,omitempty"`  // The value to compare to.
	Values     []IIBInputNodeTypeValue    `json:"values,omitempty"` // The values to compare to.
	Comparison IIBInputNodeTypeComparison `json:"comparison"`       // Operator of the comparision. You can reverse it by setting **negate** to `true`.
}

func (me *IIBInputNodeType) GetType() Type {
	return Types.IIBInputNodeType
}

func (me *IIBInputNodeType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The values to compare to. Possible values are `CALLABLE_FLOW_ASYNC_RESPONSE_NODE`, `CALLABLE_FLOW_INPUT_NODE`, `DATABASE_INPUT_NODE`, `DOTNET_INPUT_NODE`, `EMAIL_INPUT_NODE`, `EVENT_INPUT`, `EVENT_INPUT_NODE`, `FILE_INPUT_NODE`, `FTE_INPUT_NODE`, `HTTP_ASYNC_RESPONSE`, `JD_EDWARDS_INPUT_NODE`, `JMS_CLIENT_INPUT_NODE`, `LABEL_NODE`, `MQ_INPUT_NODE`, `PEOPLE_SOFT_INPUT_NODE`, `REST_ASYNC_RESPONSE`, `REST_REQUEST`, `SAP_INPUT_NODE`, `SCA_ASYNC_RESPONSE_NODE`, `SCA_INPUT_NODE`, `SIEBEL_INPUT_NODE`, `SOAP_INPUT_NODE`, `TCPIP_CLIENT_INPUT_NODE`, `TCPIP_CLIENT_REQUEST_NODE`, `TCPIP_SERVER_INPUT_NODE`, `TCPIP_SERVER_REQUEST_NODE`, `TIMEOUT_NOTIFICATION_NODE` and `WS_INPUT_NODE`",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare to. Possible values are `CALLABLE_FLOW_ASYNC_RESPONSE_NODE`, `CALLABLE_FLOW_INPUT_NODE`, `DATABASE_INPUT_NODE`, `DOTNET_INPUT_NODE`, `EMAIL_INPUT_NODE`, `EVENT_INPUT`, `EVENT_INPUT_NODE`, `FILE_INPUT_NODE`, `FTE_INPUT_NODE`, `HTTP_ASYNC_RESPONSE`, `JD_EDWARDS_INPUT_NODE`, `JMS_CLIENT_INPUT_NODE`, `LABEL_NODE`, `MQ_INPUT_NODE`, `PEOPLE_SOFT_INPUT_NODE`, `REST_ASYNC_RESPONSE`, `REST_REQUEST`, `SAP_INPUT_NODE`, `SCA_ASYNC_RESPONSE_NODE`, `SCA_INPUT_NODE`, `SIEBEL_INPUT_NODE`, `SOAP_INPUT_NODE`, `TCPIP_CLIENT_INPUT_NODE`, `TCPIP_CLIENT_REQUEST_NODE`, `TCPIP_SERVER_INPUT_NODE`, `TCPIP_SERVER_REQUEST_NODE`, `TIMEOUT_NOTIFICATION_NODE` and `WS_INPUT_NODE`",
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

func (me *IIBInputNodeType) MarshalHCL(properties hcl.Properties) error {
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

func (me *IIBInputNodeType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":   &me.Values,
		"value":    &me.Value,
		"operator": &me.Comparison,
		"unknowns": &me.Unknowns,
	})
}

func (me *IIBInputNodeType) MarshalJSON() ([]byte, error) {
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

func (me *IIBInputNodeType) UnmarshalJSON(data []byte) error {
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

// IIBInputNodeTypeComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type IIBInputNodeTypeComparison string

// IIBInputNodeTypeComparisons offers the known enum values
var IIBInputNodeTypeComparisons = struct {
	Equals      IIBInputNodeTypeComparison
	EqualsAnyOf IIBInputNodeTypeComparison
	Exists      IIBInputNodeTypeComparison
}{
	"EQUALS",
	"EQUALS_ANY_OF",
	"EXISTS",
}

// IIBInputNodeTypeValue The value to compare to.
type IIBInputNodeTypeValue string

// IIBInputNodeTypeValues offers the known enum values
var IIBInputNodeTypeValues = struct {
	CallableFlowAsyncResponseNode IIBInputNodeTypeValue
	CallableFlowInputNode         IIBInputNodeTypeValue
	DatabaseInputNode             IIBInputNodeTypeValue
	DotNetInputNode               IIBInputNodeTypeValue
	EmailInputNode                IIBInputNodeTypeValue
	EventInput                    IIBInputNodeTypeValue
	EventInputNode                IIBInputNodeTypeValue
	FileInputNode                 IIBInputNodeTypeValue
	FteInputNode                  IIBInputNodeTypeValue
	HTTPAsyncResponse             IIBInputNodeTypeValue
	JdEdwardsInputNode            IIBInputNodeTypeValue
	JmsClientInputNode            IIBInputNodeTypeValue
	LabelNode                     IIBInputNodeTypeValue
	MqInputNode                   IIBInputNodeTypeValue
	PeopleSoftInputNode           IIBInputNodeTypeValue
	RestAsyncResponse             IIBInputNodeTypeValue
	RestRequest                   IIBInputNodeTypeValue
	SAPInputNode                  IIBInputNodeTypeValue
	ScaAsyncResponseNode          IIBInputNodeTypeValue
	ScaInputNode                  IIBInputNodeTypeValue
	SiebelInputNode               IIBInputNodeTypeValue
	SoapInputNode                 IIBInputNodeTypeValue
	TcpipClientInputNode          IIBInputNodeTypeValue
	TcpipClientRequestNode        IIBInputNodeTypeValue
	TcpipServerInputNode          IIBInputNodeTypeValue
	TcpipServerRequestNode        IIBInputNodeTypeValue
	TimeoutNotificationNode       IIBInputNodeTypeValue
	WsInputNode                   IIBInputNodeTypeValue
}{
	"CALLABLE_FLOW_ASYNC_RESPONSE_NODE",
	"CALLABLE_FLOW_INPUT_NODE",
	"DATABASE_INPUT_NODE",
	"DOTNET_INPUT_NODE",
	"EMAIL_INPUT_NODE",
	"EVENT_INPUT",
	"EVENT_INPUT_NODE",
	"FILE_INPUT_NODE",
	"FTE_INPUT_NODE",
	"HTTP_ASYNC_RESPONSE",
	"JD_EDWARDS_INPUT_NODE",
	"JMS_CLIENT_INPUT_NODE",
	"LABEL_NODE",
	"MQ_INPUT_NODE",
	"PEOPLE_SOFT_INPUT_NODE",
	"REST_ASYNC_RESPONSE",
	"REST_REQUEST",
	"SAP_INPUT_NODE",
	"SCA_ASYNC_RESPONSE_NODE",
	"SCA_INPUT_NODE",
	"SIEBEL_INPUT_NODE",
	"SOAP_INPUT_NODE",
	"TCPIP_CLIENT_INPUT_NODE",
	"TCPIP_CLIENT_REQUEST_NODE",
	"TCPIP_SERVER_INPUT_NODE",
	"TCPIP_SERVER_REQUEST_NODE",
	"TIMEOUT_NOTIFICATION_NODE",
	"WS_INPUT_NODE",
}
