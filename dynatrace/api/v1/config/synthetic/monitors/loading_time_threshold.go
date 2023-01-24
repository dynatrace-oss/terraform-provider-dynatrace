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

package monitors

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LoadingTimeThresholds []*LoadingTimeThreshold

func (me *LoadingTimeThresholds) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"threshold": {
			Type:        schema.TypeList,
			Description: "The list of performance threshold rules",
			Required:    true,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(LoadingTimeThreshold).Schema()},
		},
	}
}

func (me LoadingTimeThresholds) MarshalHCL(properties hcl.Properties) error {
	entries := []any{}
	if len(me) > 0 {
		for _, entry := range me {
			marshalled := hcl.Properties{}
			if err := entry.MarshalHCL(marshalled); err == nil {
				entries = append(entries, marshalled)
			} else {
				return err
			}
		}
		properties["threshold"] = entries
	}
	return nil
}

func (me *LoadingTimeThresholds) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeSlice("threshold", me)
}

// LoadingTimeThreshold The performance threshold rule
type LoadingTimeThreshold struct {
	Type         LoadingTimeThresholdType `json:"type"`         // The type of the threshold: total loading time or action loading time
	ValueMs      int32                    `json:"valueMs"`      // Notify if monitor takes longer than *X* milliseconds to load
	RequestIndex *int32                   `json:"requestIndex"` // Specify the request to which an ACTION threshold applies
	EventIndex   *int32                   `json:"eventIndex"`   // Specify the event to which an ACTION threshold applies
}

func (me *LoadingTimeThreshold) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the threshold: `TOTAL` (total loading time) or `ACTION` (action loading time)",
			Optional:    true,
		},
		"value_ms": {
			Type:        schema.TypeInt,
			Description: "Notify if monitor takes longer than *X* milliseconds to load",
			Required:    true,
		},
		"request_index": {
			Type:        schema.TypeInt,
			Description: "Specify the request to which an ACTION threshold applies",
			Optional:    true,
		},
		"event_index": {
			Type:        schema.TypeInt,
			Description: "Specify the event to which an ACTION threshold applies",
			Optional:    true,
		},
	}
}

func (me *LoadingTimeThreshold) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("value_ms", int(me.ValueMs)); err != nil {
		return err
	}
	if err := properties.Encode("request_index", me.RequestIndex); err != nil {
		return err
	}
	if err := properties.Encode("event_index", me.EventIndex); err != nil {
		return err
	}
	return nil
}

func (me *LoadingTimeThreshold) UnmarshalHCL(decoder hcl.Decoder) error {
	if err := decoder.Decode("type", &me.Type); err != nil {
		return err
	}
	if err := decoder.Decode("value_ms", &me.ValueMs); err != nil {
		return err
	}
	if err := decoder.Decode("request_index", &me.RequestIndex); err != nil {
		return err
	}
	if err := decoder.Decode("event_index", &me.EventIndex); err != nil {
		return err
	}
	return nil
}

// LoadingTimeThresholdType The type of the threshold: total loading time or action loading time
type LoadingTimeThresholdType string

// LoadingTimeThresholdTypes offers the known enum values
var LoadingTimeThresholdTypes = struct {
	Action LoadingTimeThresholdType
	Total  LoadingTimeThresholdType
}{
	"ACTION",
	"TOTAL",
}
