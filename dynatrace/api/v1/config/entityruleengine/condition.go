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

package entityruleengine

import (
	"encoding/json"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/comparison"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/entityruleengine/condition"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Condition A condition defines how to execute matching logic for an entity.
type Condition struct {
	ComparisonInfo comparison.Comparison      `json:"comparisonInfo"` // Defines how the matching is actually performed: what and how are we comparing.  The actual set of fields and possible values of the **operator** field depend on the **type** of the comparison. \n\nFind the list of actual models in the description of the **type** field and check the description of the model you need.
	Key            condition.Key              `json:"key"`            // The key to identify the data we're matching.  The actual set of fields and possible values vary, depending on the **type** of the key.  Find the list of actual objects in the description of the **type** field.
	Unknowns       map[string]json.RawMessage `json:"-"`
}

type Validator interface {
	Validate() []string
}

func (erec *Condition) Validate() []string {
	if erec.ComparisonInfo == nil {
		return []string{}
	}
	if validator, ok := erec.ComparisonInfo.(Validator); ok {
		return validator.Validate()
	}
	return []string{}
}

func rmType(m map[string]*schema.Schema) map[string]*schema.Schema {
	delete(m, "type")
	return m
}

func (erec *Condition) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application_type_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `APPLICATION_TYPE` attributes",
			Optional:    true,
			Deprecated:  "You should use 'application_type' instead of 'application_type_comparison'. This attribute still exists for backwards compatibility.",
			Elem: &schema.Resource{
				Schema: new(comparison.ApplicationType).Schema(),
			},
		},
		"application_type": {
			Type:        schema.TypeList,
			Description: "Comparison for `APPLICATION_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.ApplicationType).Schema()),
			},
		},
		"azure_compute_mode_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `AZURE_COMPUTE_MODE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.AzureComputeMode).Schema()),
			},
		},
		"azure_compute_mode": {
			Type:        schema.TypeList,
			Description: "Comparison for `AZURE_COMPUTE_MODE` attributes",
			Deprecated:  "You should use 'azure_compute_mode' instead of 'azure_compute_mode_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.AzureComputeMode).Schema()),
			},
		},
		"azure_sku_comparision": {
			Type:        schema.TypeList,
			Description: "Comparison for `AZURE_SKU` attributes",
			Optional:    true,
			Deprecated:  "You should use 'azure_sku' instead of 'azure_sku_comparision'. This attribute still exists for backwards compatibility.",
			Elem: &schema.Resource{
				Schema: new(comparison.AzureSku).Schema(),
			},
		},
		"azure_sku": {
			Type:        schema.TypeList,
			Description: "Comparison for `AZURE_SKU` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.AzureSku).Schema()),
			},
		},
		"base_comparison_basic": {
			Type:        schema.TypeList,
			Description: "A comparison that's yet unknown to the provider. Operator and Value need to be encoded using the 'unknowns' property.",
			Deprecated:  "You should use 'comparison' instead of 'base_comparison_basic'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.BaseComparison).Schema(),
			},
		},
		"comparison": {
			Type:        schema.TypeList,
			Description: "A comparison that's yet unknown to the provider. Operator and Value need to be encoded using the 'unknowns' property.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.BaseComparison).Schema(),
			},
		},
		"bitness_comparision": {
			Type:        schema.TypeList,
			Description: "Comparison for `BITNESS` attributes",
			Optional:    true,
			Deprecated:  "You should use 'bitness' instead of 'bitness_comparision'. This attribute still exists for backwards compatibility.",
			Elem: &schema.Resource{
				Schema: new(comparison.Bitness).Schema(),
			},
		},
		"bitness": {
			Type:        schema.TypeList,
			Description: "Comparison for `BITNESS` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.Bitness).Schema()),
			},
		},
		"cloud_type_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `CLOUD_TYPE` attributes",
			Deprecated:  "You should use 'cloud_type' instead of 'cloud_type_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.CloudType).Schema(),
			},
		},
		"cloud_type": {
			Type:        schema.TypeList,
			Description: "Comparison for `CLOUD_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.CloudType).Schema()),
			},
		},
		"custom_application_type_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `CUSTOM_APPLICATION_TYPE` attributes",
			Deprecated:  "You should use 'custom_application_type' instead of 'custom_application_type_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.CustomApplicationType).Schema(),
			},
		},
		"custom_application_type": {
			Type:        schema.TypeList,
			Description: "Comparison for `CUSTOM_APPLICATION_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.CustomApplicationType).Schema()),
			},
		},
		"database_topology_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `DATABASE_TOPOLOGY` attributes",
			Deprecated:  "You should use 'database_topology' instead of 'database_topology_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.DatabaseTopology).Schema(),
			},
		},
		"database_topology": {
			Type:        schema.TypeList,
			Description: "Comparison for `DATABASE_TOPOLOGY` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.DatabaseTopology).Schema()),
			},
		},
		"dcrum_decoder_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `DCRUM_DECODER_TYPE` attributes",
			Deprecated:  "You should use 'dcrum_decoder' instead of 'dcrum_decoder_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.DCRumDecoder).Schema(),
			},
		},
		"dcrum_decoder": {
			Type:        schema.TypeList,
			Description: "Comparison for `DCRUM_DECODER_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.DCRumDecoder).Schema()),
			},
		},
		"entity_id_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `ENTITY_ID` attributes",
			Deprecated:  "You should use 'entity' instead of 'entity_id_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.EntityID).Schema(),
			},
		},
		"entity": {
			Type:        schema.TypeList,
			Description: "Comparison for `ENTITY_ID` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.EntityID).Schema()),
			},
		},
		"hypervisor_type_comparision": {
			Type:        schema.TypeList,
			Description: "`hypervisor_type_comparision` is deprecated. Use `hypervisor` instead",
			Deprecated:  "`hypervisor_type_comparision` is deprecated. Use `hypervisor` instead",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.HypervisorType).Schema(),
			},
		},
		"hypervisor": {
			Type:        schema.TypeList,
			Description: "Comparison for `HYPERVISOR_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.HypervisorType).Schema()),
			},
		},
		"indexed_name_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `INDEXED_NAME` attributes",
			Deprecated:  "You should use 'indexed_name' instead of 'indexed_name_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.IndexedName).Schema(),
			},
		},
		"indexed_name": {
			Type:        schema.TypeList,
			Description: "Comparison for `INDEXED_NAME` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.IndexedName).Schema()),
			},
		},
		"indexed_string_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `INDEXED_STRING` attributes",
			Deprecated:  "You should use 'indexed_string' instead of 'indexed_string_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.IndexedString).Schema(),
			},
		},
		"indexed_string": {
			Type:        schema.TypeList,
			Description: "Comparison for `INDEXED_STRING` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.IndexedString).Schema()),
			},
		},
		"indexed_tag_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `INDEXED_TAG` attributes",
			Deprecated:  "You should use 'indexed_tag' instead of 'indexed_tag_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.IndexedTag).Schema(),
			},
		},
		"indexed_tag": {
			Type:        schema.TypeList,
			Description: "Comparison for `INDEXED_TAG` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.IndexedTag).Schema()),
			},
		},
		"integer_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `INTEGER` attributes",
			Deprecated:  "You should use 'integer' instead of 'integer_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.Integer).Schema(),
			},
		},
		"integer": {
			Type:        schema.TypeList,
			Description: "Comparison for `INTEGER` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.Integer).Schema()),
			},
		},
		"ipaddress_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `IP_ADDRESS` attributes",
			Deprecated:  "You should use 'ipaddress' instead of 'ipaddress_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.IPAddress).Schema(),
			},
		},
		"ipaddress": {
			Type:        schema.TypeList,
			Description: "Comparison for `IP_ADDRESS` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.IPAddress).Schema()),
			},
		},
		"mobile_platform_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `MOBILE_PLATFORM` attributes",
			Deprecated:  "You should use 'mobile_platform' instead of 'mobile_platform_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.MobilePlatform).Schema(),
			},
		},
		"mobile_platform": {
			Type:        schema.TypeList,
			Description: "Comparison for `MOBILE_PLATFORM` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.MobilePlatform).Schema()),
			},
		},
		"osarchitecture_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `OS_ARCHITECTURE` attributes",
			Deprecated:  "You should use 'os_arch' instead of 'osarchitecture_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.OSArchitecture).Schema(),
			},
		},
		"os_arch": {
			Type:        schema.TypeList,
			Description: "Comparison for `OS_ARCHITECTURE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.OSArchitecture).Schema()),
			},
		},
		"ostype_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `OS_TYPE` attributes",
			Deprecated:  "You should use 'os_type' instead of 'ostype_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.OSType).Schema(),
			},
		},
		"os_type": {
			Type:        schema.TypeList,
			Description: "Comparison for `OS_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.OSType).Schema()),
			},
		},
		"paas_type_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `PAAS_TYPE` attributes",
			Deprecated:  "You should use 'paas_type' instead of 'paas_type_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.PaasType).Schema(),
			},
		},
		"paas_type": {
			Type:        schema.TypeList,
			Description: "Comparison for `PAAS_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.PaasType).Schema()),
			},
		},
		"service_topology_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `SERVICE_TOPOLOGY` attributes",
			Deprecated:  "You should use 'service_topology' instead of 'service_topology_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.ServiceTopology).Schema(),
			},
		},
		"service_topology": {
			Type:        schema.TypeList,
			Description: "Comparison for `SERVICE_TOPOLOGY` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.ServiceTopology).Schema()),
			},
		},
		"service_type_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `SERVICE_TYPE` attributes",
			Deprecated:  "You should use 'service_type' instead of 'service_type_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.ServiceType).Schema(),
			},
		},
		"service_type": {
			Type:        schema.TypeList,
			Description: "Comparison for `SERVICE_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.ServiceType).Schema()),
			},
		},
		"simple_host_tech_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `SIMPLE_HOST_TECH` attributes",
			Deprecated:  "You should use 'host_tech' instead of 'simple_host_tech_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.SimpleHostTech).Schema(),
			},
		},
		"host_tech": {
			Type:        schema.TypeList,
			Description: "Comparison for `SIMPLE_HOST_TECH` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.SimpleHostTech).Schema()),
			},
		},
		"simple_tech_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `SIMPLE_TECH` attributes",
			Deprecated:  "You should use 'tech' instead of 'simple_tech_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.SimpleTech).Schema(),
			},
		},
		"tech": {
			Type:        schema.TypeList,
			Description: "Comparison for `SIMPLE_TECH` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.SimpleTech).Schema()),
			},
		},
		"string_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `STRING` attributes",
			Deprecated:  "You should use 'string' instead of 'string_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.String).Schema(),
			},
		},
		"string": {
			Type:        schema.TypeList,
			Description: "Comparison for `STRING` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.String).Schema()),
			},
		},
		"synthetic_engine_type_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `SYNTHETIC_ENGINE_TYPE` attributes",
			Deprecated:  "You should use 'synthetic_engine' instead of 'synthetic_engine_type_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.SyntheticEngineType).Schema(),
			},
		},
		"synthetic_engine": {
			Type:        schema.TypeList,
			Description: "Comparison for `SYNTHETIC_ENGINE_TYPE` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.SyntheticEngineType).Schema()),
			},
		},
		"tag_comparison": {
			Type:        schema.TypeList,
			Description: "Comparison for `TAG` attributes",
			Deprecated:  "You should use 'tag' instead of 'tag_comparison'. This attribute still exists for backwards compatibility.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(comparison.Tag).Schema(),
			},
		},
		"tag": {
			Type:        schema.TypeList,
			Description: "Comparison for `TAG` attributes",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(comparison.Tag).Schema()),
			},
		},
		"base_condition_key": {
			Type:        schema.TypeList,
			Description: "Fallback for not yet known type",
			Deprecated:  "'base_condition_key' is deprecated. You should use 'key'",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(condition.BaseConditionKey).Schema(),
			},
		},
		"key": {
			Type:        schema.TypeList,
			Description: "Fallback for not yet known type",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(condition.BaseConditionKey).Schema(),
			},
		},
		"custom_host_metadata_condition_key": {
			Type:        schema.TypeList,
			Description: "Key for Custom Host Metadata",
			Deprecated:  "'custom_host_metadata_condition_key' is deprecated. You should use 'custom_host_metadata'",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(condition.CustomHostMetadata).Schema(),
			},
		},
		"custom_host_metadata": {
			Type:        schema.TypeList,
			Description: "Key for Custom Host Metadata",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(condition.CustomHostMetadata).Schema()),
			},
		},
		"custom_process_metadata_condition_key": {
			Type:        schema.TypeList,
			Description: "Key for Custom Process Metadata",
			Deprecated:  "'custom_process_metadata_condition_key' is deprecated. You should use 'custom_process_metadata'",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(condition.CustomProcessMetadata).Schema(),
			},
		},
		"custom_process_metadata": {
			Type:        schema.TypeList,
			Description: "Key for Custom Process Metadata",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(condition.CustomProcessMetadata).Schema()),
			},
		},
		"process_metadata_condition_key": {
			Type:        schema.TypeList,
			Description: "The key for dynamic attributes of the `PROCESS_PREDEFINED_METADATA_KEY` type",
			Deprecated:  "'process_metadata_condition_key' is deprecated. You should use 'process_metadata'",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(condition.ProcessMetadata).Schema(),
			},
		},
		"process_metadata": {
			Type:        schema.TypeList,
			Description: "The key for dynamic attributes of the `PROCESS_PREDEFINED_METADATA_KEY` type",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(condition.ProcessMetadata).Schema()),
			},
		},
		"string_condition_key": {
			Type:        schema.TypeList,
			Description: " The key for dynamic attributes of the `STRING` type",
			Deprecated:  "'string_condition_key' is deprecated. You should use 'string_key'",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: new(condition.String).Schema(),
			},
		},
		"string_key": {
			Type:        schema.TypeList,
			Description: " The key for dynamic attributes of the `STRING` type",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: rmType(new(condition.String).Schema()),
			},
		},

		"unknowns": {
			Type:        schema.TypeString,
			Description: "Any attributes that aren't yet supported by this provider",
			Optional:    true,
		},
	}
}

