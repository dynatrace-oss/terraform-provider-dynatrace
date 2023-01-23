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

package apitokens

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type APIToken struct {
	Name                string   `json:"name"`                          // The name of the token.
	Enabled             *bool    `json:"enabled,omitempty"`             // The token is enabled (true) or disabled (false), default disabled (false).
	PersonalAccessToken *bool    `json:"personalAccessToken,omitempty"` // The token is a personal access token (true) or an API token (false).
	ExpirationDate      *string  `json:"expirationDate,omitempty"`      // The expiration date of the token.
	Owner               *string  `json:"owner,omitempty"`               // The owner of the token
	CreationDate        *string  `json:"creationDate,omitempty"`        // Token creation date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z')
	ModifiedDate        *string  `json:"modifiedDate,omitempty"`        // Token last modified date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z').
	LastUsedDate        *string  `json:"lastUsedDate,omitempty"`        // Token last used date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z')
	LastUsedIpAddress   *string  `json:"lastUsedIpAddress,omitempty"`   // Token last used IP address.
	Scopes              []string `json:"scopes"`                        // A list of the scopes to be assigned to the token.
	Token               *string  `json:"token,omitempty"`               // The secret of the token.
}

func (me *APIToken) Equals(v any) (string, bool) {
	if other, ok := v.(*APIToken); ok {
		if me.Name != other.Name {
			return fmt.Sprintf("Name: expected: %s, actual: %s", me.Name, other.Name), false
		}
		if me.Enabled == nil && other.Enabled != nil {
			return fmt.Sprintf("Enabled: expected: %v, actual: %v", me.Enabled, other.Enabled), false
		}
		if me.Enabled != nil && other.Enabled == nil {
			return fmt.Sprintf("Enabled: expected: %v, actual: %v", me.Enabled, other.Enabled), false
		}
		if me.Enabled != nil && *me.Enabled != *other.Enabled {
			return fmt.Sprintf("Enabled: expected: %v, actual: %v", *me.Enabled, *other.Enabled), false
		}
		if me.PersonalAccessToken == nil && other.PersonalAccessToken != nil {
			return fmt.Sprintf("PersonalAccessToken: expected: %v, actual: %v", me.PersonalAccessToken, other.PersonalAccessToken), false
		}
		if me.PersonalAccessToken != nil && other.PersonalAccessToken == nil {
			return fmt.Sprintf("PersonalAccessToken: expected: %v, actual: %v", me.PersonalAccessToken, other.PersonalAccessToken), false
		}
		if me.PersonalAccessToken != nil && *me.PersonalAccessToken != *other.PersonalAccessToken {
			return fmt.Sprintf("PersonalAccessToken: expected: %v, actual: %v", *me.PersonalAccessToken, *other.PersonalAccessToken), false
		}
		if len(me.Scopes) != len(other.Scopes) {
			return fmt.Sprintf("Scopes: expected length: %v, actual: %v", len(me.Scopes), len(other.Scopes)), false
		}
		if len(me.Scopes) > 0 {
			for _, scope := range me.Scopes {
				found := false
				for _, oscope := range other.Scopes {
					if oscope == scope {
						found = true
						break
					}
				}
				if !found {
					return fmt.Sprintf("Scopes: expected: %v, actual: %v", me.Scopes, other.Scopes), false
				}
			}
		}
		return "", true
	}
	return "APIToken expected", false
}

type TokenList struct {
	APITokens []*struct {
		APIToken
		ID *string `json:"id,omitempty"`
	} `json:"apiTokens"` // An ordered list of api tokens
}

func (me *APIToken) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the token.",
			Required:    true,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "The token is enabled (true) or disabled (false), default disabled (false).",
			Optional:    true,
		},
		"personal_access_token": {
			Type:        schema.TypeBool,
			Description: "The token is a personal access token (true) or an API token (false).",
			Optional:    true,
		},
		"expiration_date": {
			Type:        schema.TypeString,
			Description: "The expiration date of the token.",
			Optional:    true,
		},
		"owner": {
			Type:        schema.TypeString,
			Description: "The owner of the token",
			Optional:    true,
			Computed:    true,
		},
		"creation_date": {
			Type:        schema.TypeString,
			Description: "Token creation date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z')",
			Optional:    true,
			Computed:    true,
		},
		"modified_date": {
			Type:        schema.TypeString,
			Description: "Token last modified date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z').",
			Optional:    true,
			Computed:    true,
		},
		"last_used_date": {
			Type:        schema.TypeString,
			Description: "Token last used date in ISO 8601 format (yyyy-MM-dd'T'HH:mm:ss.SSS'Z')",
			Optional:    true,
			Computed:    true,
		},
		"last_used_ip_address": {
			Type:        schema.TypeString,
			Description: "Token last used IP address.",
			Optional:    true,
			Computed:    true,
		},
		"scopes": {
			Type:        schema.TypeSet,
			Description: "A list of the scopes to be assigned to the token.",
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"token": {
			Type:        schema.TypeString,
			Description: "The secret of the token.",
			Sensitive:   true,
			Computed:    true,
		},
	}
}

func (me *APIToken) MarshalHCL(properties hcl.Properties) error {
	if err := properties.EncodeAll(map[string]any{
		"name":                  me.Name,
		"enabled":               me.Enabled,
		"personal_access_token": me.PersonalAccessToken,
		"expiration_date":       me.ExpirationDate,
		"owner":                 me.Owner,
		"creation_date":         me.CreationDate,
		"modified_date":         me.ModifiedDate,
		"last_used_date":        me.LastUsedDate,
		"last_used_ip_address":  me.LastUsedIpAddress,
		"scopes":                me.Scopes,
	}); err != nil {
		return err
	}
	if err := properties.Encode("token", me.Token); err != nil {
		return err
	}
	return nil
}

func (me *APIToken) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"name":                  &me.Name,
		"enabled":               &me.Enabled,
		"personal_access_token": &me.PersonalAccessToken,
		"expiration_date":       &me.ExpirationDate,
		"scopes":                &me.Scopes,
	})
}
