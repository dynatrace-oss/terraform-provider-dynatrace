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

package settings

import (
	"fmt"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type Users []*UserAccessor

func (_ *Users) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "User that is to be granted read or write permissions",
			Elem:        &schema.Resource{Schema: new(UserAccessor).Schema()},
		},
	}
}

func (u *Users) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("user", u)
}

func (u *Users) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("user", u)
}

type UserAccessor struct {
	UID    string
	Access string
}

func (_ *UserAccessor) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uid": {
			Type:        schema.TypeString,
			Description: "The UUID of the user, conveniently retrieved via the `uid` attribute provided by the data source `dynatrace_iam_user`",
			Required:    true,
		},
		"access": {
			Type:         schema.TypeString,
			Description:  fmt.Sprintf("Valid values: `%s`, `%s`", HCLAccessorRead, HCLAccessorWrite),
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{HCLAccessorRead, HCLAccessorWrite}, false),
		},
	}
}

func (ua *UserAccessor) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"uid":    ua.UID,
		"access": ua.Access,
	})
}

func (ua *UserAccessor) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"uid":    &ua.UID,
		"access": &ua.Access,
	})
}
