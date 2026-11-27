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

package location

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PrivateSyntheticLocation struct {
	ID                               *string           `json:"id,omitempty"`
	AutoUpdateChromium               bool              `json:"autoUpdateChromium"`
	AvailabilityLocationOutage       bool              `json:"availabilityLocationOutage"`
	AvailabilityNodeOutage           bool              `json:"availabilityNodeOutage"`
	AvailabilityNotificationsEnabled bool              `json:"availabilityNotificationsEnabled"`
	BrowserExecutionSupported        bool              `json:"browserExecutionSupported"`
	City                             *string           `json:"city,omitempty"`
	CountryCode                      *string           `json:"countryCode,omitempty"`
	CountryName                      *string           `json:"countryName,omitempty"`
	DeploymentType                   *DeploymentType   `json:"deploymentType,omitempty"`
	FipsMode                         *FipsMode         `json:"fipsMode,omitempty"`
	GeoLocationID                    *string           `json:"geoLocationId,omitempty"`
	Latitude                         float64           `json:"latitude"`
	LocationNodeOutageDelayInMinutes *int              `json:"locationNodeOutageDelayInMinutes,omitempty"`
	Longitude                        float64           `json:"longitude"`
	MaxActiveGateCount               *int              `json:"maxActiveGateCount,omitempty"`
	MinActiveGateCount               *int              `json:"minActiveGateCount,omitempty"`
	NAMExecutionSupported            bool              `json:"namExecutionSupported"`
	Name                             string            `json:"name"`
	NodeNames                        map[string]string `json:"nodeNames,omitempty"`
	NodeSize                         *NodeSize         `json:"nodeSize,omitempty"`
	Nodes                            []string          `json:"nodes"`
	RegionCode                       *string           `json:"regionCode,omitempty"`
	RegionName                       *string           `json:"regionName,omitempty"`
	Status                           *Status           `json:"status,omitempty"`
	Type                             LocationType      `json:"type"`
	UseNewKubernetesVersion          bool              `json:"useNewKubernetesVersion"`
}

func (me *PrivateSyntheticLocation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"auto_update_chromium": {
			Type:        schema.TypeBool,
			Description: "Non-containerized location property. Auto upgrade of Chromium is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"availability_location_outage": {
			Type:        schema.TypeBool,
			Description: "Alerting for location outage is enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"availability_node_outage": {
			Type:        schema.TypeBool,
			Description: "Alerting for node outage is enabled (`true`) or disabled (`false`). \n  \n  If enabled, the outage of *any* node in the location triggers an alert",
			Optional:    true,
		},
		"availability_notifications_enabled": {
			Type:        schema.TypeBool,
			Description: "Notifications for location and node outage are enabled (`true`) or disabled (`false`)",
			Optional:    true,
		},
		"browser_execution_supported": {
			Type:        schema.TypeBool,
			Description: "Containerized location property. Boolean value describes if browser monitors will be executed on this location",
			Optional:    true,
		},
		"city": {
			Type:        schema.TypeString,
			Description: "The city of the location",
			Optional:    true,
		},
		"country_code": {
			Type:        schema.TypeString,
			Description: "The country code of the location.",
			Optional:    true,
		},
		"country_name": {
			Type:        schema.TypeString,
			Description: "The country name of the location",
			Computed:    true,
		},
		"deployment_type": {
			Type:        schema.TypeString,
			Description: "The deployment type of the location. \n  \n  `STANDARD`: The location is deployed on Windows or Linux. \n  \n  `KUBERNETES`: The location is deployed on Kubernetes. \n  Possible values: `KUBERNETES`, `OPENSHIFT`, `STANDARD`, `UNKNOWN`",
			Optional:    true,
		},
		"fips_mode": {
			Type:        schema.TypeString,
			Description: "Containerized location property indicating whether FIPS mode is enabled on this location. Possible values: `DISABLED`, `ENABLED`, `ENABLED_WITH_CORPORATE_PROXY`",
			Optional:    true,
			Default:     string(FipsModes.Disabled),
		},
		"geo_location_id": {
			Type:        schema.TypeString,
			Description: "The Dynatrace GeoLocation ID",
			Computed:    true,
		},
		"latitude": {
			Type:        schema.TypeFloat,
			Description: "The latitude of the location in `DDD.dddd` format",
			Required:    true,
		},
		"location_node_outage_delay_in_minutes": {
			Type:        schema.TypeInt,
			Description: "Alert if location or node outage lasts longer than *X* minutes. \n  \n  Only applicable when `availability_location_outage` or `availability_node_outage` is set to `true`",
			Optional:    true,
		},
		"longitude": {
			Type:        schema.TypeFloat,
			Description: "The longitude of the location in `DDD.dddd` format",
			Required:    true,
		},
		"max_active_gate_count": {
			Type:        schema.TypeInt,
			Description: "Containerized location property. The maximum number of ActiveGates deployed for the location (required for a Kubernetes location)",
			Optional:    true,
		},
		"min_active_gate_count": {
			Type:        schema.TypeInt,
			Description: "Containerized location property. The minimum number of ActiveGates deployed for the location (required for a Kubernetes location)",
			Optional:    true,
		},
		"nam_execution_supported": {
			Type:        schema.TypeBool,
			Description: "Containerized location property. Boolean value describes if ICMP monitors will be executed on this location",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the location",
			Required:    true,
		},
		"node_names": {
			Type:        schema.TypeMap,
			Description: "A mapping id to name of the nodes belonging to the location",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"node_size": {
			Type:        schema.TypeString,
			Description: "Containerized location property. The size of a containerized node deployed for the location (required for a Kubernetes location). `L` is not supported in containerized locations. Possible values: `UNSUPPORTED`, `XS`, `S`, `M`",
			Optional:    true,
			Default:     string(NodeSizes.Unsupported),
		},
		"nodes": {
			Type:        schema.TypeSet,
			Description: "A list of synthetic nodes belonging to the location.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"region_code": {
			Type:        schema.TypeString,
			Description: "The region code of the location.",
			Optional:    true,
		},
		"region_name": {
			Type:        schema.TypeString,
			Description: "The region name of the location",
			Computed:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The status of the location: \n  \n  * `ENABLED`: The location is displayed as active in the UI. You can assign monitors to the location. \n  * `DISABLED`: The location is displayed as inactive in the UI. You can't assign monitors to the location. Monitors already assigned to the location will stay there and will be executed from the location. \n  * `HIDDEN`: The location is not displayed in the UI. You can't assign monitors to the location. You can only set location as `HIDDEN` when no monitor is assigned to it",
			Optional:    true,
			Default:     string(Statuses.Enabled),
		},
		"use_new_kubernetes_version": {
			Type:        schema.TypeBool,
			Description: "Containerized location property. Boolean value describes which Kubernetes version will be used. Set to `false` for version 1.23+ and older than version 1.26. Set to `true` for version 1.26+",
			Optional:    true,
			Default:     true,
		},
	}
}

