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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AlertingScopes []AlertingScope

func (me AlertingScopes) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"custom_device_group_name": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for the related custom device group name",
			Elem:        &schema.Resource{Schema: new(CustomDeviceGroupName).Schema()},
		},
		"entity": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for a monitored entity identifier",
			Elem:        &schema.Resource{Schema: new(EntityID).Schema()},
		},
		"host_group_name": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for the related host group name",
			Elem:        &schema.Resource{Schema: new(HostGroupName).Schema()},
		},
		"host_name": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for the related host name",
			Elem:        &schema.Resource{Schema: new(HostName).Schema()},
		},
		"management_zone": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for a management zone identifier",
			Elem:        &schema.Resource{Schema: new(ManagementZone).Schema()},
		},
		"name": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for a monitored entity name",
			Elem:        &schema.Resource{Schema: new(Name).Schema()},
		},
		"process_group_id": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for a process group identifier",
			Elem:        &schema.Resource{Schema: new(ProcessGroupID).Schema()},
		},
		"process_group_name": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for the related process group name",
			Elem:        &schema.Resource{Schema: new(ProcessGroupName).Schema()},
		},
		"tag": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A scope filter for tags on entities",
			Elem:        &schema.Resource{Schema: new(TagFilter).Schema()},
		},
		"scope": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A generic scope filter",
			Elem:        &schema.Resource{Schema: new(BaseAlertingScope).Schema()},
		},
	}
}

func (me AlertingScopes) MarshalHCL(properties hcl.Properties) error {
	customDeviceGroupNames := []any{}
	entityIDs := []any{}
	hostGroupNames := []any{}
	hostNames := []any{}
	managementZones := []any{}
	names := []any{}
	processGroupIDs := []any{}
	processGroupNames := []any{}
	scopes := []any{}
	tags := []any{}
	for _, scope := range me {
		switch sc := scope.(type) {
		case *CustomDeviceGroupName:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				customDeviceGroupNames = append(customDeviceGroupNames, marshalled)
			} else {
				return err
			}
		case *EntityID:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				entityIDs = append(entityIDs, marshalled)
			} else {
				return nil
			}
		case *HostGroupName:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				hostGroupNames = append(hostGroupNames, marshalled)
			} else {
				return nil
			}
		case *HostName:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				hostNames = append(hostNames, marshalled)
			} else {
				return nil
			}
		case *ManagementZone:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				managementZones = append(managementZones, marshalled)
			} else {
				return nil
			}
		case *Name:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				names = append(names, marshalled)
			} else {
				return nil
			}
		case *ProcessGroupID:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				processGroupIDs = append(processGroupIDs, marshalled)
			} else {
				return nil
			}
		case *ProcessGroupName:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				processGroupNames = append(processGroupNames, marshalled)
			} else {
				return nil
			}
		case *TagFilter:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				tags = append(tags, marshalled)
			} else {
				return nil
			}
		case *BaseAlertingScope:
			marshalled := hcl.Properties{}
			if err := sc.MarshalHCL(marshalled); err == nil {
				scopes = append(scopes, marshalled)
			} else {
				return nil
			}
		default:
		}
		if len(customDeviceGroupNames) > 0 {
			properties["custom_device_group_name"] = customDeviceGroupNames
		} else {
			properties["custom_device_group_name"] = nil
		}
		if len(entityIDs) > 0 {
			properties["entity"] = entityIDs
		} else {
			properties["entity"] = nil
		}
		if len(hostGroupNames) > 0 {
			properties["host_group_name"] = hostGroupNames
		} else {
			properties["host_group_name"] = nil
		}
		if len(hostNames) > 0 {
			properties["host_name"] = hostNames
		} else {
			properties["host_name"] = nil
		}
		if len(managementZones) > 0 {
			properties["management_zone"] = managementZones
		} else {
			properties["management_zone"] = nil
		}
		if len(names) > 0 {
			properties["name"] = names
		} else {
			properties["name"] = nil
		}
		if len(processGroupIDs) > 0 {
			properties["process_group_id"] = processGroupIDs
		} else {
			properties["process_group_id"] = nil
		}
		if len(processGroupNames) > 0 {
			properties["process_group_name"] = processGroupNames
		} else {
			properties["process_group_name"] = nil
		}
		if len(tags) > 0 {
			properties["tag"] = tags
		} else {
			properties["tag"] = nil
		}
		if len(scopes) > 0 {
			properties["scope"] = scopes
		} else {
			properties["scope"] = nil
		}
	}
	return nil
}

func (me *AlertingScopes) UnmarshalHCL(decoder hcl.Decoder) error {
	nme := AlertingScopes{}
	if result, ok := decoder.GetOk("custom_device_group_name.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(CustomDeviceGroupName)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_device_group_name", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("entity.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(EntityID)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "entity", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("host_group_name.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(HostGroupName)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "host_group_name", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("host_name.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(HostName)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "host_name", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("management_zone.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(ManagementZone)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "management_zone", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("name.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(Name)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "name", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("process_group_id.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(ProcessGroupID)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "process_group_id", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("process_group_name.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(ProcessGroupName)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "process_group_name", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("tag.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(TagFilter)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "tag", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	if result, ok := decoder.GetOk("scope.#"); ok {
		for idx := 0; idx < result.(int); idx++ {
			entry := new(BaseAlertingScope)
			if err := entry.UnmarshalHCL(hcl.NewDecoder(decoder, "scope", idx)); err != nil {
				return err
			}
			nme = append(nme, entry)
		}
	}
	*me = nme
	return nil
}

func (me *AlertingScopes) UnmarshalJSON(data []byte) error {
	dims := AlertingScopes{}
	rawMessages := []json.RawMessage{}
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return err
	}
	for _, rawMessage := range rawMessages {
		properties := map[string]json.RawMessage{}
		if err := json.Unmarshal(rawMessage, &properties); err != nil {
			return err
		}
		if rawFilterType, found := properties["filterType"]; found {
			var sFilterType string
			if err := json.Unmarshal(rawFilterType, &sFilterType); err != nil {
				return err
			}
			switch sFilterType {
			case string(FilterTypes.CustomDeviceGroupName):
				cfg := new(CustomDeviceGroupName)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.EntityID):
				cfg := new(EntityID)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.HostGroupName):
				cfg := new(HostGroupName)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.HostName):
				cfg := new(HostName)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.ManagementZone):
				cfg := new(ManagementZone)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.Name):
				cfg := new(Name)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.ProcessGroupID):
				cfg := new(ProcessGroupID)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.ProcessGroupName):
				cfg := new(ProcessGroupName)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			case string(FilterTypes.Tag):
				cfg := new(TagFilter)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			default:
				cfg := new(BaseAlertingScope)
				if err := json.Unmarshal(rawMessage, &cfg); err != nil {
					return err
				}
				dims = append(dims, cfg)
			}
		}
		*me = dims
	}
	return nil
}
