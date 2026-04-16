/**
* @license
* Copyright 2026 Dynatrace LLC
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

package anomalydetectors

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DetectionConditions []*DetectionCondition

func (me *DetectionConditions) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"detection_condition": {
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(DetectionCondition).Schema()},
		},
	}
}

func (me DetectionConditions) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("detection_condition", me)
}

func (me *DetectionConditions) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("detection_condition", me)
}

type DetectionCondition struct {
	DiskFilesystemCondition *string                    `json:"diskFilesystemCondition,omitempty"` // Disk filesystem will be included in this policy if **any** of the filters match. Disk filesystem has to match a required format.\n\n  - `$match(ext*)` – Matches string with wildcards: `*` any number (including zero) of characters and `?` exactly one character.\n - `$contains(fs)` – Matches if `fs` appears anywhere in the filesystem type.\n - `$eq(ext4)` – Matches if `ext4` matches the filesystem type exactly.\n - `$prefix(ext)` – Matches if `ext` matches the prefix of the filesystem type.\n - `$suffix(fs)` – Matches if `fs` matches the suffix of the filesystem type.\n\n  Available logic operations:\n - `$not($eq(tmpfs))` – Matches if the filesystem type is different from `tmpfs`.\n - `$and($prefix(ext),$suffix(4))` – Matches if filesystem type starts with `ext` and ends with `4`.\n - `$or($eq(xfs),$eq(btrfs))` – Matches if filesystem type equals `xfs` or `btrfs`.\n\n  Brackets **(** and **)** that are part of the matched filesystem type **must be escaped with a tilde (~)**
	DiskTotalCondition      *DiskTotalSpaceThresholds  `json:"diskTotalCondition,omitempty"`      // Specify disk total space range in GiB
	HostMetadataCondition   *HostMetadataConditionType `json:"hostMetadataCondition,omitempty"`   // Host resource attributes are dimensions enriching the host including custom metadata which are user-defined key-value pairs that you can assign to hosts monitored by Dynatrace.\n\n  By defining custom metadata, you can enrich the monitoring data with context specific to your organization's needs, such as environment names, team ownership, application versions, or any other relevant details.\n\n  See [Define tags and metadata for hosts](https://dt-url.net/w3hv0kbw).\n\n  Note: Starting from version 1.325 host resource attributes are supported in addition to host custom metadata.
	LocalDiskCondition      *DiskType                  `json:"localDiskCondition,omitempty"`      // Possible values: `LOCAL`, `REMOTE`
	Property                *DiskProp                  `json:"property,omitempty"`                // Disk property. Possible values: `DiskFilesystem`, `DiskTotalSpace`, `DiskType`
	RuleType                RuleType                   `json:"ruleType"`                          // Starting from agent 1.335 **disk** detection rules are supported. Possible values: `RuleTypeDisk`, `RuleTypeHost`
}

func (me *DetectionCondition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk_filesystem_condition": {
			Type:        schema.TypeString,
			Description: "Disk filesystem will be included in this policy if **any** of the filters match. Disk filesystem has to match a required format.\n\n  - `$match(ext*)` – Matches string with wildcards: `*` any number (including zero) of characters and `?` exactly one character.\n - `$contains(fs)` – Matches if `fs` appears anywhere in the filesystem type.\n - `$eq(ext4)` – Matches if `ext4` matches the filesystem type exactly.\n - `$prefix(ext)` – Matches if `ext` matches the prefix of the filesystem type.\n - `$suffix(fs)` – Matches if `fs` matches the suffix of the filesystem type.\n\n  Available logic operations:\n - `$not($eq(tmpfs))` – Matches if the filesystem type is different from `tmpfs`.\n - `$and($prefix(ext),$suffix(4))` – Matches if filesystem type starts with `ext` and ends with `4`.\n - `$or($eq(xfs),$eq(btrfs))` – Matches if filesystem type equals `xfs` or `btrfs`.\n\n  Brackets **(** and **)** that are part of the matched filesystem type **must be escaped with a tilde (~)**",
			Optional:    true, // precondition
		},
		"disk_total_condition": {
			Type:        schema.TypeList,
			Description: "Specify disk total space range in GiB",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(DiskTotalSpaceThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"host_metadata_condition": {
			Type:        schema.TypeList,
			Description: "Host resource attributes are dimensions enriching the host including custom metadata which are user-defined key-value pairs that you can assign to hosts monitored by Dynatrace.\n\n  By defining custom metadata, you can enrich the monitoring data with context specific to your organization's needs, such as environment names, team ownership, application versions, or any other relevant details.\n\n  See [Define tags and metadata for hosts](https://dt-url.net/w3hv0kbw).\n\n  Note: Starting from version 1.325 host resource attributes are supported in addition to host custom metadata.",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(HostMetadataConditionType).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"local_disk_condition": {
			Type:        schema.TypeString,
			Description: "Possible values: `LOCAL`, `REMOTE`",
			Optional:    true, // precondition
		},
		"property": {
			Type:        schema.TypeString,
			Description: "Disk property. Possible values: `DiskFilesystem`, `DiskTotalSpace`, `DiskType`",
			Optional:    true, // precondition
		},
		"rule_type": {
			Type:        schema.TypeString,
			Description: "Starting from agent 1.335 **disk** detection rules are supported. Possible values: `RuleTypeDisk`, `RuleTypeHost`",
			Required:    true,
		},
	}
}

func (me *DetectionCondition) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"disk_filesystem_condition": me.DiskFilesystemCondition,
		"disk_total_condition":      me.DiskTotalCondition,
		"host_metadata_condition":   me.HostMetadataCondition,
		"local_disk_condition":      me.LocalDiskCondition,
		"property":                  me.Property,
		"rule_type":                 me.RuleType,
	})
}

func (me *DetectionCondition) HandlePreconditions() error {
	if (me.DiskFilesystemCondition == nil) && ((string(me.RuleType) == "RuleTypeDisk") && (me.Property != nil && (string(*me.Property) == "DiskFilesystem"))) {
		return fmt.Errorf("'disk_filesystem_condition' must be specified if 'rule_type' is set to '%v' and 'property' is set to '%v'", me.RuleType, me.Property)
	}
	if (me.DiskFilesystemCondition != nil) && ((string(me.RuleType) != "RuleTypeDisk") || (me.Property == nil || (me.Property != nil && string(*me.Property) != "DiskFilesystem"))) {
		return fmt.Errorf("'disk_filesystem_condition' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.DiskTotalCondition == nil) && ((string(me.RuleType) == "RuleTypeDisk") && (me.Property != nil && (string(*me.Property) == "DiskTotalSpace"))) {
		return fmt.Errorf("'disk_total_condition' must be specified if 'rule_type' is set to '%v' and 'property' is set to '%v'", me.RuleType, me.Property)
	}
	if (me.DiskTotalCondition != nil) && ((string(me.RuleType) != "RuleTypeDisk") || (me.Property == nil || (me.Property != nil && string(*me.Property) != "DiskTotalSpace"))) {
		return fmt.Errorf("'disk_total_condition' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.HostMetadataCondition == nil) && (slices.Contains([]string{"RuleTypeHost"}, string(me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.HostMetadataCondition != nil) && (!slices.Contains([]string{"RuleTypeHost"}, string(me.RuleType))) {
		return fmt.Errorf("'host_metadata_condition' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.LocalDiskCondition == nil) && ((string(me.RuleType) == "RuleTypeDisk") && (me.Property != nil && (string(*me.Property) == "DiskType"))) {
		return fmt.Errorf("'local_disk_condition' must be specified if 'rule_type' is set to '%v' and 'property' is set to '%v'", me.RuleType, me.Property)
	}
	if (me.LocalDiskCondition != nil) && ((string(me.RuleType) != "RuleTypeDisk") || (me.Property == nil || (me.Property != nil && string(*me.Property) != "DiskType"))) {
		return fmt.Errorf("'local_disk_condition' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.Property == nil) && (string(me.RuleType) == "RuleTypeDisk") {
		return fmt.Errorf("'property' must be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	if (me.Property != nil) && (string(me.RuleType) != "RuleTypeDisk") {
		return fmt.Errorf("'property' must not be specified if 'rule_type' is set to '%v'", me.RuleType)
	}
	return nil
}

func (me *DetectionCondition) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"disk_filesystem_condition": &me.DiskFilesystemCondition,
		"disk_total_condition":      &me.DiskTotalCondition,
		"host_metadata_condition":   &me.HostMetadataCondition,
		"local_disk_condition":      &me.LocalDiskCondition,
		"property":                  &me.Property,
		"rule_type":                 &me.RuleType,
	})
}