func (me *PrivateSyntheticLocation) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"auto_update_chromium":                  me.AutoUpdateChromium,
		"availability_location_outage":          me.AvailabilityLocationOutage,
		"availability_node_outage":              me.AvailabilityNodeOutage,
		"availability_notifications_enabled":    me.AvailabilityNotificationsEnabled,
		"browser_execution_supported":           me.BrowserExecutionSupported,
		"city":                                  me.City,
		"country_code":                          me.CountryCode,
		"country_name":                          me.CountryName,
		"deployment_type":                       me.DeploymentType,
		"fips_mode":                             me.FipsMode,
		"geo_location_id":                       me.GeoLocationID,
		"latitude":                              me.Latitude,
		"location_node_outage_delay_in_minutes": me.LocationNodeOutageDelayInMinutes,
		"longitude":                             me.Longitude,
		"max_active_gate_count":                 me.MaxActiveGateCount,
		"min_active_gate_count":                 me.MinActiveGateCount,
		"nam_execution_supported":               me.NAMExecutionSupported,
		"name":                                  me.Name,
		"node_names":                            me.NodeNames,
		"node_size":                             me.NodeSize,
		"nodes":                                 me.Nodes,
		"region_code":                           me.RegionCode,
		"region_name":                           me.RegionName,
		"status":                                me.Status,
		"use_new_kubernetes_version":            me.UseNewKubernetesVersion,
	})
}

// UnmarshalHCL skips the computed properties, which are read-only and must not be sent to the API.
func (me *PrivateSyntheticLocation) UnmarshalHCL(decoder hcl.Decoder) error {
	me.Type = LocationTypes.Private
	if err := decoder.DecodeAll(map[string]any{
		"auto_update_chromium":                  &me.AutoUpdateChromium,
		"availability_location_outage":          &me.AvailabilityLocationOutage,
		"availability_node_outage":              &me.AvailabilityNodeOutage,
		"availability_notifications_enabled":    &me.AvailabilityNotificationsEnabled,
		"browser_execution_supported":           &me.BrowserExecutionSupported,
		"city":                                  &me.City,
		"country_code":                          &me.CountryCode,
		"deployment_type":                       &me.DeploymentType,
		"fips_mode":                             &me.FipsMode,
		"latitude":                              &me.Latitude,
		"location_node_outage_delay_in_minutes": &me.LocationNodeOutageDelayInMinutes,
		"longitude":                             &me.Longitude,
		"max_active_gate_count":                 &me.MaxActiveGateCount,
		"min_active_gate_count":                 &me.MinActiveGateCount,
		"nam_execution_supported":               &me.NAMExecutionSupported,
		"name":                                  &me.Name,
		"node_size":                             &me.NodeSize,
		"nodes":                                 &me.Nodes,
		"region_code":                           &me.RegionCode,
		"status":                                &me.Status,
		"use_new_kubernetes_version":            &me.UseNewKubernetesVersion,
	}); err != nil {
		return err
	}
	// nodes is required by the API, so send an empty list instead of omitting it.
	if me.Nodes == nil {
		me.Nodes = []string{}
	}
	return nil
}
