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

type Wrapper struct {
	Negate     bool
	Comparison ComparisonInfo
}

func (me *Wrapper) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"negate": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Reverse the comparison **operator**. For example, it turns **equals** into **does not equal**",
		},
		"boolean": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Boolean Comparison for `BOOLEAN` attributes",
			Elem:        &schema.Resource{Schema: new(Boolean).Schema()},
		},
		"esb_input_node_type": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Type-specific comparison information for attributes of type 'ESB_INPUT_NODE_TYPE'",
			Elem:        &schema.Resource{Schema: new(ESBInputNodeType).Schema()},
		},
		"failed_state": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `FAILED_STATE` attributes",
			Elem:        &schema.Resource{Schema: new(FailedState).Schema()},
		},
		"failure_reason": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `FAILURE_REASON` attributes",
			Elem:        &schema.Resource{Schema: new(FailureReason).Schema()},
		},
		"fast_string": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `FAST_STRING` attributes. Use it for all service property attributes",
			Elem:        &schema.Resource{Schema: new(FastString).Schema()},
		},
		"flaw_state": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `FLAW_STATE` attributes",
			Elem:        &schema.Resource{Schema: new(FlawState).Schema()},
		},
		"http_method": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `HTTP_METHOD` attributes",
			Elem:        &schema.Resource{Schema: new(HTTPMethod).Schema()},
		},
		"http_status_class": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `HTTP_STATUS_CLASS` attributes",
			Elem:        &schema.Resource{Schema: new(HTTPStatusClass).Schema()},
		},
		"iib_input_node_type": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `IIB_INPUT_NODE_TYPE` attributes",
			Elem:        &schema.Resource{Schema: new(IIBInputNodeType).Schema()},
		},
		"number": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `NUMBER` attributes",
			Elem:        &schema.Resource{Schema: new(Number).Schema()},
		},
		"number_request_attribute": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `NUMBER_REQUEST_ATTRIBUTE` attributes",
			Elem:        &schema.Resource{Schema: new(NumberRequestAttribute).Schema()},
		},
		"service_type": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `SERVICE_TYPE` attributes",
			Elem:        &schema.Resource{Schema: new(ServiceType).Schema()},
		},
		"string": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `STRING` attributes",
			Elem:        &schema.Resource{Schema: new(String).Schema()},
		},
		"string_request_attribute": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `STRING_REQUEST_ATTRIBUTE` attributes",
			Elem:        &schema.Resource{Schema: new(StringRequestAttribute).Schema()},
		},
		"tag": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `TAG` attributes",
			Elem:        &schema.Resource{Schema: new(Tag).Schema()},
		},
		"zos_call_type": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `ZOS_CALL_TYPE` attributes",
			Elem:        &schema.Resource{Schema: new(ZOSCallType).Schema()},
		},
		"generic": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			MinItems:    1,
			Description: "Comparison for `NUMBER` attributes",
			Elem:        &schema.Resource{Schema: new(BaseComparisonInfo).Schema()},
		},
	}
}

