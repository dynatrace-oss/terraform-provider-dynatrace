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

package hosts

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Network struct {
	HighNetworkDetection               *HighNetworkDetection               `json:"highNetworkDetection"`
	NetworkDroppedPacketsDetection     *NetworkDroppedPacketsDetection     `json:"networkDroppedPacketsDetection"`
	NetworkErrorsDetection             *NetworkErrorsDetection             `json:"networkErrorsDetection"`
	NetworkHighRetransmissionDetection *NetworkHighRetransmissionDetection `json:"networkHighRetransmissionDetection"`
	NetworkTcpProblemsDetection        *NetworkTcpProblemsDetection        `json:"networkTcpProblemsDetection"`
}

func (me *Network) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"high_network_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(HighNetworkDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"network_dropped_packets_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(NetworkDroppedPacketsDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"network_errors_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(NetworkErrorsDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"network_high_retransmission_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(NetworkHighRetransmissionDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
		"network_tcp_problems_detection": {
			Type:        schema.TypeList,
			Description: "no documentation available",
			Required:    true,
			Elem:        &schema.Resource{Schema: new(NetworkTcpProblemsDetection).Schema()},
			MinItems:    1,
			MaxItems:    1,
		},
	}
}

func (me *Network) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"high_network_detection":                me.HighNetworkDetection,
		"network_dropped_packets_detection":     me.NetworkDroppedPacketsDetection,
		"network_errors_detection":              me.NetworkErrorsDetection,
		"network_high_retransmission_detection": me.NetworkHighRetransmissionDetection,
		"network_tcp_problems_detection":        me.NetworkTcpProblemsDetection,
	})
}

func (me *Network) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"high_network_detection":                &me.HighNetworkDetection,
		"network_dropped_packets_detection":     &me.NetworkDroppedPacketsDetection,
		"network_errors_detection":              &me.NetworkErrorsDetection,
		"network_high_retransmission_detection": &me.NetworkHighRetransmissionDetection,
		"network_tcp_problems_detection":        &me.NetworkTcpProblemsDetection,
	})
}
