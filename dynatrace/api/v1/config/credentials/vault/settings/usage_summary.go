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

package vault

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type UsageSummary []*CredentialUsageObj

type CredentialUsageObj struct {
	MonitorType MonitorType `json:"type"`
	Count       int32       `json:"count"`
}

func (me *CredentialUsageObj) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "Type of usage, `HTTP_MONITOR` or `BROWSER_MONITOR`",
			Required:    true,
		},
		"count": {
			Type:        schema.TypeInt,
			Description: "The number of uses",
			Required:    true,
		},
	}
}

func (me *CredentialUsageObj) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", string(me.MonitorType)); err != nil {
		return err
	}
	if err := properties.Encode("count", int(me.Count)); err != nil {
		return err
	}
	return nil
}

func (me *CredentialUsageObj) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("type"); ok {
		me.MonitorType = MonitorType(value.(string))
	}
	if value, ok := decoder.GetOk("count"); ok {
		me.Count = int32(value.(int))
	}
	return nil
}
