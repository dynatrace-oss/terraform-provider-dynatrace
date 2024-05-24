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

package anomalydetectors

import (
	"fmt"
	"slices"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Alerts []*Alert

func (me *Alerts) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alert": {
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Description: "",
			Elem:        &schema.Resource{Schema: new(Alert).Schema()},
		},
	}
}

func (me Alerts) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeSlice("alert", me)
}

func (me *Alerts) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("alert", me)
}

type Alert struct {
	SampleCountThresholds            *SampleCountThresholds            `json:"sampleCountThresholds,omitempty"`
	SampleCountThresholdsImmediately *SampleCountThresholdsImmediately `json:"sampleCountThresholdsImmediately,omitempty"`
	ThresholdMebibytes               *float64                          `json:"thresholdMebibytes,omitempty"`
	ThresholdMilliseconds            *float64                          `json:"thresholdMilliseconds,omitempty"`
	ThresholdNumber                  *float64                          `json:"thresholdNumber,omitempty"`
	ThresholdPercent                 *float64                          `json:"thresholdPercent,omitempty"`
	Trigger                          Trigger                           `json:"trigger"` // Possible Values: `AVAILABLE_DISK_SPACE_MEBIBYTES_BELOW`, `AVAILABLE_DISK_SPACE_PERCENT_BELOW`, `AVAILABLE_INODES_NUMBER_BELOW`, `AVAILABLE_INODES_PERCENT_BELOW`, `READ_ONLY_FILE_SYSTEM`, `READ_TIME_EXCEEDING`, `WRITE_TIME_EXCEEDING`
}

func (me *Alert) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sample_count_thresholds": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SampleCountThresholds).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"sample_count_thresholds_immediately": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Optional:    true, // precondition
			Elem:        &schema.Resource{Schema: new(SampleCountThresholdsImmediately).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"threshold_mebibytes": {
			Type:        schema.TypeFloat,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"threshold_milliseconds": {
			Type:        schema.TypeFloat,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"threshold_number": {
			Type:        schema.TypeFloat,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"threshold_percent": {
			Type:        schema.TypeFloat,
			Description: "no documentation available",
			Optional:    true, // precondition
		},
		"trigger": {
			Type:        schema.TypeString,
			Description: "Possible Values: `AVAILABLE_DISK_SPACE_MEBIBYTES_BELOW`, `AVAILABLE_DISK_SPACE_PERCENT_BELOW`, `AVAILABLE_INODES_NUMBER_BELOW`, `AVAILABLE_INODES_PERCENT_BELOW`, `READ_ONLY_FILE_SYSTEM`, `READ_TIME_EXCEEDING`, `WRITE_TIME_EXCEEDING`",
			Required:    true,
		},
	}
}

func (me *Alert) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"sample_count_thresholds":             me.SampleCountThresholds,
		"sample_count_thresholds_immediately": me.SampleCountThresholdsImmediately,
		"threshold_mebibytes":                 me.ThresholdMebibytes,
		"threshold_milliseconds":              me.ThresholdMilliseconds,
		"threshold_number":                    me.ThresholdNumber,
		"threshold_percent":                   me.ThresholdPercent,
		"trigger":                             me.Trigger,
	})
}

func (me *Alert) HandlePreconditions() error {
	if (me.ThresholdMebibytes == nil) && (slices.Contains([]string{"AVAILABLE_DISK_SPACE_MEBIBYTES_BELOW"}, string(me.Trigger))) {
		me.ThresholdMebibytes = opt.NewFloat64(0.0)
	}
	if (me.ThresholdMilliseconds == nil) && (slices.Contains([]string{"READ_TIME_EXCEEDING", "WRITE_TIME_EXCEEDING"}, string(me.Trigger))) {
		me.ThresholdMilliseconds = opt.NewFloat64(0.0)
	}
	if (me.ThresholdNumber == nil) && (slices.Contains([]string{"AVAILABLE_INODES_NUMBER_BELOW"}, string(me.Trigger))) {
		me.ThresholdNumber = opt.NewFloat64(0.0)
	}
	if (me.ThresholdPercent == nil) && (slices.Contains([]string{"AVAILABLE_DISK_SPACE_PERCENT_BELOW", "AVAILABLE_INODES_PERCENT_BELOW"}, string(me.Trigger))) {
		me.ThresholdPercent = opt.NewFloat64(0.0)
	}
	if (me.SampleCountThresholds == nil) && (slices.Contains([]string{"AVAILABLE_INODES_NUMBER_BELOW", "AVAILABLE_INODES_PERCENT_BELOW", "AVAILABLE_DISK_SPACE_PERCENT_BELOW", "AVAILABLE_DISK_SPACE_MEBIBYTES_BELOW", "READ_TIME_EXCEEDING", "WRITE_TIME_EXCEEDING"}, string(me.Trigger))) {
		return fmt.Errorf("'sample_count_thresholds' must be specified if 'trigger' is set to '%v'", me.Trigger)
	}
	if (me.SampleCountThresholds != nil) && (!slices.Contains([]string{"AVAILABLE_INODES_NUMBER_BELOW", "AVAILABLE_INODES_PERCENT_BELOW", "AVAILABLE_DISK_SPACE_PERCENT_BELOW", "AVAILABLE_DISK_SPACE_MEBIBYTES_BELOW", "READ_TIME_EXCEEDING", "WRITE_TIME_EXCEEDING"}, string(me.Trigger))) {
		return fmt.Errorf("'sample_count_thresholds' must not be specified if 'trigger' is set to '%v'", me.Trigger)
	}
	if (me.SampleCountThresholdsImmediately == nil) && (slices.Contains([]string{"READ_ONLY_FILE_SYSTEM"}, string(me.Trigger))) {
		return fmt.Errorf("'sample_count_thresholds_immediately' must be specified if 'trigger' is set to '%v'", me.Trigger)
	}
	if (me.SampleCountThresholdsImmediately != nil) && (!slices.Contains([]string{"READ_ONLY_FILE_SYSTEM"}, string(me.Trigger))) {
		return fmt.Errorf("'sample_count_thresholds_immediately' must not be specified if 'trigger' is set to '%v'", me.Trigger)
	}
	return nil
}

func (me *Alert) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"sample_count_thresholds":             &me.SampleCountThresholds,
		"sample_count_thresholds_immediately": &me.SampleCountThresholdsImmediately,
		"threshold_mebibytes":                 &me.ThresholdMebibytes,
		"threshold_milliseconds":              &me.ThresholdMilliseconds,
		"threshold_number":                    &me.ThresholdNumber,
		"threshold_percent":                   &me.ThresholdPercent,
		"trigger":                             &me.Trigger,
	})
}
