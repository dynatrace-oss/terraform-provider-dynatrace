---
page_title: "dynatrace_management_zone Resource - terraform-provider-dynatrace"
subcategory: ""
description: |-
  
---

# Resource `dynatrace_management_zone`





## Schema

### Optional

- **id** (String) The ID of this resource.
- **metadata** (Block List, Max: 1) (see [below for nested schema](#nestedblock--metadata))
- **name** (String)
- **rules** (Block List) (see [below for nested schema](#nestedblock--rules))

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- **cluster_version** (String)
- **configuration_versions** (List of Number)
- **current_configuration_versions** (List of String)


<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Optional:

- **conditions** (Block List) (see [below for nested schema](#nestedblock--rules--conditions))
- **enabled** (Boolean)
- **propagation_types** (List of String)
- **type** (String)

<a id="nestedblock--rules--conditions"></a>
### Nested Schema for `rules.conditions`

Optional:

- **application_type_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--application_type_comparison))
- **azure_compute_mode_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--azure_compute_mode_comparison))
- **azure_sku_comparision** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--azure_sku_comparision))
- **base_comparison_basic** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--base_comparison_basic))
- **base_condition_key** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--base_condition_key))
- **bitness_comparision** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--bitness_comparision))
- **cloud_type_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--cloud_type_comparison))
- **custom_application_type_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--custom_application_type_comparison))
- **custom_host_metadata_condition_key** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--custom_host_metadata_condition_key))
- **custom_process_metadata_condition_key** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--custom_process_metadata_condition_key))
- **database_topology_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--database_topology_comparison))
- **dcrum_decoder_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--dcrum_decoder_comparison))
- **entity_id_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--entity_id_comparison))
- **hypervisor_type_comparision** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--hypervisor_type_comparision))
- **indexed_name_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--indexed_name_comparison))
- **indexed_string_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--indexed_string_comparison))
- **indexed_tag_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--indexed_tag_comparison))
- **integer_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--integer_comparison))
- **ipaddress_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--ipaddress_comparison))
- **mobile_platform_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--mobile_platform_comparison))
- **osarchitecture_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--osarchitecture_comparison))
- **ostype_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--ostype_comparison))
- **paas_type_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--paas_type_comparison))
- **process_metadata_condition_key** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--process_metadata_condition_key))
- **service_topology_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--service_topology_comparison))
- **service_type_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--service_type_comparison))
- **simple_host_tech_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--simple_host_tech_comparison))
- **simple_tech_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--simple_tech_comparison))
- **string_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--string_comparison))
- **string_condition_key** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--string_condition_key))
- **synthetic_engine_type_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--synthetic_engine_type_comparison))
- **tag_comparison** (Block List) (see [below for nested schema](#nestedblock--rules--conditions--tag_comparison))

<a id="nestedblock--rules--conditions--application_type_comparison"></a>
### Nested Schema for `rules.conditions.application_type_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--azure_compute_mode_comparison"></a>
### Nested Schema for `rules.conditions.azure_compute_mode_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--azure_sku_comparision"></a>
### Nested Schema for `rules.conditions.azure_sku_comparision`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--base_comparison_basic"></a>
### Nested Schema for `rules.conditions.base_comparison_basic`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--base_condition_key"></a>
### Nested Schema for `rules.conditions.base_condition_key`

Optional:

- **attribute** (String)
- **dynamic_key** (String)
- **type** (String)


<a id="nestedblock--rules--conditions--bitness_comparision"></a>
### Nested Schema for `rules.conditions.bitness_comparision`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--cloud_type_comparison"></a>
### Nested Schema for `rules.conditions.cloud_type_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--custom_application_type_comparison"></a>
### Nested Schema for `rules.conditions.custom_application_type_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--custom_host_metadata_condition_key"></a>
### Nested Schema for `rules.conditions.custom_host_metadata_condition_key`

Optional:

- **attribute** (String)
- **dynamic_key** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--conditions--custom_host_metadata_condition_key--dynamic_key))
- **type** (String)

<a id="nestedblock--rules--conditions--custom_host_metadata_condition_key--dynamic_key"></a>
### Nested Schema for `rules.conditions.custom_host_metadata_condition_key.type`

Optional:

- **key** (String)
- **source** (String)



<a id="nestedblock--rules--conditions--custom_process_metadata_condition_key"></a>
### Nested Schema for `rules.conditions.custom_process_metadata_condition_key`

