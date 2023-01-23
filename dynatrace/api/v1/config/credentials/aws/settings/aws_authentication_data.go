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
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AWSAuthenticationData A credentials for the AWS authentication.
type AWSAuthenticationData struct {
	KeyBasedAuthentication  *KeyBasedAuthentication    `json:"keyBasedAuthentication,omitempty"`  // The credentials for the key-based authentication.
	RoleBasedAuthentication *RoleBasedAuthentication   `json:"roleBasedAuthentication,omitempty"` // The credentials for the role-based authentication.
	Type                    Type                       `json:"type"`                              // The type of the authentication: role-based or key-based.
	Unknowns                map[string]json.RawMessage `json:"-"`
}

func (aad *AWSAuthenticationData) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Description: "the access key",
			Optional:    true,
		},
		"secret_key": {
			Type:        schema.TypeString,
			Description: "the secret access key",
			Optional:    true,
			Sensitive:   true,
		},
		"account_id": {
			Type:        schema.TypeString,
			Description: "the ID of the Amazon account",
			Optional:    true,
		},
		"external_id": {
			Type:        schema.TypeString,
			Description: "the external ID token for setting an IAM role. You can obtain it with the `GET /aws/iamExternalId` request",
			Optional:    true,
		},
		"iam_role": {
			Type:        schema.TypeString,
			Description: "the IAM role to be used by Dynatrace to get monitoring data",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (aad *AWSAuthenticationData) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), aad); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &aad.Unknowns); err != nil {
			return err
		}
		delete(aad.Unknowns, "access_key")
		delete(aad.Unknowns, "secret_key")
		delete(aad.Unknowns, "account_id")
		delete(aad.Unknowns, "external_id")
		delete(aad.Unknowns, "iam_role")
		if len(aad.Unknowns) == 0 {
			aad.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("access_key"); ok {
		if aad.KeyBasedAuthentication == nil {
			aad.KeyBasedAuthentication = new(KeyBasedAuthentication)
		}
		aad.Type = Types.Keys
		aad.KeyBasedAuthentication.AccessKey = value.(string)
	}
	if value, ok := decoder.GetOk("secret_key"); ok {
		if aad.KeyBasedAuthentication == nil {
			aad.KeyBasedAuthentication = new(KeyBasedAuthentication)
		}
		aad.Type = Types.Keys
		aad.KeyBasedAuthentication.SecretKey = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("account_id"); ok {
		if aad.RoleBasedAuthentication == nil {
			aad.RoleBasedAuthentication = new(RoleBasedAuthentication)
		}
		aad.Type = Types.Role
		aad.RoleBasedAuthentication.AccountID = value.(string)
	}
	if value, ok := decoder.GetOk("external_id"); ok {
		if aad.RoleBasedAuthentication == nil {
			aad.RoleBasedAuthentication = new(RoleBasedAuthentication)
		}
		aad.Type = Types.Role
		aad.RoleBasedAuthentication.ExternalID = opt.NewString(value.(string))
	}
	if value, ok := decoder.GetOk("iam_role"); ok {
		if aad.RoleBasedAuthentication == nil {
			aad.RoleBasedAuthentication = new(RoleBasedAuthentication)
		}
		aad.Type = Types.Role
		aad.RoleBasedAuthentication.IamRole = value.(string)
	}
	return nil
}

func (aad *AWSAuthenticationData) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(aad.Unknowns); err != nil {
		return err
	}
	if aad.KeyBasedAuthentication != nil {
		if err := properties.Encode("access_key", aad.KeyBasedAuthentication.AccessKey); err != nil {
			return err
		}
		if err := properties.Encode("secret_key", "${state.secret_value}"); err != nil {
			return err
		}
	}
	if aad.RoleBasedAuthentication != nil {
		if err := properties.Encode("account_id", aad.RoleBasedAuthentication.AccountID); err != nil {
			return err
		}
		if err := properties.Encode("external_id", aad.RoleBasedAuthentication.ExternalID); err != nil {
			return err
		}
		if err := properties.Encode("iam_role", aad.RoleBasedAuthentication.IamRole); err != nil {
			return err
		}

	}
	return nil
}

// UnmarshalJSON provides custom JSON deserialization
func (aad *AWSAuthenticationData) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["keyBasedAuthentication"]; found {
		if err := json.Unmarshal(v, &aad.KeyBasedAuthentication); err != nil {
			return err
		}
	}
	if v, found := m["roleBasedAuthentication"]; found {
		if err := json.Unmarshal(v, &aad.RoleBasedAuthentication); err != nil {
			return err
		}
	}
	if v, found := m["type"]; found {
		if err := json.Unmarshal(v, &aad.Type); err != nil {
			return err
		}
	} else {
		if aad.RoleBasedAuthentication != nil {
			aad.Type = Types.Role
		} else if aad.KeyBasedAuthentication != nil {
			aad.Type = Types.Keys
		}
	}
	delete(m, "keyBasedAuthentication")
	delete(m, "roleBasedAuthentication")
	delete(m, "type")
	if len(m) > 0 {
		aad.Unknowns = m
	}
	return nil
}

// MarshalJSON provides custom JSON serialization
func (aad *AWSAuthenticationData) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(aad.Unknowns) > 0 {
		for k, v := range aad.Unknowns {
			m[k] = v
		}
	}
	if aad.KeyBasedAuthentication != nil {
		rawMessage, err := json.Marshal(aad.KeyBasedAuthentication)
		if err != nil {
			return nil, err
		}
		m["keyBasedAuthentication"] = rawMessage
	}
	if aad.RoleBasedAuthentication != nil {
		rawMessage, err := json.Marshal(aad.RoleBasedAuthentication)
		if err != nil {
			return nil, err
		}
		m["roleBasedAuthentication"] = rawMessage
	}
	rawMessage, err := json.Marshal(aad.Type)
	if err != nil {
		return nil, err
	}
	m["type"] = rawMessage
	return json.Marshal(m)
}
