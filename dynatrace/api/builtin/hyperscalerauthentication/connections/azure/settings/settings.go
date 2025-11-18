/**
* @license
* Copyright 2025 Dynatrace LLC
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

package azure

import (
	"fmt"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	ClientSecret                *ClientSecretConfig          `json:"clientSecret,omitempty"`
	FederatedIdentityCredential *FederatedIdentityCredential `json:"federatedIdentityCredential,omitempty"`
	Name                        string                       `json:"name"` // The name of the connection
	Type                        Type                         `json:"type"` // Azure Authentication mechanism to be used by the connection. Possible Values: `clientSecret`, `federatedIdentityCredential`
}

const DefaultTimeout = 2 * time.Minute

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_secret": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(ClientSecretConfig).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"federated_identity_credential": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(FederatedIdentityCredential).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the connection",
			Required:    true,
			ForceNew:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Azure Authentication mechanism to be used by the connection. Possible Values: `clientSecret`, `federatedIdentityCredential`",
			Required:    true,
		},
	}
}

func (me *Settings) Timeouts() *schema.ResourceTimeout {
	return &schema.ResourceTimeout{
		Create: schema.DefaultTimeout(DefaultTimeout),
		Update: schema.DefaultTimeout(DefaultTimeout),
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"client_secret":                 me.ClientSecret,
		"federated_identity_credential": me.FederatedIdentityCredential,
		"name":                          me.Name,
		"type":                          me.Type,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.ClientSecret == nil) && (string(me.Type) == "clientSecret") {
		return fmt.Errorf("'client_secret' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.ClientSecret != nil) && (string(me.Type) != "clientSecret") {
		return fmt.Errorf("'client_secret' must not be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FederatedIdentityCredential == nil) && (string(me.Type) == "federatedIdentityCredential") {
		return fmt.Errorf("'federated_identity_credential' must be specified if 'type' is set to '%v'", me.Type)
	}
	if (me.FederatedIdentityCredential != nil) && (string(me.Type) != "federatedIdentityCredential") {
		return fmt.Errorf("'federated_identity_credential' must not be specified if 'type' is set to '%v'", me.Type)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"client_secret":                 &me.ClientSecret,
		"federated_identity_credential": &me.FederatedIdentityCredential,
		"name":                          &me.Name,
		"type":                          &me.Type,
	})
}
