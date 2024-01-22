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
	EnableCrashDumpAnalytics bool   `json:"enableCrashDumpAnalytics"` // Control the automatic crash dump analytics feature. It's strongly recommended to keep this option enabled. Disabling this feature will stop receiving information about potential problems
	HostID                   string `json:"-" scope:"hostId"`         // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
}

func (me *Settings) Name() string {
	return me.HostID
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_crash_dump_analytics": {
			Type:        schema.TypeBool,
			Description: "Control the automatic crash dump analytics feature. It's strongly recommended to keep this option enabled. Disabling this feature will stop receiving information about potential problems",
			Required:    true,
		},
		"host_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
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
