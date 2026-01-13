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

package options

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	DisableNfsDiskMonitoring bool          `json:"disableNfsDiskMonitoring"` // Deactivate NFS monitoring on all supported systems
	Exclusions               DiskComplexes `json:"exclusions,omitempty"`     // OneAgent automatically detects and monitors all your mount points, however you can create exception rules to remove disks from the monitoring list.
	MonitorTmpfs             bool          `json:"monitorTmpfs"`             // Activate tmpfs monitoring on Linux systems
	NfsShowAll               *bool         `json:"nfsShowAll,omitempty"`     // When disabled OneAgent will try to deduplicate some of nfs mount points. Disabled by default, applies only to Linux hosts.
	Scope                    *string       `json:"-" scope:"scope"`          // The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.
}

func (me *Settings) Name() string {
	return *me.Scope
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disable_nfs_disk_monitoring": {
			Type:        schema.TypeBool,
			Description: "Deactivate NFS monitoring on all supported systems",
			Optional:    true,
		},
		"exclusions": {
			Type:        schema.TypeList,
			Description: "OneAgent automatically detects and monitors all your mount points, however you can create exception rules to remove disks from the monitoring list.",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Resource{Schema: new(DiskComplexes).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"monitor_tmpfs": {
			Type:        schema.TypeBool,
			Description: "Activate tmpfs monitoring on Linux systems",
			Optional:    true, // new required property. Default to false (default bool value)
		},
		"nfs_show_all": {
			Type:        schema.TypeBool,
			Description: "When disabled OneAgent will try to deduplicate some of nfs mount points. Disabled by default, applies only to Linux hosts.",
			Optional:    true, // precondition
		},
		"scope": {
			Type:        schema.TypeString,
			Description: "The scope of this setting (HOST, HOST_GROUP). Omit this property if you want to cover the whole environment.",
			Optional:    true,
			Default:     "environment",
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"disable_nfs_disk_monitoring": me.DisableNfsDiskMonitoring,
		"exclusions":                  me.Exclusions,
		"monitor_tmpfs":               me.MonitorTmpfs,
		"nfs_show_all":                me.NfsShowAll,
		"scope":                       me.Scope,
	})
}

func (me *Settings) HandlePreconditions() error {
	if (me.NfsShowAll == nil) && (!me.DisableNfsDiskMonitoring) {
		me.NfsShowAll = opt.NewBool(false)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"disable_nfs_disk_monitoring": &me.DisableNfsDiskMonitoring,
		"exclusions":                  &me.Exclusions,
		"monitor_tmpfs":               &me.MonitorTmpfs,
		"nfs_show_all":                &me.NfsShowAll,
		"scope":                       &me.Scope,
	})
}
