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

package httpparameters

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type HttpResponseCodes struct {
	ClientSideErrors                    string `json:"clientSideErrors"`                    // HTTP response codes which indicate client side errors
	FailOnMissingResponseCodeClientSide bool   `json:"failOnMissingResponseCodeClientSide"` // Treat missing HTTP response code as client side error
	FailOnMissingResponseCodeServerSide bool   `json:"failOnMissingResponseCodeServerSide"` // Treat missing HTTP response code as server side errors
	ServerSideErrors                    string `json:"serverSideErrors"`                    // HTTP response codes which indicate an error on the server side
}

func (me *HttpResponseCodes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_side_errors": {
			Type:        schema.TypeString,
			Description: "HTTP response codes which indicate client side errors",
			Required:    true,
		},
		"fail_on_missing_response_code_client_side": {
			Type:        schema.TypeBool,
			Description: "Treat missing HTTP response code as client side error",
			Required:    true,
		},
		"fail_on_missing_response_code_server_side": {
			Type:        schema.TypeBool,
			Description: "Treat missing HTTP response code as server side errors",
			Required:    true,
		},
		"server_side_errors": {
			Type:        schema.TypeString,
			Description: "HTTP response codes which indicate an error on the server side",
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