Optional:

- **attribute** (String)
- **dynamic_key** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--conditions--custom_process_metadata_condition_key--dynamic_key))
- **type** (String)

<a id="nestedblock--rules--conditions--custom_process_metadata_condition_key--dynamic_key"></a>
### Nested Schema for `rules.conditions.custom_process_metadata_condition_key.type`

Optional:

- **key** (String)
- **source** (String)



<a id="nestedblock--rules--conditions--database_topology_comparison"></a>
### Nested Schema for `rules.conditions.database_topology_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--dcrum_decoder_comparison"></a>
### Nested Schema for `rules.conditions.dcrum_decoder_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--entity_id_comparison"></a>
### Nested Schema for `rules.conditions.entity_id_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--hypervisor_type_comparision"></a>
### Nested Schema for `rules.conditions.hypervisor_type_comparision`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--indexed_name_comparison"></a>
### Nested Schema for `rules.conditions.indexed_name_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--indexed_string_comparison"></a>
### Nested Schema for `rules.conditions.indexed_string_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--indexed_tag_comparison"></a>
### Nested Schema for `rules.conditions.indexed_tag_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--conditions--indexed_tag_comparison--value))

<a id="nestedblock--rules--conditions--indexed_tag_comparison--value"></a>
### Nested Schema for `rules.conditions.indexed_tag_comparison.value`

Optional:

- **context** (String)
- **key** (String)
- **value** (String)



<a id="nestedblock--rules--conditions--integer_comparison"></a>
### Nested Schema for `rules.conditions.integer_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (Number)


<a id="nestedblock--rules--conditions--ipaddress_comparison"></a>
### Nested Schema for `rules.conditions.ipaddress_comparison`

Optional:

- **case_sensitive** (Boolean)
- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--mobile_platform_comparison"></a>
### Nested Schema for `rules.conditions.mobile_platform_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--osarchitecture_comparison"></a>
### Nested Schema for `rules.conditions.osarchitecture_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--ostype_comparison"></a>
### Nested Schema for `rules.conditions.ostype_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--paas_type_comparison"></a>
### Nested Schema for `rules.conditions.paas_type_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--process_metadata_condition_key"></a>
### Nested Schema for `rules.conditions.process_metadata_condition_key`

Optional:

- **attribute** (String)
- **dynamic_key** (String)
- **type** (String)


<a id="nestedblock--rules--conditions--service_topology_comparison"></a>
### Nested Schema for `rules.conditions.service_topology_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--service_type_comparison"></a>
### Nested Schema for `rules.conditions.service_type_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--simple_host_tech_comparison"></a>
### Nested Schema for `rules.conditions.simple_host_tech_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--conditions--simple_host_tech_comparison--value))

<a id="nestedblock--rules--conditions--simple_host_tech_comparison--value"></a>
### Nested Schema for `rules.conditions.simple_host_tech_comparison.value`

Optional:

- **type** (String)
- **verbatim_type** (String)



<a id="nestedblock--rules--conditions--simple_tech_comparison"></a>
### Nested Schema for `rules.conditions.simple_tech_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--conditions--simple_tech_comparison--value))

<a id="nestedblock--rules--conditions--simple_tech_comparison--value"></a>
### Nested Schema for `rules.conditions.simple_tech_comparison.value`

Optional:

- **type** (String)
- **verbatim_type** (String)



<a id="nestedblock--rules--conditions--string_comparison"></a>
### Nested Schema for `rules.conditions.string_comparison`

Optional:

- **case_sensitive** (Boolean)
- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--string_condition_key"></a>
### Nested Schema for `rules.conditions.string_condition_key`

Optional:

- **attribute** (String)
- **dynamic_key** (String)
- **type** (String)


<a id="nestedblock--rules--conditions--synthetic_engine_type_comparison"></a>
### Nested Schema for `rules.conditions.synthetic_engine_type_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (String)


<a id="nestedblock--rules--conditions--tag_comparison"></a>
### Nested Schema for `rules.conditions.tag_comparison`

Optional:

- **negate** (Boolean)
- **operator** (String)
- **type** (String)
- **value** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--conditions--tag_comparison--value))

<a id="nestedblock--rules--conditions--tag_comparison--value"></a>
### Nested Schema for `rules.conditions.tag_comparison.value`

Optional:

- **context** (String)
- **key** (String)
- **value** (String)


