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

package analytics

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	EnableCrashDumpAnalytics bool   `json:"enableCrashDumpAnalytics"` // Disable the feature to stop receiving information about crash details and potential problems. We recommend keeping the feature enabled.
	HostID                   string `json:"-" scope:"hostId"`         // The scope of this setting (HOST HOST_GROUP environment)
}

func (me *Settings) Name() string {
	return me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_crash_dump_analytics": {
			Type:        schema.TypeBool,
			Description: "Disable the feature to stop receiving information about crash details and potential problems. We recommend keeping the feature enabled.",
			Required:    true,
		},
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST HOST_GROUP environment)",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"enable_crash_dump_analytics": me.EnableCrashDumpAnalytics,
		"host_id":                     me.HostID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"enable_crash_dump_analytics": &me.EnableCrashDumpAnalytics,
		"host_id":                     &me.HostID,
	})
}
