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

package diskrules

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"
)

type Settings struct {
	DiskNameFilter        *DiskNameFilter `json:"diskNameFilter"`                  // Only apply to disks whose name matches
	Enabled               bool            `json:"enabled"`                         // This setting is enabled (`true`) or disabled (`false`)
	HostGroupID           *string         `json:"-" scope:"hostGroupId"`           // The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.
	Metric                DiskMetric      `json:"metric"`                          // Possible Values: `LOW_DISK_SPACE`, `LOW_INODES`, `READ_TIME_EXCEEDING`, `WRITE_TIME_EXCEEDING`
	Name                  string          `json:"name"`                            // Name
	SampleLimit           *SampleLimit    `json:"sampleLimit"`                     // Only alert if the threshold was violated in at least *n* of the last *m* samples
	TagFilters            []string        `json:"tagFilters,omitempty"`            // Only apply to hosts that have the following tags
	ThresholdMilliseconds *float64        `json:"thresholdMilliseconds,omitempty"` // Alert if higher than
	ThresholdPercent      *float64        `json:"thresholdPercent,omitempty"`      // Alert if lower than
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"disk_name_filter": {
			Type:        schema.TypeList,
			Description: "Only apply to disks whose name matches",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(DiskNameFilter).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"enabled": {
			Type:        schema.TypeBool,
			Description: "This setting is enabled (`true`) or disabled (`false`)",
			Required:    true,
		},
		"host_group_id": {
			Type:        schema.TypeString,
			Description: "The scope of this settings. If the settings should cover the whole environment, just don't specify any scope.",
			Optional:    true,
			Default:     "environment",
		},
		"metric": {
			Type:        schema.TypeString,
			Description: "Possible Values: `LOW_DISK_SPACE`, `LOW_INODES`, `READ_TIME_EXCEEDING`, `WRITE_TIME_EXCEEDING`",
			Required:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name",
			Required:    true,
		},
		"sample_limit": {
			Type:        schema.TypeList,
			Description: "Only alert if the threshold was violated in at least *n* of the last *m* samples",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(SampleLimit).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"tag_filters": {
			Type:        schema.TypeSet,
			Description: "Only apply to hosts that have the following tags",
			Optional:    true, // minobjects == 0
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"threshold_milliseconds": {
			Type:        schema.TypeFloat,
			Description: "Alert if higher than",
			Optional:    true, // precondition
		},
		"threshold_percent": {
			Type:        schema.TypeFloat,
			Description: "Alert if lower than",
			Optional:    true, // precondition
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"disk_name_filter":       me.DiskNameFilter,
		"enabled":                me.Enabled,
		"host_group_id":          me.HostGroupID,
		"metric":                 me.Metric,
		"name":                   me.Name,
		"sample_limit":           me.SampleLimit,
		"tag_filters":            me.TagFilters,
		"threshold_milliseconds": me.ThresholdMilliseconds,
		"threshold_percent":      me.ThresholdPercent,
	})
}

func (me *Settings) HandlePreconditions() error {
	if me.ThresholdMilliseconds == nil && slices.Contains([]string{"READ_TIME_EXCEEDING", "WRITE_TIME_EXCEEDING"}, string(me.Metric)) {
		me.ThresholdMilliseconds = opt.NewFloat64(0.0)
	}
	if me.ThresholdPercent == nil && slices.Contains([]string{"LOW_DISK_SPACE", "LOW_INODES"}, string(me.Metric)) {
		me.ThresholdPercent = opt.NewFloat64(0.0)
	}
	return nil
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"disk_name_filter":       &me.DiskNameFilter,
		"enabled":                &me.Enabled,
		"host_group_id":          &me.HostGroupID,
		"metric":                 &me.Metric,
		"name":                   &me.Name,
		"sample_limit":           &me.SampleLimit,
		"tag_filters":            &me.TagFilters,
		"threshold_milliseconds": &me.ThresholdMilliseconds,
		"threshold_percent":      &me.ThresholdPercent,
	})
}
