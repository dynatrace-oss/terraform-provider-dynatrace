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

package logcustomattributes

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	AggregableAttribute bool   `json:"aggregableAttribute"` // Change applies only to newly ingested log events. Any log events ingested before this option was toggled on will not be searchable by this attribute.
	Key                 string `json:"key"`                 // The attribute key is case insensitive in log data ingestion.
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sidebar": {
			Type:        schema.TypeBool,
			Description: "Show attribute values in side bar",
			Required:    true,
		},
		"key": {
			Type:        schema.TypeString,
			Description: "The attribute key is case insensitive in log data ingestion.",
			Required:    true,
			ForceNew:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"sidebar": me.AggregableAttribute,
		"key":     me.Key,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"sidebar": &me.AggregableAttribute,
		"key":     &me.Key,
	})
}
