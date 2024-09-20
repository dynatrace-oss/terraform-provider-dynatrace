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

package awsconnection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/export/sensitive"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type WebIdentity struct {
	PolicyArns []string `json:"policyArns,omitempty"` // An optional list of policies that can be used to restrict the AWS role
	RoleArn    string   `json:"roleArn"`              // The ARN of the AWS role that should be assumed
}

func (me *WebIdentity) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"policy_arns": {
			Type:        schema.TypeList,
			Description: "An optional list of policies that can be used to restrict the AWS role",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
			Sensitive:   true,
		},
		"role_arn": {
			Type:        schema.TypeString,
			Description: "The ARN of the AWS role that should be assumed",
			Required:    true,
			Sensitive:   true,
			ForceNew:    true,
		},
	}
}

func (me *WebIdentity) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(sensitive.ConditionalIgnoreChangesMap(
		me.Schema(), map[string]any{
			"policy_arns": sensitive.SecretifyExact(me.PolicyArns),
			"role_arn":    sensitive.SecretValueExact,
		},
	))
}

func (me *WebIdentity) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"policy_arns": &me.PolicyArns,
		"role_arn":    &me.RoleArn,
	})
}

func (me *WebIdentity) FillDemoValues() []string {
	me.PolicyArns = []string{"#######"}
	me.RoleArn = "#######"
	return []string{"Please fill in the policy and role ARNs"}
}
