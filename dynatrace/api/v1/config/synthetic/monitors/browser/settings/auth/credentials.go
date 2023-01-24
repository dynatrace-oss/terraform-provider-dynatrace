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

package auth

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Credentials The login credentials to bypass the browser login mask during a Navigate event
type Credentials struct {
	Type       string     `json:"type"`       // The type of authentication
	Credential Credential `json:"credential"` // A reference to the entry within the credential vault
}

func (me *Credentials) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of authentication",
			Required:    true,
		},
		"creds": {
			Type:        schema.TypeString,
			Description: "A reference to the entry within the credential vault",
			Required:    true,
		},
	}
}

func (me *Credentials) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", me.Type); err != nil {
		return err
	}
	if err := properties.Encode("creds", me.Credential.ID); err != nil {
		return err
	}
	return nil
}

func (me *Credentials) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	cred := new(Credential)
	if err := decoder.Decode("creds", &cred.ID); err != nil {
		return err
	}
	if len(cred.ID) > 0 {
		me.Credential = *cred
	}
	return nil
}

type Credential struct {
	ID string `json:"id"`
}
