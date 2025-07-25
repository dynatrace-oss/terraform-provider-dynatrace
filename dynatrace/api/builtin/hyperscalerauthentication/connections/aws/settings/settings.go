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

type Settings struct {
	Name                       string                            `json:"name"`
	Type                       Type                              `json:"type"`
	AWSRoleBasedAuthentication *AwsRoleBasedAuthenticationConfig `json:"awsRoleBasedAuthentication,omitempty"`
	AWSWebIdentity             *AWSWebIdentity                   `json:"awsWebIdentity,omitempty"`
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
			ForceNew:    true,
		},
		"web_identity": {
			Type:         schema.TypeList,
			Description:  "Configuration required for authenticating via AWS Web Identity",
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(AWSWebIdentity).Schema()},
			ExactlyOneOf: []string{"web_identity", "role_based_auth"},
			MinItems:     1,
			MaxItems:     1,
			ForceNew:     true,
		},
		"role_based_auth": {
			Type:         schema.TypeList,
			Description:  "Configuration required for authenticating via AWS Role Based Authentication",
			Optional:     true,
			Elem:         &schema.Resource{Schema: new(AwsRoleBasedAuthenticationConfig).Schema()},
			ExactlyOneOf: []string{"web_identity", "role_based_auth"},
			MinItems:     1,
			MaxItems:     1,
			ForceNew:     true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":            me.Name,
		"web_identity":    me.AWSWebIdentity,
		"role_based_auth": me.AWSRoleBasedAuthentication,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.DecodeAll(map[string]any{
		"name":            &me.Name,
		"web_identity":    &me.AWSWebIdentity,
		"role_based_auth": &me.AWSRoleBasedAuthentication,
	}); err != nil {
		return err
	}

	if me.AWSWebIdentity != nil {
		me.Type = Types.AWSWebIdentity
	} else if me.AWSRoleBasedAuthentication != nil {
		me.Type = Types.AWSRoleBasedAuthentication
	}

	return nil
}
