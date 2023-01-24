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

package http

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AuthenticationType is a type alias, nailing down currently supported AuthenticationTypes to `BASIC_AUTHENTICATION`, `NTLM` and `KERBEROS`. Additional values ARE however possible.
type AuthenticationType string

// AuthenticationTypes hints the currently supported AuthenticationTypes to `BASIC_AUTHENTICATION`, `NTLM` and `KERBEROS`. Additional values ARE however possible.
var AuthenticationTypes = struct {
	Basic    AuthenticationType
	NTML     AuthenticationType
	Kerberos AuthenticationType
}{
	AuthenticationType("BASIC_AUTHENTICATION"),
	AuthenticationType("NTML"),
	AuthenticationType("KERBEROS"),
}

// Authentication represents authentication options for a HTTP Request
type Authentication struct {
	Type        AuthenticationType `json:"type"`                // The type of authentication. Possible values are `BASIC_AUTHENTICATION`, `NTLM` and `KERBEROS`
	Credentials string             `json:"credentials"`         // The ID of the credentials within the Dynatrace Credentials Vault.
	RealmName   *string            `json:"realmName,omitempty"` // The Realm Name. Valid and required only if the type of authentication is `KERBEROS`
	KdcIP       *string            `json:"kdcIp,omitempty"`     // The KDC IP. Valid and required only if the type of authentication is `KERBEROS`
}

// Schema provides the schema map for the terraform provider
func (me *Authentication) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of authentication. Possible values are `BASIC_AUTHENTICATION`, `NTLM` and `KERBEROS`.",
			Required:    true,
		},
		"credentials": {
			Type:        schema.TypeString,
			Description: "The ID of the credentials within the Dynatrace Credentials Vault.",
			Required:    true,
		},
		"realm_name": {
			Type:        schema.TypeString,
			Description: "The Realm Name. Valid and required only if the type of authentication is `KERBEROS`.",
			Optional:    true,
		},
		"kdc_ip": {
			Type:        schema.TypeString,
			Description: "The KDC IP. Valid and required only if the type of authentication is `KERBEROS`.",
			Optional:    true,
		},
	}
}

// MarshalHCL serializes the fields of an Authentication struct into a map, using the keys specified within the Schema function
func (me *Authentication) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("credentials", me.Credentials); err != nil {
		return err
	}
	if err := properties.Encode("realm_name", me.RealmName); err != nil {
		return err
	}
	if err := properties.Encode("kdc_ip", me.KdcIP); err != nil {
		return err
	}
	return nil
}

// UnmarshalHCL deserializes data available via terraform provider into the Authentication struct. The keys to be used are defined by the Schema function
func (me *Authentication) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("credentials", &me.Credentials); err != nil {
		return err
	}
	if err := decoder.Decode("realm_name", &me.RealmName); err != nil {
		return err
	}
	if err := decoder.Decode("kdc_ip", &me.KdcIP); err != nil {
		return err
	}
	return nil
}
