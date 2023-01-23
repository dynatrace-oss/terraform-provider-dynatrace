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

package dashboards

import "encoding/json"

// TileType has no documentation
type TileType string

// TileTypes offers the known enum values
var TileTypes = struct {
	Application             TileType
	Applications            TileType
	ApplicationMethod       TileType
	ApplicationWorldMap     TileType
	AWS                     TileType
	BounceRate              TileType
	CustomApplication       TileType
	CustomCharting          TileType
	Database                TileType
	DatabasesOverview       TileType
	DemKeyUserAction        TileType
	DeviceApplicationMethod TileType
	DTAQL                   TileType
	Host                    TileType
	Hosts                   TileType
	LogAnalytics            TileType
	Markdown                TileType
	MobileApplication       TileType
	OpenStack               TileType
	OpenStackAVZone         TileType
	OpenStackProject        TileType
	ProcessGroupsOne        TileType
	Resources               TileType
	Services                TileType
	ServiceVersatile        TileType
	SessionMetrics          TileType
	SyntheticHTTPMonitor    TileType
	SyntheticSingleExtTest  TileType
	SyntheticSingleWebCheck TileType
	SyntheticTests          TileType
	ThirdPartyMostActive    TileType
	UEMConversionsOverall   TileType
	UEMConversionsPerGoal   TileType
	UEMJserrorsOverall      TileType
	UEMKeyUserActions       TileType
	Users                   TileType
	Virtualization          TileType
	Header                  TileType
}{
	"APPLICATION",
	"APPLICATIONS",
	"APPLICATION_METHOD",
	"APPLICATION_WORLDMAP",
	"AWS",
	"BOUNCE_RATE",
	"CUSTOM_APPLICATION",
	"CUSTOM_CHARTING",
	"DATABASE",
	"DATABASES_OVERVIEW",
	"DEM_KEY_USER_ACTION",
	"DEVICE_APPLICATION_METHOD",
	"DTAQL",
	"HOST",
	"HOSTS",
	"LOG_ANALYTICS",
	"MARKDOWN",
	"MOBILE_APPLICATION",
	"OPENSTACK",
	"OPENSTACK_AV_ZONE",
	"OPENSTACK_PROJECT",
	"PROCESS_GROUPS_ONE",
	"RESOURCES",
	"SERVICES",
	"SERVICE_VERSATILE",
	"SESSION_METRICS",
	"SYNTHETIC_HTTP_MONITOR",
	"SYNTHETIC_SINGLE_EXT_TEST",
	"SYNTHETIC_SINGLE_WEBCHECK",
	"SYNTHETIC_TESTS",
	"THIRD_PARTY_MOST_ACTIVE",
	"UEM_CONVERSIONS_OVERALL",
	"UEM_CONVERSIONS_PER_GOAL",
	"UEM_JSERRORS_OVERALL",
	"UEM_KEY_USER_ACTIONS",
	"USERS",
	"VIRTUALIZATION",
	"HEADER",
}

// UnmarshalJSON performs custom unmarshalling of this enum type
func (t *TileType) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	*t = TileType(name)
	return nil
}
