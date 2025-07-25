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

package role_arn

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Name            string
	AWSConnectionID string
	RoleARN         string `json:"roleArn"` // The ARN of the AWS role that should be assumed
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"aws_connection_id": {
			Type:        schema.TypeString,
			Description: "The ID of a `dynatrace_aws_connection` resource instance for which to define the AWS Role ARN",
			Required:    true,
			ForceNew:    true,
		},
		"role_arn": {
			Type:        schema.TypeString,
			Description: "The ARN of the AWS role that should be assumed.",
			ForceNew:    true,
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"aws_connection_id": me.AWSConnectionID,
		"role_arn":          me.RoleARN,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"aws_connection_id": &me.AWSConnectionID,
		"role_arn":          &me.RoleARN,
	})
}
