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

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Headers []*Header

func (me *Headers) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"header": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "An additional HTTP Header to include when sending requests",
			Elem:        &schema.Resource{Schema: new(Header).Schema()},
		},
	}
}

func (me Headers) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("header", me)
}

func (me *Headers) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("header", me)
}

type Header struct {
	Name        string  `json:"name"`                  // The name of the HTTP header
	Secret      bool    `json:"secret"`                // The value of this HTTP header is a secret (`true`) or not (`false`).
	Value       *string `json:"value,omitempty"`       // The value of the HTTP header. May contain an empty value
	SecretValue *string `json:"secretValue,omitempty"` // The secret value of the HTTP header. May contain an empty value
}

func (me *Header) Equals(v any) bool {
	if v == nil {
		return false
	}
	if other, ok := v.(*Header); ok {
		if me.Name != other.Name {
			return false
		}
		if me.Secret != other.Secret {
			return false
		}
		if me.Value == nil && other.Value != nil {
			return false
		}
		if me.Value != nil && other.Value == nil {
			return false
		}
		if me.Value != nil && *me.Value != *other.Value {
			return false
		}
		if me.SecretValue == nil && other.SecretValue != nil {
			return false
		}
		if me.SecretValue != nil && other.SecretValue == nil {
			return false
		}
		if me.SecretValue != nil && *me.SecretValue != *other.SecretValue {
			return false
		}
		return true
	}
	return false
}

func (me *Header) FillDemoValues() []string {
	if me.Secret {
		me.SecretValue = opt.NewString("#######")
		return []string{"Please fill in the secret header value"}
	}
	return []string{}
}

func (me *Header) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the HTTP header",
			Required:    true,
		},
		"secret_value": {
			Type:        schema.TypeString,
			Description: "The value of the HTTP header as a sensitive property. May contain an empty value. `secret_value` and `value` are mutually exclusive. Only one of those two is allowed to be specified.",
			Sensitive:   true,
			Optional:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the HTTP header. May contain an empty value. `secret_value` and `value` are mutually exclusive. Only one of those two is allowed to be specified.",
			Optional:    true,
		},
	}
}

func (me *Header) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if me.Secret {
		if err := properties.Encode("secret_value", "${state.secret_value}"); err != nil {
			return err
		}
		if err := properties.Encode("value", ""); err != nil {
			return err
		}
	} else {
		if err := properties.Encode("secret_value", nil); err != nil {
			return err
		}
		if err := properties.Encode("value", me.Value); err != nil {
			return err
		}
	}

	return nil
}

func (me *Header) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = opt.NewString(value.(string))
		me.Secret = false
	}
	if value, ok := decoder.GetOk("secret_value"); ok {
		me.SecretValue = opt.NewString(value.(string))
		me.Secret = true
	}
	if me.Secret {
		me.Value = nil
	} else {
		me.SecretValue = nil
	}
	return nil
}
