/**
* @license
* Copyright 2026 Dynatrace LLC
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

package dataforwarding

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AwsConnection struct {
	Arn          string `json:"arn"`          // S3 Bucket ARN
	ConnectionID string `json:"connectionId"` // AWS connection
}

func (me *AwsConnection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"arn": {
			Type:        schema.TypeString,
			Description: "S3 Bucket ARN",
			Required:    true,
		},
		"connection_id": {
			Type:        schema.TypeString,
			Description: "AWS connection",
			Required:    true,
		},
	}
}

func (me *AwsConnection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"arn":           me.Arn,
		"connection_id": me.ConnectionID,
	})
}

func (me *AwsConnection) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"arn":           &me.Arn,
		"connection_id": &me.ConnectionID,
	})
}
