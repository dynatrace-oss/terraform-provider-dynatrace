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

package options

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DiskComplexes []*DiskComplex

func (me *DiskComplexes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"exclusion": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DiskComplex).Schema()},
		},
	}
}

func (me DiskComplexes) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("exclusion", me)
}

func (me *DiskComplexes) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("exclusion", me)
}

type DiskComplex struct {
	Filesystem *string    `json:"filesystem,omitempty"` // **File system type field:** the type of the file system to be excluded from monitoring. Examples:\n\n* ext4\n* ext3\n* btrfs\n* ext*\n\n⚠️ Starting from **OneAgent 1.299+** file system types are not case sensitive! \n\nThe wildcard in the last example means to exclude matching file systems such as types ext4 and ext3
	Mountpoint *string    `json:"mountpoint,omitempty"` // **Disk or mount point path field:** the path to where the disk to be excluded from monitoring is mounted. Examples:\n\n* /mnt/my_disk\n* /staff/emp1\n* C:\\\n* /staff/*\n* /disk*\n\n ⚠️ Mount point paths are case sensitive! \n\nThe wildcard in **/staff/*** means to exclude every child folder of /staff.\n\nThe wildcard in **/disk*** means to exclude every mount point starting with /disk, for example /disk1, /disk99,  /diskabc
	Os         OsTypeEnum `json:"os"`                   // Possible Values: `OS_TYPE_AIX`, `OS_TYPE_DARWIN`, `OS_TYPE_HPUX`, `OS_TYPE_LINUX`, `OS_TYPE_SOLARIS`, `OS_TYPE_UNKNOWN`, `OS_TYPE_WINDOWS`, `OS_TYPE_ZOS`
}

func (me *DiskComplex) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filesystem": {
			Type:        schema.TypeString,
			Description: "**File system type field:** the type of the file system to be excluded from monitoring. Examples:\n\n* ext4\n* ext3\n* btrfs\n* ext*\n\n⚠️ Starting from **OneAgent 1.299+** file system types are not case sensitive! \n\nThe wildcard in the last example means to exclude matching file systems such as types ext4 and ext3",
			Optional:    true,
		},
		"mountpoint": {
			Type:        schema.TypeString,
			Description: "**Disk or mount point path field:** the path to where the disk to be excluded from monitoring is mounted. Examples:\n\n* /mnt/my_disk\n* /staff/emp1\n* C:\\\n* /staff/*\n* /disk*\n\n ⚠️ Mount point paths are case sensitive! \n\nThe wildcard in **/staff/*** means to exclude every child folder of /staff.\n\nThe wildcard in **/disk*** means to exclude every mount point starting with /disk, for example /disk1, /disk99,  /diskabc",
			Optional:    true,
		},
		"os": {
			Type:        schema.TypeString,
			Description: "Possible Values: `OS_TYPE_AIX`, `OS_TYPE_DARWIN`, `OS_TYPE_HPUX`, `OS_TYPE_LINUX`, `OS_TYPE_SOLARIS`, `OS_TYPE_UNKNOWN`, `OS_TYPE_WINDOWS`, `OS_TYPE_ZOS`",
			Required:    true,
		},
	}
}

func (me *DiskComplex) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"filesystem": me.Filesystem,
		"mountpoint": me.Mountpoint,
		"os":         me.Os,
	})
}

func (me *DiskComplex) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"filesystem": &me.Filesystem,
		"mountpoint": &me.Mountpoint,
		"os":         &me.Os,
	})
}
