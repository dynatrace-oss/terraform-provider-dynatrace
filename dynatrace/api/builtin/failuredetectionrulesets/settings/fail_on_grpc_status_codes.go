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

package failuredetectionrulesets

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type FailOnGrpcStatusCodes struct {
	StatusCodes string `json:"statusCodes"` // Status codes which indicate a failure on the server side
}

func (me *FailOnGrpcStatusCodes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"status_codes": {
			Type:        schema.TypeString,
			Description: "Status codes which indicate a failure on the server side",
			Required:    true,
		},
	}
}

func (me *FailOnGrpcStatusCodes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"status_codes": me.StatusCodes,
	})
}

func (me *FailOnGrpcStatusCodes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"status_codes": &me.StatusCodes,
	})
}
