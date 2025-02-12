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

package automaticinjection

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CacheControlHeaders struct {
	CacheControlHeaders bool `json:"cacheControlHeaders"` // [How to ensure timely configuration updates for automatic injection?](https://dt-url.net/m9039ea)
}

func (me *CacheControlHeaders) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cache_control_headers": {
			Type:        schema.TypeBool,
			Description: "[How to ensure timely configuration updates for automatic injection?](https://dt-url.net/m9039ea)",
			Required:    true,
		},
	}
}

func (me *CacheControlHeaders) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"cache_control_headers": me.CacheControlHeaders,
	})
}

func (me *CacheControlHeaders) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"cache_control_headers": &me.CacheControlHeaders,
	})
}
