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

import "encoding/json"

// RoleBasedAuthentication The credentials for the role-based authentication.
type RoleBasedAuthentication struct {
	AccountID  string                     `json:"accountId"`            // The ID of the Amazon account.
	ExternalID *string                    `json:"externalId,omitempty"` // The external ID token for setting an IAM role.   You can obtain it with the `GET /aws/iamExternalId` request.
	IamRole    string                     `json:"iamRole"`              // The IAM role to be used by Dynatrace to get monitoring data.
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (rba *RoleBasedAuthentication) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["accountId"]; found {
		if err := json.Unmarshal(v, &rba.AccountID); err != nil {
			return err
		}
	}
	if v, found := m["iamRole"]; found {
		if err := json.Unmarshal(v, &rba.IamRole); err != nil {
			return err
		}
	}
	if rba.ExternalID != nil {
		if v, found := m["externalId"]; found {
			if err := json.Unmarshal(v, &rba.ExternalID); err != nil {
				return err
			}
		}
	}
	delete(m, "accountId")
	delete(m, "iamRole")
	delete(m, "externalId")
	if len(m) > 0 {
		rba.Unknowns = m
	}
	return nil
}

func (rba *RoleBasedAuthentication) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(rba.Unknowns) > 0 {
		for k, v := range rba.Unknowns {
			m[k] = v
		}
	}
	{
		rawMessage, err := json.Marshal(rba.AccountID)
		if err != nil {
			return nil, err
		}
		m["accountId"] = rawMessage
	}
	{
		rawMessage, err := json.Marshal(rba.IamRole)
		if err != nil {
			return nil, err
		}
		m["iamRole"] = rawMessage
	}
	if rba.ExternalID != nil {
		rawMessage, err := json.Marshal(rba.ExternalID)
		if err != nil {
			return nil, err
		}
		m["externalId"] = rawMessage
	}
	return json.Marshal(m)
}
