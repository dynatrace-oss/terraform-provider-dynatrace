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

package parameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// httpResponseCodes. HTTP response code settings that control which response codes are treated as failures on the server side and client side.
type HttpResponseCodes struct {
	ClientSideErrors                    string `json:"clientSideErrors"`                    // A list of HTTP response code ranges and individual values that are treated as client-side errors. The format is a comma-separated list of ranges and values (e.g., `400-499, 503, 510-599`). Default: `400-599`.
	FailOnMissingResponseCodeClientSide bool   `json:"failOnMissingResponseCodeClientSide"` // If `true`, a missing HTTP response code on the client side is treated as a failure. Missing response codes can indicate a fire-and-forget call, a timeout, or an error. Default: `false`.
	FailOnMissingResponseCodeServerSide bool   `json:"failOnMissingResponseCodeServerSide"` // If `true`, a missing HTTP response code on the server side is treated as a failure. Missing response codes can indicate a fire-and-forget call, a timeout, or an error. Default: `false`.
	ServerSideErrors                    string `json:"serverSideErrors"`                    // A list of HTTP response code ranges and individual values that are treated as server-side errors. The format is a comma-separated list of ranges and values (e.g., `500-599, 402, 405-499`). Default: `500-599`.
}

func (me *HttpResponseCodes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_side_errors": {
			Type:        schema.TypeString,
			Description: "A list of HTTP response code ranges and individual values that are treated as client-side errors. The format is a comma-separated list of ranges and values (e.g., `400-499, 503, 510-599`). Default: `400-599`.",
			Required:    true,
		},
		"fail_on_missing_response_code_client_side": {
			Type:        schema.TypeBool,
			Description: "If `true`, a missing HTTP response code on the client side is treated as a failure. Missing response codes can indicate a fire-and-forget call, a timeout, or an error. Default: `false`.",
			Required:    true,
		},
		"fail_on_missing_response_code_server_side": {
			Type:        schema.TypeBool,
			Description: "If `true`, a missing HTTP response code on the server side is treated as a failure. Missing response codes can indicate a fire-and-forget call, a timeout, or an error. Default: `false`.",
			Required:    true,
		},
		"server_side_errors": {
			Type:        schema.TypeString,
			Description: "A list of HTTP response code ranges and individual values that are treated as server-side errors. The format is a comma-separated list of ranges and values (e.g., `500-599, 402, 405-499`). Default: `500-599`.",
			Required:    true,
		},
	}
}

func (me *HttpResponseCodes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"client_side_errors":                        me.ClientSideErrors,
		"fail_on_missing_response_code_client_side": me.FailOnMissingResponseCodeClientSide,
		"fail_on_missing_response_code_server_side": me.FailOnMissingResponseCodeServerSide,
		"server_side_errors":                        me.ServerSideErrors,
	})
}

func (me *HttpResponseCodes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"client_side_errors":                        &me.ClientSideErrors,
		"fail_on_missing_response_code_client_side": &me.FailOnMissingResponseCodeClientSide,
		"fail_on_missing_response_code_server_side": &me.FailOnMissingResponseCodeServerSide,
		"server_side_errors":                        &me.ServerSideErrors,
	})
}
