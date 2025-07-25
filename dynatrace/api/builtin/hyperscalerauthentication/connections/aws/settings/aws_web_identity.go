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

package aws

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AWSWebIdentity struct {
	RoleARN   string                      `json:"roleArn"`   // The ARN of the AWS role that should be assumed
	Consumers []ConsumersOfAwsWebIdentity `json:"consumers"` // Dynatrace integrations that can use this connection
}

func (me *AWSWebIdentity) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// "role_arn": {
		// 	Type:        schema.TypeString,
		// 	Description: "The ARN of the AWS role that should be assumed. Omit this attribute and configure it using `dynatrace_aws_connection_role_arn` later to avoid circular dependencies.",
		// 	Optional:    true,
		// 	Computed:    true,
		// },
		"consumers": {
			Type:        schema.TypeSet,
			Description: "Dynatrace integrations that can use this connection. Possible values: `APP:dynatrace.aws.connector` (AWS Connector) and `APP:dynatrace.biz.carbon` (Cost & Carbon Optimization)",
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			MaxItems:    1,
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (me *AWSWebIdentity) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		// "role_arn":  me.RoleARN,
		"consumers": me.Consumers,
	})
}

func (me *AWSWebIdentity) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		// "role_arn":  &me.RoleARN,
		"consumers": &me.Consumers,
	})
}
