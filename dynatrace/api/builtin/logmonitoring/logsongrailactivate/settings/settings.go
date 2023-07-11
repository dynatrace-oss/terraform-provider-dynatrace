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

package logsongrailactivate

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	Activated            bool                 `json:"activated"`            // Activate logs powered by Grail.
	ParallelIngestPeriod ParallelIngestPeriod `json:"parallelIngestPeriod"` // Possible Values: `NONE`, `SEVEN_DAYS`, `THIRTY_FIVE_DAYS`
}

func (me *Settings) Name() string {
	return "log_grail"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"activated": {
			Type:        schema.TypeBool,
			Description: "Activate logs powered by Grail.",
			Required:    true,
		},
		"parallel_ingest_period": {
			Type:        schema.TypeString,
			Description: "Possible Values: `NONE`, `SEVEN_DAYS`, `THIRTY_FIVE_DAYS`",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"activated":              me.Activated,
		"parallel_ingest_period": me.ParallelIngestPeriod,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"activated":              &me.Activated,
		"parallel_ingest_period": &me.ParallelIngestPeriod,
	})
}
