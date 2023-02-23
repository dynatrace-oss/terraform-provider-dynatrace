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

package txmonitoring

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	GroupCicsRegions                         bool `json:"groupCicsRegions"`                         // If enabled, CICS regions belonging to the same CICSPlex will be grouped into a single process group. If disabled, a process group will be created for each CICS region.
	GroupImsRegions                          bool `json:"groupImsRegions"`                          // If enabled, IMS regions belonging to the same subsystem will be grouped into a single process group. If disabled, a process group will be created for each IMS region.
	MonitorAllCtgProtocols                   bool `json:"monitorAllCtgProtocols"`                   // If enabled, the CICS Transaction Gateway sensor will trace all EXCI requests including those that are using the TCP/IP or SNA protocol.
	MonitorAllIncomingWebRequests            bool `json:"monitorAllIncomingWebRequests"`            // Dynatrace automatically traces incoming web requests when they are called by already-monitored services. Enable this setting to monitor all incoming web requests. We recommend enabling it only over a short period of time.
	NodeLimit                                int  `json:"nodeLimit"`                                // We recommend the default limit of 500 nodes. The value 0 means unlimited number of nodes.
	ZosCicsServiceDetectionUsesTransactionID bool `json:"zosCicsServiceDetectionUsesTransactionId"` // If enabled, a CICS service will be created for each monitored transaction ID within a process group. If disabled, a CICS service will be created for each monitored CICS region within a process group. We recommend enabling it only when the CICS regions are grouped by their CICSPlex.
	ZosImsServiceDetectionUsesTransactionID  bool `json:"zosImsServiceDetectionUsesTransactionId"`  // If enabled, an IMS service will be created for each monitored transaction ID within a process group. If disabled, an IMS service will be created for each monitored IMS region within a process group. We recommend enabling it only when the IMS regions are grouped by their subsystem.
}

func (me *Settings) Name() string {
	return "mainframe_transaction_monitoring"
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_cics_regions": {
			Type:        schema.TypeBool,
			Description: "If enabled, CICS regions belonging to the same CICSPlex will be grouped into a single process group. If disabled, a process group will be created for each CICS region.",
			Required:    true,
		},
		"group_ims_regions": {
			Type:        schema.TypeBool,
			Description: "If enabled, IMS regions belonging to the same subsystem will be grouped into a single process group. If disabled, a process group will be created for each IMS region.",
			Required:    true,
		},
		"monitor_all_ctg_protocols": {
			Type:        schema.TypeBool,
			Description: "If enabled, the CICS Transaction Gateway sensor will trace all EXCI requests including those that are using the TCP/IP or SNA protocol.",
			Required:    true,
		},
		"monitor_all_incoming_web_requests": {
			Type:        schema.TypeBool,
			Description: "Dynatrace automatically traces incoming web requests when they are called by already-monitored services. Enable this setting to monitor all incoming web requests. We recommend enabling it only over a short period of time.",
			Required:    true,
		},
		"node_limit": {
			Type:        schema.TypeInt,
			Description: "We recommend the default limit of 500 nodes. The value 0 means unlimited number of nodes.",
			Required:    true,
		},
		"zos_cics_service_detection_uses_transaction_id": {
			Type:        schema.TypeBool,
			Description: "If enabled, a CICS service will be created for each monitored transaction ID within a process group. If disabled, a CICS service will be created for each monitored CICS region within a process group. We recommend enabling it only when the CICS regions are grouped by their CICSPlex.",
			Required:    true,
		},
		"zos_ims_service_detection_uses_transaction_id": {
			Type:        schema.TypeBool,
			Description: "If enabled, an IMS service will be created for each monitored transaction ID within a process group. If disabled, an IMS service will be created for each monitored IMS region within a process group. We recommend enabling it only when the IMS regions are grouped by their subsystem.",
			Required:    true,
		},
	}
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"group_cics_regions":                             me.GroupCicsRegions,
		"group_ims_regions":                              me.GroupImsRegions,
		"monitor_all_ctg_protocols":                      me.MonitorAllCtgProtocols,
		"monitor_all_incoming_web_requests":              me.MonitorAllIncomingWebRequests,
		"node_limit":                                     me.NodeLimit,
		"zos_cics_service_detection_uses_transaction_id": me.ZosCicsServiceDetectionUsesTransactionID,
		"zos_ims_service_detection_uses_transaction_id":  me.ZosImsServiceDetectionUsesTransactionID,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"group_cics_regions":                             &me.GroupCicsRegions,
		"group_ims_regions":                              &me.GroupImsRegions,
		"monitor_all_ctg_protocols":                      &me.MonitorAllCtgProtocols,
		"monitor_all_incoming_web_requests":              &me.MonitorAllIncomingWebRequests,
		"node_limit":                                     &me.NodeLimit,
		"zos_cics_service_detection_uses_transaction_id": &me.ZosCicsServiceDetectionUsesTransactionID,
		"zos_ims_service_detection_uses_transaction_id":  &me.ZosImsServiceDetectionUsesTransactionID,
	})
}
