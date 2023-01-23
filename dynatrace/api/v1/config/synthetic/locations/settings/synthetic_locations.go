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

package locations

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SyntheticLocations is a list of synthetic locations
type SyntheticLocations struct {
	Locations []*SyntheticLocation `json:"locations"` // A list of synthetic locations
}

func (me *SyntheticLocations) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location": {
			Type:        schema.TypeList,
			Description: "The name of the location",
			Optional:    true,
			Elem:        &schema.Resource{Schema: new(SyntheticLocation).Schema()},
		},
	}
}

func (me *SyntheticLocations) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("location", me.Locations); err != nil {
		return err
	}
	return nil
}

type SyntheticLocation struct {
	ID            string         `json:"entityId"`      // The Dynatrace entity ID of the location
	Name          string         `json:"name"`          // The name of the location
	Type          LocationType   `json:"type"`          // The type of the location
	Status        *Status        `json:"status"`        // The status of the location: \n\n* `ENABLED`: The location is displayed as active in the UI. You can assign monitors to the location. \n* `DISABLED`: The location is displayed as inactive in the UI. You can't assign monitors to the location. Monitors already assigned to the location will stay there and will be executed from the location. \n* `HIDDEN`: The location is not displayed in the UI. You can't assign monitors to the location. You can only set location as `HIDDEN` when no monitor is assigned to it
	CloudPlatform *CloudPlatform `json:"cloudPlatform"` // The cloud provider where the location is hosted. \n\n Only applicable to `PUBLIC` locations
	IPs           []string       `json:"ips"`           // The list of IP addresses assigned to the location. \n\n Only applicable to `PUBLIC` locations
	Stage         *Stage         `json:"stage"`         // The release stage of the location
}

func (me *SyntheticLocation) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_id": {
			Type:        schema.TypeString,
			Description: "The unique ID of the location",
			Optional:    true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "The name of the location",
			Optional:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "The type of the location. Supported values are `PUBLIC`, `PRIVATE` and `CLUSTER`",
			Optional:    true,
		},
		"status": {
			Type:        schema.TypeString,
			Description: "The status of the location: \n\n* `ENABLED`: The location is displayed as active in the UI. You can assign monitors to the location. \n* `DISABLED`: The location is displayed as inactive in the UI. You can't assign monitors to the location. Monitors already assigned to the location will stay there and will be executed from the location. \n* `HIDDEN`: The location is not displayed in the UI. You can't assign monitors to the location. You can only set location as `HIDDEN` when no monitor is assigned to it",
			Optional:    true,
			Computed:    true,
		},
		"stage": {
			Type:        schema.TypeString,
			Description: "The release stage of the location",
			Optional:    true,
			Computed:    true,
		},
		"cloud_platform": {
			Type:        schema.TypeString,
			Description: "The cloud provider where the location is hosted. \n\n Only applicable to `PUBLIC` locations",
			Optional:    true,
			Computed:    true,
		},
		"ips": {
			Type:        schema.TypeList,
			Description: "The list of IP addresses assigned to the location. \n\n Only applicable to `PUBLIC` locations",
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *SyntheticLocation) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Encode("entity_id", me.ID); err != nil {
		return err
	}
	if err := properties.Encode("name", me.Name); err != nil {
		return err
	}
	if err := properties.Encode("type", string(me.Type)); err != nil {
		return err
	}
	if err := properties.Encode("stage", me.Stage); err != nil {
		return err
	}
	if err := properties.Encode("status", me.Status); err != nil {
		return err
	}
	if err := properties.Encode("cloud_platform", me.CloudPlatform); err != nil {
		return err
	}
	if err := properties.Encode("ips", append([]string{}, me.IPs...)); err != nil {
		return err
	}
	return nil
}

func (me *SyntheticLocation) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]any{
		"entity_id":      &me.ID,
		"name":           &me.Name,
		"type":           &me.Type,
		"stage":          &me.Stage,
		"status":         &me.Status,
		"cloud_platform": &me.CloudPlatform,
		"ips":            &me.IPs,
	})
}

type Stage string

var Stages = struct {
	Beta       Stage
	ComingSoon Stage
	GA         Stage
}{
	Stage("BETA"),
	Stage("COMING_SOON"),
	Stage("GA"),
}

type CloudPlatform string

var CloudPlatforms = struct {
	Alibaba        CloudPlatform
	AmaconEC2      CloudPlatform
	Azure          CloudPlatform
	DynatraceCloud CloudPlatform
	GoogleCloud    CloudPlatform
	Interoute      CloudPlatform
	Other          CloudPlatform
	Undefined      CloudPlatform
}{
	CloudPlatform("ALIBABA"),
	CloudPlatform("AMAZON_EC2"),
	CloudPlatform("AZURE"),
	CloudPlatform("DYNATRACE_CLOUD"),
	CloudPlatform("GOOGLE_CLOUD"),
	CloudPlatform("INTEROUTE"),
	CloudPlatform("OTHER"),
	CloudPlatform("UNDEFINED"),
}

type Status string

var Statuses = struct {
	Disabled Status
	Enabled  Status
	Hidden   Status
}{
	Status("DISABLED"),
	Status("ENABLED"),
	Status("HIDDEN"),
}

type LocationType string

var LocationTypes = struct {
	Public  LocationType
	Private LocationType
	Cluster LocationType
}{
	LocationType("PUBLIC"),
	LocationType("PRIVATE"),
	LocationType("CLUSTER"),
}
