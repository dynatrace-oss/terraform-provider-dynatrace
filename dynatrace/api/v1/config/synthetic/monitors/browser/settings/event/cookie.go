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

package event

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/synthetic/monitors/request"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Cookie struct {
	EventBase
	Cookies request.Cookies `json:"cookies"` // Every cookie must be unique within the list. However, you can use the same cookie again in other event
}

func (me *Cookie) GetType() Type {
	return Types.Cookie
}

func (me *Cookie) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cookies": {
			Type:        schema.TypeList,
			Description: "Every cookie must be unique within the list. However, you can use the same cookie again in other event",
			Required:    true,
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(request.Cookies).Schema()},
		},
	}
}

func (me *Cookie) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("cookies", me.Cookies); err != nil {
		return err
	}
	return nil
}

func (me *Cookie) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Type = Types.Tap
	if err := decoder.Decode("cookies", &me.Cookies); err != nil {
		return err
	}
	return nil
}
