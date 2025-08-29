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

package ingestsources

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AzureLogForwardingAttributes struct {
	FieldExtraction   *FieldExtraction `json:"fieldExtraction"` // Field Extraction
	ForwarderConfigID string           `json:"forwarderConfigId"`
}

func (me *AzureLogForwardingAttributes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field_extraction": {
			Type:        schema.TypeList,
			Description: "Field Extraction",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(FieldExtraction).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"forwarder_config_id": {
			Type:        schema.TypeString,
			Description: "no documentation available",
			Required:    true,
		},
	}
}

func (me *AzureLogForwardingAttributes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"field_extraction":    me.FieldExtraction,
		"forwarder_config_id": me.ForwarderConfigID,
	})
}

func (me *AzureLogForwardingAttributes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"field_extraction":    &me.FieldExtraction,
		"forwarder_config_id": &me.ForwarderConfigID,
	})
}