func (erec *Condition) MarshalHCL(properties hcl.Properties) error {
	if err := properties.Unknowns(erec.Unknowns); err != nil {
		return err
	}

	switch comparison := erec.ComparisonInfo.(type) {
	case *comparison.CustomApplicationType:
		if err := properties.Encode("custom_application_type", comparison); err != nil {
			return err
		}
	case *comparison.MobilePlatform:
		if err := properties.Encode("mobile_platform", comparison); err != nil {
			return err
		}
	case *comparison.ApplicationType:
		if err := properties.Encode("application_type", comparison); err != nil {
			return err
		}
	case *comparison.Bitness:
		if err := properties.Encode("bitness", comparison); err != nil {
			return err
		}
	case *comparison.PaasType:
		if err := properties.Encode("paas_type", comparison); err != nil {
			return err
		}
	case *comparison.OSArchitecture:
		if err := properties.Encode("os_arch", comparison); err != nil {
			return err
		}
	case *comparison.ServiceTopology:
		if err := properties.Encode("service_topology", comparison); err != nil {
			return err
		}
	case *comparison.String:
		if err := properties.Encode("string", comparison); err != nil {
			return err
		}
	case *comparison.DatabaseTopology:
		if err := properties.Encode("database_topology", comparison); err != nil {
			return err
		}
	case *comparison.DCRumDecoder:
		if err := properties.Encode("dcrum_decoder", comparison); err != nil {
			return err
		}
	case *comparison.IndexedTag:
		if err := properties.Encode("indexed_tag", comparison); err != nil {
			return err
		}
	case *comparison.Tag:
		if err := properties.Encode("tag", comparison); err != nil {
			return err
		}
	case *comparison.HypervisorType:
		if err := properties.Encode("hypervisor", comparison); err != nil {
			return err
		}
	case *comparison.CloudType:
		if err := properties.Encode("cloud_type", comparison); err != nil {
			return err
		}
	case *comparison.IndexedName:
		if err := properties.Encode("indexed_name", comparison); err != nil {
			return err
		}
	case *comparison.IndexedString:
		if err := properties.Encode("indexed_string", comparison); err != nil {
			return err
		}
	case *comparison.SimpleTech:
		if err := properties.Encode("tech", comparison); err != nil {
			return err
		}
	case *comparison.AzureSku:
		if err := properties.Encode("azure_sku", comparison); err != nil {
			return err
		}
	case *comparison.EntityID:
		if err := properties.Encode("entity", comparison); err != nil {
			return err
		}
	case *comparison.SimpleHostTech:
		if err := properties.Encode("host_tech", comparison); err != nil {
			return err
		}
	case *comparison.Integer:
		if err := properties.Encode("integer", comparison); err != nil {
			return err
		}
	case *comparison.ServiceType:
		if err := properties.Encode("service_type", comparison); err != nil {
			return err
		}
	case *comparison.OSType:
		if err := properties.Encode("os_type", comparison); err != nil {
			return err
		}
	case *comparison.SyntheticEngineType:
		if err := properties.Encode("synthetic_engine", comparison); err != nil {
			return err
		}
	case *comparison.IPAddress:
		if err := properties.Encode("ipaddress", comparison); err != nil {
			return err
		}
	case *comparison.AzureComputeMode:
		if err := properties.Encode("azure_compute_mode", comparison); err != nil {
			return err
		}
	default:
		if err := properties.Encode("comparison", comparison); err != nil {
			return err
		}
	}

	switch key := erec.Key.(type) {
	case *condition.CustomHostMetadata:
		if err := properties.Encode("custom_host_metadata", key); err != nil {
			return err
		}
	case *condition.String:
		if err := properties.Encode("string_key", key); err != nil {
			return err
		}
	case *condition.CustomProcessMetadata:
		if err := properties.Encode("custom_process_metadata", key); err != nil {
			return err
		}
	case *condition.ProcessMetadata:
		if err := properties.Encode("process_metadata", key); err != nil {
			return err
		}
	default:
		if err := properties.Encode("key", key); err != nil {
			return err
		}
	}
	return nil
}

