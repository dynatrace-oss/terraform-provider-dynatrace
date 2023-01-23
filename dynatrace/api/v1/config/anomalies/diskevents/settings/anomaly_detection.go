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

package diskevents

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/common"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AnomalyDetection has no documentation
type AnomalyDetection struct {
	Name             string          `json:"name"`                     // The name of the disk event rule.
	HostGroupID      *string         `json:"hostGroupId,omitempty"`    // Narrows the rule usage down to disks that run on hosts that themselves run on the specified host group.
	Threshold        float64         `json:"threshold"`                // The threshold to trigger disk event.   * A percentage for `LowDiskSpace` or `LowInodes` metrics.   * In milliseconds for `ReadTimeExceeding` or `WriteTimeExceeding` metrics.
	DiskNameFilter   *DiskNameFilter `json:"diskNameFilter,omitempty"` // Narrows the rule usage down to disks, matching the specified criteria.
	Enabled          bool            `json:"enabled"`                  // Disk event rule enabled/disabled.
	Samples          int32           `json:"samples"`                  // The number of samples to evaluate.
	ViolatingSamples int32           `json:"violatingSamples"`         // The number of samples that must violate the threshold to trigger an event. Must not exceed the number of evaluated samples.
	Metric           Metric          `json:"metric"`                   // The metric to monitor.
	TagFilters       TagFilters      `json:"tagFilters,omitempty"`     // Narrows the rule usage down to the hosts matching the specified tags.
}

func (me *AnomalyDetection) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the disk event rule",
		},
		"host_group_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Narrows the rule usage down to disks that run on hosts that themselves run on the specified host group",
		},
		"threshold": {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "The threshold to trigger disk event.   * A percentage for `LowDiskSpace` or `LowInodes` metrics.   * In milliseconds for `ReadTimeExceeding` or `WriteTimeExceeding` metrics",
		},
		"disk_name": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Narrows the rule usage down to disks, matching the specified criteria",
			Elem:        &schema.Resource{Schema: new(DiskNameFilter).Schema()},
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Disk event rule enabled/disabled",
		},
		"samples": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of samples to evaluate",
		},
		"violating_samples": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "The number of samples that must violate the threshold to trigger an event. Must not exceed the number of evaluated samples",
		},
		"metric": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The metric to monitor. Possible values are: `LOW_DISK_SPACE`, `LOW_INODES`, `READ_TIME_EXCEEDING` and `WRITE_TIME_EXCEEDING`",
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Narrows the rule usage down to the hosts matching the specified tags",
			Elem:        &schema.Resource{Schema: new(common.TagFilters).Schema()},
		},
	}
}

func (me *AnomalyDetection) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"name":              me.Name,
		"threshold":         me.Threshold,
		"enabled":           me.Enabled,
		"violating_samples": me.ViolatingSamples,
		"samples":           me.Samples,
		"metric":            me.Metric,
		"host_group_id":     me.HostGroupID,
		"disk_name":         me.DiskNameFilter,
		"tags":              me.TagFilters,
	})
}

func (me *AnomalyDetection) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("name"); ok {
		me.Name = value.(string)
	}
	if value, ok := decoder.GetOk("host_group_id"); ok {
		me.HostGroupID = opt.NewString(value.(string))
	}
	if me.HostGroupID != nil && len(*me.HostGroupID) == 0 {
		me.HostGroupID = nil
	}
	if value, ok := decoder.GetOk("threshold"); ok {
		me.Threshold = value.(float64)
	}
	if value, ok := decoder.GetOk("enabled"); ok {
		me.Enabled = value.(bool)
	}
	if value, ok := decoder.GetOk("violating_samples"); ok {
		me.ViolatingSamples = int32(value.(int))
	}
	if value, ok := decoder.GetOk("samples"); ok {
		me.Samples = int32(value.(int))
	}
	if value, ok := decoder.GetOk("metric"); ok {
		me.Metric = Metric(value.(string))
	}

	if _, ok := decoder.GetOk("disk_name.#"); ok {
		me.DiskNameFilter = new(DiskNameFilter)
		if err := me.DiskNameFilter.UnmarshalHCL(hcl.NewDecoder(decoder, "disk_name", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("tags.#"); ok {
		me.TagFilters = TagFilters{}
		if err := me.TagFilters.UnmarshalHCL(hcl.NewDecoder(decoder, "tags", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (me *AnomalyDetection) MarshalJSON() ([]byte, error) {
	properties := xjson.Properties{}
	if err := properties.MarshalAll(map[string]any{
		"name":             me.Name,
		"hostGroupId":      me.HostGroupID,
		"threshold":        me.Threshold,
		"enabled":          me.Enabled,
		"samples":          me.Samples,
		"violatingSamples": me.ViolatingSamples,
		"metric":           me.Metric,
		"diskNameFilter":   me.DiskNameFilter,
		"tagFilters":       me.TagFilters,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *AnomalyDetection) UnmarshalJSON(data []byte) error {
	properties := xjson.Properties{}
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"name":             &me.Name,
		"hostGroupId":      &me.HostGroupID,
		"threshold":        &me.Threshold,
		"enabled":          &me.Enabled,
		"samples":          &me.Samples,
		"violatingSamples": &me.ViolatingSamples,
		"metric":           &me.Metric,
		"diskNameFilter":   &me.DiskNameFilter,
		"tagFilters":       &me.TagFilters,
	}); err != nil {
		return err
	}
	return nil
}

// Metric The metric to monitor.
type Metric string

// Metrics offers the known enum values
var Metrics = struct {
	LowDiskSpace       Metric
	LowInodes          Metric
	ReadTimeExceeding  Metric
	WriteTimeExceeding Metric
}{
	"LOW_DISK_SPACE",
	"LOW_INODES",
	"READ_TIME_EXCEEDING",
	"WRITE_TIME_EXCEEDING",
}
