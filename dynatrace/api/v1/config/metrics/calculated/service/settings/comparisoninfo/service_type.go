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

// ServiceType Comparison for `SERVICE_TYPE` attributes.
type ServiceType struct {
	BaseComparisonInfo
	Comparison ServiceTypeComparison `json:"comparison"`       // Operator of the comparision. You can reverse it by setting **negate** to `true`.
	Value      *ServiceTypeValue     `json:"value,omitempty"`  // The value to compare to.
	Values     []ServiceTypeValue    `json:"values,omitempty"` // The values to compare to.
}

func (me *ServiceType) GetType() Type {
	return Types.ServiceType
}

func (me *ServiceType) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"values": {
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    1,
			Description: "The values to compare to. Possible values are `BACKGROUND_ACTIVITY`, `CICS_SERVICE`, `CUSTOM_SERVICE`, `DATABASE_SERVICE`, `ENTERPRISE_SERVICE_BUS_SERVICE`, `EXTERNAL`, `IBM_INTEGRATION_BUS_SERVICE`, `IMS_SERVICE`, `MESSAGING_SERVICE`, `RMI_SERVICE`, `RPC_SERVICE`, `WEB_REQUEST_SERVICE` and `WEB_SERVICE`",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The value to compare to. Possible values are `BACKGROUND_ACTIVITY`, `CICS_SERVICE`, `CUSTOM_SERVICE`, `DATABASE_SERVICE`, `ENTERPRISE_SERVICE_BUS_SERVICE`, `EXTERNAL`, `IBM_INTEGRATION_BUS_SERVICE`, `IMS_SERVICE`, `MESSAGING_SERVICE`, `RMI_SERVICE`, `RPC_SERVICE`, `WEB_REQUEST_SERVICE` and `WEB_SERVICE`",
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

func (me *ServiceType) MarshalHCL(properties hcl.Properties) error {
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

func (me *ServiceType) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"values":   &me.Values,
		"value":    &me.Value,
		"operator": &me.Comparison,
		"unknowns": &me.Unknowns,
	})
}

func (me *ServiceType) MarshalJSON() ([]byte, error) {
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

func (me *ServiceType) UnmarshalJSON(data []byte) error {
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

// ServiceTypeComparison Operator of the comparision. You can reverse it by setting **negate** to `true`.
type ServiceTypeComparison string

// ServiceTypeComparisons offers the known enum values
var ServiceTypeComparisons = struct {
	Equals      ServiceTypeComparison
	EqualsAnyOf ServiceTypeComparison
	Exists      ServiceTypeComparison
}{
	"EQUALS",
	"EQUALS_ANY_OF",
	"EXISTS",
}

// ServiceTypeValue The value to compare to.
type ServiceTypeValue string

// ServiceTypeValues offers the known enum values
var ServiceTypeValues = struct {
	BackgroundActivity          ServiceTypeValue
	CICSService                 ServiceTypeValue
	CustomService               ServiceTypeValue
	DatabaseService             ServiceTypeValue
	EnterpriseServiceBusService ServiceTypeValue
	External                    ServiceTypeValue
	IBMIntegrationBusService    ServiceTypeValue
	IMSService                  ServiceTypeValue
	MessagingService            ServiceTypeValue
	RMIService                  ServiceTypeValue
	RPCService                  ServiceTypeValue
	WebRequestService           ServiceTypeValue
	WebService                  ServiceTypeValue
}{
	"BACKGROUND_ACTIVITY",
	"CICS_SERVICE",
	"CUSTOM_SERVICE",
	"DATABASE_SERVICE",
	"ENTERPRISE_SERVICE_BUS_SERVICE",
	"EXTERNAL",
	"IBM_INTEGRATION_BUS_SERVICE",
	"IMS_SERVICE",
	"MESSAGING_SERVICE",
	"RMI_SERVICE",
	"RPC_SERVICE",
	"WEB_REQUEST_SERVICE",
	"WEB_SERVICE",
}
