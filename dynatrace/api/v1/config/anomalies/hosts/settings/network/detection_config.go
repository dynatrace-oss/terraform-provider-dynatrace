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

package network

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/droppedpackets"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/errors"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/retransmission"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/tcp"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/anomalies/hosts/settings/network/utilization"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DetectionConfig struct {
	NetworkDroppedPacketsDetection     *droppedpackets.DetectionConfig `json:"networkDroppedPacketsDetection"`     // Configuration of high number of dropped packets detection.
	HighNetworkDetection               *utilization.DetectionConfig    `json:"highNetworkDetection"`               // Configuration of high network utilization detection.
	NetworkHighRetransmissionDetection *retransmission.DetectionConfig `json:"networkHighRetransmissionDetection"` // Configuration of high retransmission rate detection.
	NetworkTcpProblemsDetection        *tcp.DetectionConfig            `json:"networkTcpProblemsDetection"`        // Configuration of TCP connectivity problems detection.
	NetworkErrorsDetection             *errors.DetectionConfig         `json:"networkErrorsDetection"`             // Configuration of high number of network errors detection.
}

func (me *DetectionConfig) IsConfigured() bool {
	if me.NetworkDroppedPacketsDetection != nil && me.NetworkDroppedPacketsDetection.Enabled {
		return true
	}
	if me.HighNetworkDetection != nil && me.HighNetworkDetection.Enabled {
		return true
	}
	if me.NetworkHighRetransmissionDetection != nil && me.NetworkHighRetransmissionDetection.Enabled {
		return true
	}
	if me.NetworkTcpProblemsDetection != nil && me.NetworkTcpProblemsDetection.Enabled {
		return true
	}
	if me.NetworkErrorsDetection != nil && me.NetworkErrorsDetection.Enabled {
		return true
	}
	return false
}

func (me *DetectionConfig) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dropped_packets": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high number of dropped packets detection",
			Elem:        &schema.Resource{Schema: new(droppedpackets.DetectionConfig).Schema()},
		},
		"utilization": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high network utilization detection",
			Elem:        &schema.Resource{Schema: new(utilization.DetectionConfig).Schema()},
		},
		"retransmission": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high retransmission rate detection",
			Elem:        &schema.Resource{Schema: new(retransmission.DetectionConfig).Schema()},
		},
		"connectivity": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of TCP connectivity problems detection",
			Elem:        &schema.Resource{Schema: new(tcp.DetectionConfig).Schema()},
		},
		"errors": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Configuration of high number of network errors detection",
			Elem:        &schema.Resource{Schema: new(errors.DetectionConfig).Schema()},
		},
	}
}

func (me *DetectionConfig) MarshalHCL(properties hcl.Properties) error {
	if !me.IsConfigured() {
		return nil
	}
	return properties.EncodeAll(map[string]any{
		"dropped_packets": me.NetworkDroppedPacketsDetection,
		"utilization":     me.HighNetworkDetection,
		"retransmission":  me.NetworkHighRetransmissionDetection,
		"connectivity":    me.NetworkTcpProblemsDetection,
		"errors":          me.NetworkErrorsDetection,
	})
}

func (me *DetectionConfig) UnmarshalHCL(decoder hcl.Decoder) error {
	me.NetworkDroppedPacketsDetection = &droppedpackets.DetectionConfig{Enabled: false}
	me.HighNetworkDetection = &utilization.DetectionConfig{Enabled: false}
	me.NetworkHighRetransmissionDetection = &retransmission.DetectionConfig{Enabled: false}
	me.NetworkTcpProblemsDetection = &tcp.DetectionConfig{Enabled: false}
	me.NetworkErrorsDetection = &errors.DetectionConfig{Enabled: false}
	if _, ok := decoder.GetOk("dropped_packets.#"); ok {
		me.NetworkDroppedPacketsDetection = new(droppedpackets.DetectionConfig)
		if err := me.NetworkDroppedPacketsDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "dropped_packets", 0)); err != nil {
			return err
		}
	}

	if _, ok := decoder.GetOk("utilization.#"); ok {
		me.HighNetworkDetection = new(utilization.DetectionConfig)
		if err := me.HighNetworkDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "utilization", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("retransmission.#"); ok {
		me.NetworkHighRetransmissionDetection = new(retransmission.DetectionConfig)
		if err := me.NetworkHighRetransmissionDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "retransmission", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("connectivity.#"); ok {
		me.NetworkTcpProblemsDetection = new(tcp.DetectionConfig)
		if err := me.NetworkTcpProblemsDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "connectivity", 0)); err != nil {
			return err
		}
	}
	if _, ok := decoder.GetOk("errors.#"); ok {
		me.NetworkErrorsDetection = new(errors.DetectionConfig)
		if err := me.NetworkErrorsDetection.UnmarshalHCL(hcl.NewDecoder(decoder, "errors", 0)); err != nil {
			return err
		}
	}
	return nil
}