func (me *Wrapper) MarshalHCL(properties hcl.Properties) error {
	properties.Encode("negate", me.Comparison.IsNegate())
	switch cmp := me.Comparison.(type) {
	case *Boolean:
		if err := properties.Encode("boolean", cmp); err != nil {
			return nil
		}
		return nil
	case *ESBInputNodeType:
		if err := properties.Encode("esb_input_node_type", cmp); err != nil {
			return nil
		}
		return nil
	case *FailedState:
		if err := properties.Encode("failed_state", cmp); err != nil {
			return nil
		}
		return nil
	case *FailureReason:
		if err := properties.Encode("failure_reason", cmp); err != nil {
			return nil
		}
		return nil
	case *FastString:
		if err := properties.Encode("fast_string", cmp); err != nil {
			return nil
		}
		return nil
	case *FlawState:
		if err := properties.Encode("flaw_state", cmp); err != nil {
			return nil
		}
		return nil
	case *HTTPMethod:
		if err := properties.Encode("http_method", cmp); err != nil {
			return nil
		}
		return nil
	case *HTTPStatusClass:
		if err := properties.Encode("http_status_class", cmp); err != nil {
			return nil
		}
		return nil
	case *IIBInputNodeType:
		if err := properties.Encode("iib_input_node_type", cmp); err != nil {
			return nil
		}
		return nil
	case *NumberRequestAttribute:
		if err := properties.Encode("number_request_attribute", cmp); err != nil {
			return nil
		}
		return nil
	case *Number:
		if err := properties.Encode("number", cmp); err != nil {
			return nil
		}
		return nil
	case *ServiceType:
		if err := properties.Encode("service_type", cmp); err != nil {
			return nil
		}
		return nil
	case *StringRequestAttribute:
		if err := properties.Encode("string_request_attribute", cmp); err != nil {
			return nil
		}
		return nil
	case *String:
		if err := properties.Encode("string", cmp); err != nil {
			return nil
		}
		return nil
	case *Tag:
		if err := properties.Encode("tag", cmp); err != nil {
			return nil
		}
		return nil
	case *ZOSCallType:
		if err := properties.Encode("zos_call_type", cmp); err != nil {
			return nil
		}
		return nil
	case *BaseComparisonInfo:
		if err := properties.Encode("generic", cmp); err != nil {
			return nil
		}
		return nil
	default:
		return fmt.Errorf("cannot HCL marshal objects (xxx) of type %T", cmp)
	}
}

func (me *Wrapper) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("negate", &me.Negate); err != nil {
		return err
	}
	var err error
	var cmp any
	if cmp, err = decoder.DecodeAny(map[string]any{
		"boolean":                  new(Boolean),
		"esb_input_node_type":      new(ESBInputNodeType),
		"failed_state":             new(FailedState),
		"failure_reason":           new(FailureReason),
		"fast_string":              new(FastString),
		"flaw_state":               new(FlawState),
		"http_method":              new(HTTPMethod),
		"http_status_class":        new(HTTPStatusClass),
		"iib_input_node_type":      new(IIBInputNodeType),
		"number":                   new(Number),
		"number_request_attribute": new(NumberRequestAttribute),
		"service_type":             new(ServiceType),
		"string":                   new(String),
		"string_request_attribute": new(StringRequestAttribute),
		"tag":                      new(Tag),
		"zos_call_type":            new(ZOSCallType),
		"generic":                  new(BaseComparisonInfo)}); err != nil {
		return err
	}
	if cmp != nil {
		me.Comparison = cmp.(ComparisonInfo)
		me.Comparison.SetNegate(me.Negate)
	}
	return nil
}

func (me *Wrapper) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	var compType string
	if err := properties.UnmarshalAll(map[string]any{
		"negate": &me.Negate,
		"type":   &compType,
	}); err != nil {
		return err
	}
	switch compType {
	case "STRING":
		cfg := new(String)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "NUMBER":
		cfg := new(Number)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "BOOLEAN":
		cfg := new(Boolean)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "HTTP_METHOD":
		cfg := new(HTTPMethod)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "STRING_REQUEST_ATTRIBUTE":
		cfg := new(StringRequestAttribute)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "NUMBER_REQUEST_ATTRIBUTE":
		cfg := new(NumberRequestAttribute)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "ZOS_CALL_TYPE":
		cfg := new(ZOSCallType)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "IIB_INPUT_NODE_TYPE":
		cfg := new(IIBInputNodeType)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "ESB_INPUT_NODE_TYPE":
		cfg := new(ESBInputNodeType)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "FAILED_STATE":
		cfg := new(FailedState)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "FLAW_STATE":
		cfg := new(FlawState)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "FAILURE_REASON":
		cfg := new(FailureReason)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "HTTP_STATUS_CLASS":
		cfg := new(HTTPStatusClass)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "TAG":
		cfg := new(Tag)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "FAST_STRING":
		cfg := new(FastString)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	case "SERVICE_TYPE":
		cfg := new(ServiceType)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	default:
		cfg := new(BaseComparisonInfo)
		if err := json.Unmarshal(data, &cfg); err != nil {
			return err
		}
		cfg.Negate = me.Negate
		me.Comparison = cfg
	}
	return nil
}
