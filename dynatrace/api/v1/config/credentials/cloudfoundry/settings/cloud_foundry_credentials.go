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

package cloudfoundry

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CloudFoundryCredentials Configuration for specific Cloud Foundry credentials.
type CloudFoundryCredentials struct {
	LoginURL string                     `json:"loginUrl"`           // The login URL of the Cloud Foundry foundation credentials.  The URL must be valid according to RFC 2396.  Leading or trailing whitespaces are not allowed.
	Password *string                    `json:"password,omitempty"` // The password of the Cloud Foundry foundation credentials.
	Active   bool                       `json:"active"`             // The monitoring is enabled (`true`) or disabled (`false`) for given credentials configuration.  If not set on creation, the `true` value is used.  If the field is omitted during an update, the old value remains unaffected.
	Name     string                     `json:"name"`               // The name of the Cloud Foundry foundation credentials.  Allowed characters are letters, numbers, whitespaces, and the following characters: `.+-_`. Leading or trailing whitespace is not allowed.
	Username string                     `json:"username"`           // The username of the Cloud Foundry foundation credentials.  Leading and trailing whitespaces are not allowed.
	APIURL   string                     `json:"apiUrl"`             // The URL of the Cloud Foundry foundation credentials.  The URL must be valid according to RFC 2396.  Leading or trailing whitespaces are not allowed.
	Unknowns map[string]json.RawMessage `json:"-"`
}

func (me *CloudFoundryCredentials) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"login_url": {
			Type:        schema.TypeString,
			Description: "The login URL of the Cloud Foundry foundation credentials. The URL must be valid according to RFC 2396.  Leading or trailing whitespaces are not allowed.",
			Required:    true,
		},
		"api_url": {
			Type:        schema.TypeString,
			Description: "The URL of the Cloud Foundry foundation credentials.  The URL must be valid according to RFC 2396.  Leading or trailing whitespaces are not allowed.",
			Required:    true,
		},
		"username": {
			Type:        schema.TypeString,
			Description: "The username of the Cloud Foundry foundation credentials.  Leading and trailing whitespaces are not allowed.",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the Cloud Foundry foundation credentials.  Allowed characters are letters, numbers, whitespaces, and the following characters: `.+-_`. Leading or trailing whitespace is not allowed.",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "The password of the Cloud Foundry foundation credentials.",
			Optional:    true,
			Sensitive:   true,
		},
		"active": {
			Type:        schema.TypeBool,
			Description: "The monitoring is enabled (`true`) or disabled (`false`) for given credentials configuration.  If not set on creation, the `true` value is used.  If the field is omitted during an update, the old value remains unaffected.",
			Optional:    true,
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (me *CloudFoundryCredentials) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.EncodeAll(map[string]any{
		"login_url": me.LoginURL,
		"api_url":   me.APIURL,
		"password":  "${state.secret_value}",
		"active":    me.Active,
		"name":      me.Name,
		"username":  me.Username,
		"unknowns":  me.Unknowns,
	})
}

func (me *CloudFoundryCredentials) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"login_url": &me.LoginURL,
		"api_url":   &me.APIURL,
		"password":  &me.Password,
		"active":    &me.Active,
		"name":      &me.Name,
		"username":  &me.Username,
		"unknowns":  &me.Unknowns,
	})
	return err
}

func (me *CloudFoundryCredentials) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"loginUrl": me.LoginURL,
		"apiUrl":   me.APIURL,
		"password": me.Password,
		"active":   me.Active,
		"name":     me.Name,
		"username": me.Username,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *CloudFoundryCredentials) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	err := properties.UnmarshalAll(map[string]any{
		"loginUrl": &me.LoginURL,
		"apiUrl":   &me.APIURL,
		"password": &me.Password,
		"active":   &me.Active,
		"name":     &me.Name,
		"username": &me.Username,
	})
	return err
}

const credsNotProvided = "REST API didn't provide credential data"

func (me *CloudFoundryCredentials) FillDemoValues() []string {
	me.Password = opt.NewString("################")
	return []string{credsNotProvided}
}
