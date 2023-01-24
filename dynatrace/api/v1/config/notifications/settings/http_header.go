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

package notifications

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HTTPHeader The HTTP header.
type HTTPHeader struct {
	Name  string  `json:"name"`            // The name of the HTTP header.
	Value *string `json:"value,omitempty"` // The value of the HTTP header. May contain an empty value.   Required when creating a new notification.  For the **Authorization** header, GET requests return the `null` value.  If you want update a notification configuration with an **Authorization** header which you want to remain intact, set the **Authorization** header with the `null` value.
}

func (me *HTTPHeader) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the HTTP header",
			Required:    true,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "The value of the HTTP header. May contain an empty value.   Required when creating a new notification.  For the **Authorization** header, GET requests return the `null` value.  If you want update a notification configuration with an **Authorization** header which you want to remain intact, set the **Authorization** header with the `null` value",
			Optional:    true,
		},
	}
}

func (me *HTTPHeader) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("value", me.Value); err != nil {
		return err
	}

	return nil
}

func (me *HTTPHeader) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("value"); ok {
		me.Value = opt.NewString(value.(string))
	}
	return nil
}

func (me *HTTPHeader) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"name":  me.Name,
		"value": me.Value,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *HTTPHeader) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"name":  &me.Name,
		"value": &me.Value,
	}); err != nil {
		return err
	}
	return nil
}
