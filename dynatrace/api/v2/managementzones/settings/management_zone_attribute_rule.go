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

package managementzones

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// No documentation available
type ManagementZoneAttributeRule struct {
	PgToServicePropagation                     *bool                `json:"pgToServicePropagation,omitempty"`                     // Apply to all services provided by the process groups
	EntityType                                 ManagementZoneMeType `json:"entityType"`                                           // Rule applies to
	ServiceToHostPropagation                   *bool                `json:"serviceToHostPropagation,omitempty"`                   // Apply to underlying hosts of matching services
	PgToHostPropagation                        *bool                `json:"pgToHostPropagation,omitempty"`                        // Apply to underlying hosts of matching process groups
	AzureToPGPropagation                       *bool                `json:"azureToPGPropagation,omitempty"`                       // Apply to process groups connected to matching Azure entities
	Conditions                                 AttributeConditions  `json:"conditions"`                                           // Conditions
	ServiceToPGPropagation                     *bool                `json:"serviceToPGPropagation,omitempty"`                     // Apply to underlying process groups of matching services
	HostToPGPropagation                        *bool                `json:"hostToPGPropagation,omitempty"`                        // Apply to processes running on matching hosts
	CustomDeviceGroupToCustomDevicePropagation *bool                `json:"customDeviceGroupToCustomDevicePropagation,omitempty"` // Apply to custom devices in a custom device group
	AzureToServicePropagation                  *bool                `json:"azureToServicePropagation,omitempty"`                  // Apply to services provided by matching Azure entities
}

func (me *ManagementZoneAttributeRule) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"pg_to_service_propagation": {
			Type:        schema.TypeBool,
			Description: "Apply to all services provided by the process groups",
			Optional:    true,
		},
		"entity_type": {
			Type:        schema.TypeString,
			Description: "Rule applies to",
			Required:    true,
		},
		"service_to_host_propagation": {
			Type:        schema.TypeBool,
			Description: "Apply to underlying hosts of matching services",
			Optional:    true,
		},
		"pg_to_host_propagation": {
			Type:        schema.TypeBool,
			Description: "Apply to underlying hosts of matching process groups",
			Optional:    true,
		},
		"azure_to_pgpropagation": {
			Type:        schema.TypeBool,
			Description: "Apply to process groups connected to matching Azure entities",
			Optional:    true,
		},
		"attribute_conditions": {
			Type:        schema.TypeList,
			Description: "Conditions",
			MaxItems:    1,
			MinItems:    1,
			Elem:        &schema.Resource{Schema: new(AttributeConditions).Schema()},
			Required:    true,
		},
		"service_to_pgpropagation": {
			Type:        schema.TypeBool,
			Description: "Apply to underlying process groups of matching services",
			Optional:    true,
		},
		"host_to_pgpropagation": {
			Type:        schema.TypeBool,
			Description: "Apply to processes running on matching hosts",
			Optional:    true,
		},
		"custom_device_group_to_custom_device_propagation": {
			Type:        schema.TypeBool,
			Description: "Apply to custom devices in a custom device group",
			Optional:    true,
		},
		"azure_to_service_propagation": {
			Type:        schema.TypeBool,
			Description: "Apply to services provided by matching Azure entities",
			Optional:    true,
		},
	}
}

func (me *ManagementZoneAttributeRule) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]any{
		"pg_to_service_propagation":                        me.PgToServicePropagation,
		"entity_type":                                      me.EntityType,
		"service_to_host_propagation":                      me.ServiceToHostPropagation,
		"pg_to_host_propagation":                           me.PgToHostPropagation,
		"azure_to_pgpropagation":                           me.AzureToPGPropagation,
		"attribute_conditions":                             me.Conditions,
		"service_to_pgpropagation":                         me.ServiceToPGPropagation,
		"host_to_pgpropagation":                            me.HostToPGPropagation,
		"custom_device_group_to_custom_device_propagation": me.CustomDeviceGroupToCustomDevicePropagation,
		"azure_to_service_propagation":                     me.AzureToServicePropagation,
	})
}

func (me *ManagementZoneAttributeRule) UnmarshalHCL(decoder hcl.Decoder) error {
	err := decoder.DecodeAll(map[string]any{
		"pg_to_service_propagation":                        &me.PgToServicePropagation,
		"entity_type":                                      &me.EntityType,
		"service_to_host_propagation":                      &me.ServiceToHostPropagation,
		"pg_to_host_propagation":                           &me.PgToHostPropagation,
		"azure_to_pgpropagation":                           &me.AzureToPGPropagation,
		"attribute_conditions":                             &me.Conditions,
		"service_to_pgpropagation":                         &me.ServiceToPGPropagation,
		"host_to_pgpropagation":                            &me.HostToPGPropagation,
		"custom_device_group_to_custom_device_propagation": &me.CustomDeviceGroupToCustomDevicePropagation,
		"azure_to_service_propagation":                     &me.AzureToServicePropagation,
	})
	if me.PgToServicePropagation == nil && me.EntityType == ManagementZoneMeTypes.ProcessGroup {
		me.PgToServicePropagation = opt.NewBool(false)
	}
	if me.ServiceToHostPropagation == nil && me.EntityType == ManagementZoneMeTypes.Service {
		me.ServiceToHostPropagation = opt.NewBool(false)
	}
	if me.PgToHostPropagation == nil && me.EntityType == ManagementZoneMeTypes.ProcessGroup {
		me.PgToHostPropagation = opt.NewBool(false)
	}
	if me.AzureToPGPropagation == nil && me.EntityType == ManagementZoneMeTypes.Azure {
		me.AzureToPGPropagation = opt.NewBool(false)
	}
	if me.ServiceToPGPropagation == nil && me.EntityType == ManagementZoneMeTypes.Service {
		me.ServiceToPGPropagation = opt.NewBool(false)
	}
	if me.HostToPGPropagation == nil && me.EntityType == ManagementZoneMeTypes.Host {
		me.HostToPGPropagation = opt.NewBool(false)
	}
	if me.CustomDeviceGroupToCustomDevicePropagation == nil && me.EntityType == ManagementZoneMeTypes.CustomDeviceGroup {
		me.CustomDeviceGroupToCustomDevicePropagation = opt.NewBool(false)
	}
	if me.AzureToServicePropagation == nil && me.EntityType == ManagementZoneMeTypes.Azure {
		me.AzureToServicePropagation = opt.NewBool(false)
	}
	return err
}