func (erec *Condition) UnmarshalHCL(decoder hcl.Decoder) error {
	if value, ok := decoder.GetOk("unknowns"); ok {
		if err := json.Unmarshal([]byte(value.(string)), erec); err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(value.(string)), &erec.Unknowns); err != nil {
			return err
		}
		delete(erec.Unknowns, "hypervisor")
		delete(erec.Unknowns, "integer")
		delete(erec.Unknowns, "service_type")
		delete(erec.Unknowns, "string")
		delete(erec.Unknowns, "synthetic_engine")
		delete(erec.Unknowns, "string_key")
		delete(erec.Unknowns, "azure_compute_mode")
		delete(erec.Unknowns, "bitness")
		delete(erec.Unknowns, "cloud_type_comparison")
		delete(erec.Unknowns, "entity_id_comparison")
		delete(erec.Unknowns, "mobile_platform")
		delete(erec.Unknowns, "custom_application_type")
		delete(erec.Unknowns, "indexed_tag")
		delete(erec.Unknowns, "paas_type")
		delete(erec.Unknowns, "base_condition_key")
		delete(erec.Unknowns, "azure_compute_mode_comparison")
		delete(erec.Unknowns, "bitness_comparision")
		delete(erec.Unknowns, "simple_tech_comparison")
		delete(erec.Unknowns, "os_arch")
		delete(erec.Unknowns, "service_topology_comparison")
		delete(erec.Unknowns, "service_type_comparison")
		delete(erec.Unknowns, "application_type")
		delete(erec.Unknowns, "cloud_type")
		delete(erec.Unknowns, "database_topology")
		delete(erec.Unknowns, "indexed_tag_comparison")
		delete(erec.Unknowns, "integer_comparison")
		delete(erec.Unknowns, "custom_application_type_comparison")
		delete(erec.Unknowns, "azure_sku")
		delete(erec.Unknowns, "entity")
		delete(erec.Unknowns, "mobile_platform_comparison")
		delete(erec.Unknowns, "custom_host_metadata")
		delete(erec.Unknowns, "process_metadata")
		delete(erec.Unknowns, "custom_process_metadata")
		delete(erec.Unknowns, "database_topology_comparison")
		delete(erec.Unknowns, "indexed_name_comparison")
		delete(erec.Unknowns, "indexed_name")
		delete(erec.Unknowns, "host_tech")
		delete(erec.Unknowns, "synthetic_engine_type_comparison")
		delete(erec.Unknowns, "string_condition_key")
		delete(erec.Unknowns, "application_type_comparison")
		delete(erec.Unknowns, "hypervisor_type_comparision")
		delete(erec.Unknowns, "indexed_string")
		delete(erec.Unknowns, "ipaddress_comparison")
		delete(erec.Unknowns, "ostype_comparison")
		delete(erec.Unknowns, "dcrum_decoder")
		delete(erec.Unknowns, "unknowns")
		delete(erec.Unknowns, "azure_sku_comparision")
		delete(erec.Unknowns, "comparison")
		delete(erec.Unknowns, "osarchitecture_comparison")
		delete(erec.Unknowns, "simple_host_tech_comparison")
		delete(erec.Unknowns, "tag_comparison")
		delete(erec.Unknowns, "ipaddress")
		delete(erec.Unknowns, "os_type")
		delete(erec.Unknowns, "tech")
		delete(erec.Unknowns, "indexed_string_comparison")
		delete(erec.Unknowns, "paas_type_comparison")
		delete(erec.Unknowns, "key")
		delete(erec.Unknowns, "custom_host_metadata_condition_key")
		delete(erec.Unknowns, "custom_process_metadata_condition_key")
		delete(erec.Unknowns, "tag")
		delete(erec.Unknowns, "base_comparison_basic")
		delete(erec.Unknowns, "dcrum_decoder_comparison")
		delete(erec.Unknowns, "service_topology")
		delete(erec.Unknowns, "string_comparison")
		delete(erec.Unknowns, "process_metadata_condition_key")
		if len(erec.Unknowns) == 0 {
			erec.Unknowns = nil
		}
	}

	if _, ok := decoder.GetOk("application_type_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.ApplicationType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "application_type_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("application_type.#"); ok {
		erec.ComparisonInfo = new(comparison.ApplicationType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "application_type", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("azure_compute_mode_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.AzureComputeMode)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "azure_compute_mode_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("azure_compute_mode.#"); ok {
		erec.ComparisonInfo = new(comparison.AzureComputeMode)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "azure_compute_mode", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("azure_sku_comparision.#"); ok {
		erec.ComparisonInfo = new(comparison.AzureSku)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "azure_sku_comparision", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("azure_sku.#"); ok {
		erec.ComparisonInfo = new(comparison.AzureSku)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "azure_sku", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("base_comparison_basic.#"); ok {
		erec.ComparisonInfo = new(comparison.BaseComparison)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "base_comparison_basic", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.BaseComparison)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("bitness_comparision.#"); ok {
		erec.ComparisonInfo = new(comparison.Bitness)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "bitness_comparision", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("bitness.#"); ok {
		erec.ComparisonInfo = new(comparison.Bitness)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "bitness", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("cloud_type_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.CloudType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "cloud_type_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("cloud_type.#"); ok {
		erec.ComparisonInfo = new(comparison.CloudType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "cloud_type", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("custom_application_type_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.CustomApplicationType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_application_type_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("custom_application_type.#"); ok {
		erec.ComparisonInfo = new(comparison.CustomApplicationType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_application_type", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("hypervisor.#"); ok {
		erec.ComparisonInfo = new(comparison.HypervisorType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "hypervisor", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("hypervisor_type_comparision.#"); ok {
		erec.ComparisonInfo = new(comparison.HypervisorType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "hypervisor_type_comparision", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("database_topology_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.DatabaseTopology)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "database_topology_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("database_topology.#"); ok {
		erec.ComparisonInfo = new(comparison.DatabaseTopology)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "database_topology", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("dcrum_decoder.#"); ok {
		erec.ComparisonInfo = new(comparison.DCRumDecoder)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "dcrum_decoder", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("dcrum_decoder_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.DCRumDecoder)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "dcrum_decoder_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("entity.#"); ok {
		erec.ComparisonInfo = new(comparison.EntityID)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "entity", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("entity_id_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.EntityID)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "entity_id_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("entity_id_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.EntityID)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "entity_id_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("host_tech.#"); ok {
		erec.ComparisonInfo = new(comparison.SimpleHostTech)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "host_tech", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("simple_host_tech_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.SimpleHostTech)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "simple_host_tech_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("indexed_name.#"); ok {
		erec.ComparisonInfo = new(comparison.IndexedName)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "indexed_name", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("indexed_name_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.IndexedName)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "indexed_name_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("indexed_string.#"); ok {
		erec.ComparisonInfo = new(comparison.IndexedString)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "indexed_string", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("indexed_string_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.IndexedString)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "indexed_string_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("indexed_tag.#"); ok {
		erec.ComparisonInfo = new(comparison.IndexedTag)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "indexed_tag", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("indexed_tag_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.IndexedTag)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "indexed_tag_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("integer.#"); ok {
		erec.ComparisonInfo = new(comparison.Integer)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "integer", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("integer_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.Integer)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "integer_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("ipaddress.#"); ok {
		erec.ComparisonInfo = new(comparison.IPAddress)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "ipaddress", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("ipaddress_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.IPAddress)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "ipaddress_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("mobile_platform.#"); ok {
		erec.ComparisonInfo = new(comparison.MobilePlatform)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "mobile_platform", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("mobile_platform_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.MobilePlatform)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "mobile_platform_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("os_arch.#"); ok {
		erec.ComparisonInfo = new(comparison.OSArchitecture)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "os_arch", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("osarchitecture_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.OSArchitecture)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "osarchitecture_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("os_type.#"); ok {
		erec.ComparisonInfo = new(comparison.OSType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "os_type", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("ostype_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.OSType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "ostype_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("paas_type.#"); ok {
		erec.ComparisonInfo = new(comparison.PaasType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "paas_type", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("paas_type_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.PaasType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "paas_type_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("service_topology.#"); ok {
		erec.ComparisonInfo = new(comparison.ServiceTopology)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "service_topology", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("service_topology_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.ServiceTopology)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "service_topology_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("service_type.#"); ok {
		erec.ComparisonInfo = new(comparison.ServiceType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "service_type", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("service_type_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.ServiceType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "service_type_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("simple_tech_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.SimpleTech)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "simple_tech_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("tech.#"); ok {
		erec.ComparisonInfo = new(comparison.SimpleTech)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "tech", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("string.#"); ok {
		erec.ComparisonInfo = new(comparison.String)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "string", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("string_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.String)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "string_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("synthetic_engine.#"); ok {
		erec.ComparisonInfo = new(comparison.SyntheticEngineType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "synthetic_engine", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("synthetic_engine_type_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.SyntheticEngineType)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "synthetic_engine_type_comparison", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("tag.#"); ok {
		erec.ComparisonInfo = new(comparison.Tag)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "tag", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("tag_comparison.#"); ok {
		erec.ComparisonInfo = new(comparison.Tag)
		if err := erec.ComparisonInfo.UnmarshalHCL(hcl.NewDecoder(decoder, "tag_comparison", 0)); err != nil {
			return err
		}
	}

	if _, ok := decoder.GetOk("string_key.#"); ok {
		erec.Key = new(condition.String)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "string_key", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("string_condition_key.#"); ok {
		erec.Key = new(condition.String)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "string_condition_key", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("process_metadata.#"); ok {
		erec.Key = new(condition.ProcessMetadata)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "process_metadata", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("process_metadata_condition_key.#"); ok {
		erec.Key = new(condition.ProcessMetadata)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "process_metadata_condition_key", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("key.#"); ok {
		erec.Key = new(condition.BaseConditionKey)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "key", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("base_condition_key.#"); ok {
		erec.Key = new(condition.BaseConditionKey)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "base_condition_key", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("custom_host_metadata.#"); ok {
		erec.Key = new(condition.CustomHostMetadata)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_host_metadata", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("custom_host_metadata_condition_key.#"); ok {
		erec.Key = new(condition.CustomHostMetadata)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_host_metadata_condition_key", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("custom_process_metadata.#"); ok {
		erec.Key = new(condition.CustomProcessMetadata)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_process_metadata", 0)); err != nil {
			return err
		}
	} else if _, ok := decoder.GetOk("custom_process_metadata_condition_key.#"); ok {
		erec.Key = new(condition.CustomProcessMetadata)
		if err := erec.Key.UnmarshalHCL(hcl.NewDecoder(decoder, "custom_process_metadata_condition_key", 0)); err != nil {
			return err
		}
	}
	return nil
}

func (erec *Condition) MarshalJSON() ([]byte, error) {
	m := map[string]json.RawMessage{}
	if len(erec.Unknowns) > 0 {
		for k, v := range erec.Unknowns {
			m[k] = v
		}
	}
	if erec.ComparisonInfo != nil {
		rawMessage, err := json.Marshal(erec.ComparisonInfo)
		if err != nil {
			return nil, err
		}
		m["comparisonInfo"] = rawMessage
	}
	if erec.Key != nil {
		rawMessage, err := json.Marshal(erec.Key)
		if err != nil {
			return nil, err
		}
		m["key"] = rawMessage
	}
	return json.Marshal(m)
}

func (erec *Condition) UnmarshalJSON(data []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if v, found := m["comparisonInfo"]; found {
		x := struct {
			Type comparison.ComparisonBasicType `json:"type"`
		}{}
		if err := json.Unmarshal(v, &x); err != nil {
			return err
		}
		switch x.Type {
		case comparison.ComparisonBasicTypes.ApplicationType:
			cmp := new(comparison.ApplicationType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.AzureComputeMode:
			cmp := new(comparison.AzureComputeMode)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.AzureSku:
			cmp := new(comparison.AzureSku)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.Bitness:
			cmp := new(comparison.Bitness)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.CloudType:
			cmp := new(comparison.CloudType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.CustomApplicationType:
			cmp := new(comparison.CustomApplicationType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.DatabaseTopology:
			cmp := new(comparison.DatabaseTopology)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.DCRumDecoderType:
			cmp := new(comparison.DCRumDecoder)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.EntityID:
			cmp := new(comparison.EntityID)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.HypervisorType:
			cmp := new(comparison.HypervisorType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.IndexedName:
			cmp := new(comparison.IndexedName)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.IndexedString:
			cmp := new(comparison.IndexedString)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.IndexedTag:
			cmp := new(comparison.IndexedTag)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.Integer:
			cmp := new(comparison.Integer)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.IPAddress:
			cmp := new(comparison.IPAddress)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.MobilePlatform:
			cmp := new(comparison.MobilePlatform)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.OSArchitecture:
			cmp := new(comparison.OSArchitecture)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.OSType:
			cmp := new(comparison.OSType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.PaasType:
			cmp := new(comparison.PaasType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.ServiceTopology:
			cmp := new(comparison.ServiceTopology)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.ServiceType:
			cmp := new(comparison.ServiceType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.SimpleHostTech:
			cmp := new(comparison.SimpleHostTech)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.SimpleTech:
			cmp := new(comparison.SimpleTech)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.String:
			cmp := new(comparison.String)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.SyntheticEngineType:
			cmp := new(comparison.SyntheticEngineType)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		case comparison.ComparisonBasicTypes.Tag:
			cmp := new(comparison.Tag)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		default:
			cmp := new(comparison.BaseComparison)
			if err := json.Unmarshal(v, cmp); err != nil {
				return err
			}
			erec.ComparisonInfo = cmp
		}
	}
	if v, found := m["key"]; found {
		x := struct {
			Type condition.ConditionKeyType `json:"type"`
		}{}
		if err := json.Unmarshal(v, &x); err != nil {
			return err
		}

		switch x.Type {
		case condition.ConditionKeyTypes.HostCustomMetadataKey:
			key := new(condition.CustomHostMetadata)
			if err := json.Unmarshal(v, key); err != nil {
				return err
			}
			erec.Key = key
		case condition.ConditionKeyTypes.ProcessCustomMetadataKey:
			key := new(condition.CustomProcessMetadata)
			if err := json.Unmarshal(v, key); err != nil {
				return err
			}
			erec.Key = key
		case condition.ConditionKeyTypes.ProcessPredefinedMetadataKey:
			key := new(condition.ProcessMetadata)
			if err := json.Unmarshal(v, key); err != nil {
				return err
			}
			erec.Key = key
		case condition.ConditionKeyTypes.String:
			key := new(condition.String)
			if err := json.Unmarshal(v, key); err != nil {
				return err
			}
			erec.Key = key
		default:
			key := new(condition.BaseConditionKey)
			if err := json.Unmarshal(v, key); err != nil {
				return err
			}
			erec.Key = key
		}
	}

	delete(m, "comparisonInfo")
	delete(m, "key")
	if len(m) > 0 {
		erec.Unknowns = m
	}
	return nil
}
