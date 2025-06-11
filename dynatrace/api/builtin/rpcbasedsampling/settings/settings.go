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

package rpcbasedsampling

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Enabled                           bool                 `json:"enabled"`                           // This setting is enabled (`true`) or disabled (`false`)
	EndpointName                      *string              `json:"endpointName,omitempty"`            // Specify the RPC endpoint name. If the endpoint name is empty, either remote operation name or remote service name must be specified that can be used for RPC matching.
	EndpointNameComparisonType        ComparisonType       `json:"endpointNameComparisonType"`        // Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
	Factor                            *SamplingScaleFactor `json:"factor,omitempty"`                  // Possible Values: `0`, `1`, `2`, `3`, `4`, `5`, `6`, `8`, `9`, `10`, `11`, `12`, `13`, `14` \n\n  - `0`: Increase capturing 128 times\n\n  - `1`: Increase capturing 64 times\n\n  - `2`: Increase capturing 32 times\n\n  - `3`: Increase capturing 16 times\n\n  - `4`: Increase capturing 8 times\n\n  - `5`: Increase capturing 4 times\n\n  - `6`: Increase capturing 2 times\n\n  - `8`: Reduce capturing by factor 2\n\n  - `9`: Reduce capturing by factor 4\n\n  - `10`: Reduce capturing by factor 8\n\n  - `11`: Reduce capturing by factor 16\n\n  - `12`: Reduce capturing by factor 32\n\n  - `13`: Reduce capturing by factor 64\n\n  - `14`: Reduce capturing by factor 128
	Ignore                            bool                 `json:"ignore"`                            // No Traces will be captured for matching RPC requests. This applies always, even if Adaptive Traffic Management is inactive.
	RemoteOperationName               *string              `json:"remoteOperationName,omitempty"`     // Specify the RPC operation name. If the remote operation name is empty, either remote service name or endpoint name must be specified that can be used for RPC matching.
	RemoteOperationNameComparisonType ComparisonType       `json:"remoteOperationNameComparisonType"` // Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
	RemoteServiceName                 *string              `json:"remoteServiceName,omitempty"`       // Specify the RPC remote service name. If the remote service name is empty, either remote operation name or endpoint name must be specified that can be used for RPC matching.
	RemoteServiceNameComparisonType   ComparisonType       `json:"remoteServiceNameComparisonType"`   // Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`
	Scope                             *string              `json:"-" scope:"scope"`                   // The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP, CLOUD_APPLICATION, CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.
	WireProtocolType                  WireProtocolType     `json:"wireProtocolType"`                  // Possible Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10` \n\n  - `1`: ADK\n\n  - `2`: DOTNET_REMOTING\n\n  - `3`: DOTNET_REMOTING_TCP\n\n  - `4`: DOTNET_REMOTING_HTTP\n\n  - `5`: DOTNET_REMOTING_XMLRPC\n\n  - `6`: GRPC\n\n  - `7`: GRPC_BIDI\n\n  - `8`: GRPC_UNARY\n\n  - `9`: GRPC_SERVERSTREAM\n\n  - `10`: GRPC_CLIENTSTREAM
	InsertAfter                       string               `json:"-"`
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"endpoint_name": {
			Type:        schema.TypeString,
			Description: "Specify the RPC endpoint name. If the endpoint name is empty, either remote operation name or remote service name must be specified that can be used for RPC matching.",
			Optional:    true, // nullable
		},
		"endpoint_name_comparison_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`",
			Required:    true,
		},
		"factor": {
			Type:        schema.TypeString,
			Description: "Possible Values: `0`, `1`, `2`, `3`, `4`, `5`, `6`, `8`, `9`, `10`, `11`, `12`, `13`, `14` \n\n  - `0`: Increase capturing 128 times\n\n  - `1`: Increase capturing 64 times\n\n  - `2`: Increase capturing 32 times\n\n  - `3`: Increase capturing 16 times\n\n  - `4`: Increase capturing 8 times\n\n  - `5`: Increase capturing 4 times\n\n  - `6`: Increase capturing 2 times\n\n  - `8`: Reduce capturing by factor 2\n\n  - `9`: Reduce capturing by factor 4\n\n  - `10`: Reduce capturing by factor 8\n\n  - `11`: Reduce capturing by factor 16\n\n  - `12`: Reduce capturing by factor 32\n\n  - `13`: Reduce capturing by factor 64\n\n  - `14`: Reduce capturing by factor 128",
			Optional:    true, // precondition
		},
		"ignore": {
			Type:        schema.TypeBool,
			Description: "No Traces will be captured for matching RPC requests. This applies always, even if Adaptive Traffic Management is inactive.",
			Required:    true,
		},
		"remote_operation_name": {
			Type:        schema.TypeString,
			Description: "Specify the RPC operation name. If the remote operation name is empty, either remote service name or endpoint name must be specified that can be used for RPC matching.",
			Optional:    true, // nullable
		},
		"remote_operation_name_comparison_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`",
			Required:    true,
		},
		"remote_service_name": {
			Type:        schema.TypeString,
			Description: "Specify the RPC remote service name. If the remote service name is empty, either remote operation name or endpoint name must be specified that can be used for RPC matching.",
			Optional:    true, // nullable
		},
		"remote_service_name_comparison_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `CONTAINS`, `DOES_NOT_CONTAIN`, `DOES_NOT_END_WITH`, `DOES_NOT_EQUAL`, `DOES_NOT_START_WITH`, `ENDS_WITH`, `EQUALS`, `STARTS_WITH`",
			Required:    true,
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (PROCESS_GROUP_INSTANCE, PROCESS_GROUP, CLOUD_APPLICATION, CLOUD_APPLICATION_NAMESPACE, KUBERNETES_CLUSTER, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
		"wire_protocol_type": {
			Type:        schema.TypeString,
			Description: "Possible Values: `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10` \n\n  - `1`: ADK\n\n  - `2`: DOTNET_REMOTING\n\n  - `3`: DOTNET_REMOTING_TCP\n\n  - `4`: DOTNET_REMOTING_HTTP\n\n  - `5`: DOTNET_REMOTING_XMLRPC\n\n  - `6`: GRPC\n\n  - `7`: GRPC_BIDI\n\n  - `8`: GRPC_UNARY\n\n  - `9`: GRPC_SERVERSTREAM\n\n  - `10`: GRPC_CLIENTSTREAM",
			Required:    true,
		},
		"insert_after": {
			Type:        schema.TypeString,
			Description: "Because this resource allows for ordering you may specify the ID of the resource instance that comes before this instance regarding order. If not specified when creating the setting will be added to the end of the list. If not specified during update the order will remain untouched",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enabled":                               me.Enabled,
		"endpoint_name":                         me.EndpointName,
		"endpoint_name_comparison_type":         me.EndpointNameComparisonType,
		"factor":                                me.Factor,
		"ignore":                                me.Ignore,
		"remote_operation_name":                 me.RemoteOperationName,
		"remote_operation_name_comparison_type": me.RemoteOperationNameComparisonType,
		"remote_service_name":                   me.RemoteServiceName,
		"remote_service_name_comparison_type":   me.RemoteServiceNameComparisonType,
		"scope":                                 me.Scope,
		"wire_protocol_type":                    me.WireProtocolType,
		"insert_after":                          me.InsertAfter,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.Factor == nil) && (!me.Ignore) {
		return fmt.Errorf("'factor' must be specified if 'ignore' is set to '%v'", me.Ignore)
	}
	if (me.Factor != nil) && (me.Ignore) {
		return fmt.Errorf("'factor' must not be specified if 'ignore' is set to '%v'", me.Ignore)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enabled":                               &me.Enabled,
		"endpoint_name":                         &me.EndpointName,
		"endpoint_name_comparison_type":         &me.EndpointNameComparisonType,
		"factor":                                &me.Factor,
		"ignore":                                &me.Ignore,
		"remote_operation_name":                 &me.RemoteOperationName,
		"remote_operation_name_comparison_type": &me.RemoteOperationNameComparisonType,
		"remote_service_name":                   &me.RemoteServiceName,
		"remote_service_name_comparison_type":   &me.RemoteServiceNameComparisonType,
		"scope":                                 &me.Scope,
		"wire_protocol_type":                    &me.WireProtocolType,
		"insert_after":                          &me.InsertAfter,
	})
}
