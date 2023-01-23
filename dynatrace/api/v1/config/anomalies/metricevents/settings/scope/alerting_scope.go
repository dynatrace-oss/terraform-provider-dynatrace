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

package scope

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/xjson"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AlertingScope A single filter for the alerting scope.
// This is the base version of the filter, depending on the type,
// the actual JSON may contain additional fields.
type AlertingScope interface {
	GetType() FilterType
}

// BaseAlertingScope A single filter for the alerting scope.
// This is the base version of the filter, depending on the type,
// the actual JSON may contain additional fields.
type BaseAlertingScope struct {
	FilterType FilterType                 `json:"filterType"` // Defines the actual set of fields depending on the value. See one of the following objects:  * `ENTITY_ID` -> EntityIdAlertingScope  * `MANAGEMENT_ZONE` -> ManagementZoneAlertingScope  * `TAG` -> TagFilterAlertingScope  * `NAME` -> NameAlertingScope  * `CUSTOM_DEVICE_GROUP_NAME` -> CustomDeviceGroupNameAlertingScope  * `HOST_GROUP_NAME` -> HostGroupNameAlertingScope  * `HOST_NAME` -> HostNameAlertingScope  * `PROCESS_GROUP_ID` -> ProcessGroupIdAlertingScope  * `PROCESS_GROUP_NAME` -> ProcessGroupNameAlertingScope
	Unknowns   map[string]json.RawMessage `json:"-"`
}

func (me *BaseAlertingScope) GetType() FilterType {
	return me.FilterType
}

func (me *BaseAlertingScope) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Defines the actual set of fields depending on the value",
		},
		"unknowns": {
			Type:        schema.TypeString,
			Description: "allows for configuring properties that are not explicitly supported by the current version of this provider",
			Optional:    true,
		},
	}
}

func (me *BaseAlertingScope) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(me.Unknowns); err != nil {
		return err
	}
	return properties.Encode("type", me.FilterType)
}

func (me *BaseAlertingScope) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), me); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &me.Unknowns); err != nil {
			return err
		}
		delete(me.Unknowns, "filterType")

		if len(me.Unknowns) == 0 {
			me.Unknowns = nil
		}
	}
	if value, ok := decoder.GetOk("type"); ok {
		me.FilterType = FilterType(value.(string))
	}
	return nil
}

func (me *BaseAlertingScope) MarshalJSON() ([]byte, error) {
	properties := xjson.NewProperties(me.Unknowns)
	if err := properties.MarshalAll(map[string]any{
		"filterType": me.FilterType,
	}); err != nil {
		return nil, err
	}
	return json.Marshal(properties)
}

func (me *BaseAlertingScope) UnmarshalJSON(data []byte) error {
	properties := xjson.NewProperties(me.Unknowns)
	if err := json.Unmarshal(data, &properties); err != nil {
		return err
	}
	if err := properties.UnmarshalAll(map[string]any{
		"filterType": &me.FilterType,
	}); err != nil {
		return err
	}
	return nil
}

// FilterType Defines the actual set of fields depending on the value. See one of the following objects:
// * `ENTITY_ID` -> EntityIdAlertingScope
// * `MANAGEMENT_ZONE` -> ManagementZoneAlertingScope
// * `TAG` -> TagFilterAlertingScope
// * `NAME` -> NameAlertingScope
// * `CUSTOM_DEVICE_GROUP_NAME` -> CustomDeviceGroupNameAlertingScope
// * `HOST_GROUP_NAME` -> HostGroupNameAlertingScope
// * `HOST_NAME` -> HostNameAlertingScope
// * `PROCESS_GROUP_ID` -> ProcessGroupIdAlertingScope
// * `PROCESS_GROUP_NAME` -> ProcessGroupNameAlertingScope
type FilterType string

// FilterTypes offers the known enum values
var FilterTypes = struct {
	CustomDeviceGroupName FilterType
	EntityID              FilterType
	HostGroupName         FilterType
	HostName              FilterType
	ManagementZone        FilterType
	Name                  FilterType
	ProcessGroupID        FilterType
	ProcessGroupName      FilterType
	Tag                   FilterType
}{
	"CUSTOM_DEVICE_GROUP_NAME",
	"ENTITY_ID",
	"HOST_GROUP_NAME",
	"HOST_NAME",
	"MANAGEMENT_ZONE",
	"NAME",
	"PROCESS_GROUP_ID",
	"PROCESS_GROUP_NAME",
	"TAG",
}
