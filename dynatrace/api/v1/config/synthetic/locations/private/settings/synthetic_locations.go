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
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
)

type SyntheticLocations struct {
	Locations LocationCollectionElements `json:"locations"`
}

func (me *SyntheticLocations) ToStubs() settings.Stubs {
	stubs := settings.Stubs{}
	if len(me.Locations) == 0 {
		return stubs
	}
	for _, location := range me.Locations {
		stub := settings.Stub{ID: *location.ID, Name: location.Name}
		stubs = append(stubs, &stub)
	}
	return stubs
}

type LocationCollectionElements []LocationCollectionElement

// LocationCollectionElement represents a synthetic location
type LocationCollectionElement struct {
	Name          string         `json:"name"`                    // The name of the location
	ID            *string        `json:"entityId"`                // The Dynatrace entity ID of the location
	Type          LocationType   `json:"type"`                    // The type of the location
	CloudPlatform *CloudPlatform `json:"cloudPlatform,omitempty"` // The cloud provider where the location is hosted. Only applicable to `PUBLIC` locations
	IPs           []string       `json:"ips,omitempty"`           // The list of IP addresses assigned to the location. Only applicable to `PUBLIC` locations
	Stage         *Stage         `json:"stage,omitempty"`         // The release stage of the location
}
