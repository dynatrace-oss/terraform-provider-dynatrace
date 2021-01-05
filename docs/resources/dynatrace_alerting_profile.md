---
page_title: "dynatrace_alerting_profile Resource - terraform-provider-dynatrace-1"
subcategory: ""
description: |-
  
---

# Resource `dynatrace_alerting_profile`





## Schema

### Optional

- **display_name** (String)
- **event_type_filters** (Block List) (see [below for nested schema](#nestedblock--event_type_filters))
- **id** (String) The ID of this resource.
- **metadata** (Block List, Max: 1) (see [below for nested schema](#nestedblock--metadata))
- **mz_id** (String)
- **rules** (Block List) (see [below for nested schema](#nestedblock--rules))

<a id="nestedblock--event_type_filters"></a>
### Nested Schema for `event_type_filters`

Optional:

- **custom_event_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--event_type_filters--custom_event_filter))
- **predefined_event_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--event_type_filters--predefined_event_filter))

<a id="nestedblock--event_type_filters--custom_event_filter"></a>
### Nested Schema for `event_type_filters.custom_event_filter`

Optional:

- **custom_description_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--event_type_filters--custom_event_filter--custom_description_filter))
- **custom_title_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--event_type_filters--custom_event_filter--custom_title_filter))

<a id="nestedblock--event_type_filters--custom_event_filter--custom_description_filter"></a>
### Nested Schema for `event_type_filters.custom_event_filter.custom_description_filter`

Optional:

- **case_insensitive** (Boolean)
- **enabled** (Boolean)
- **negate** (Boolean)
- **operator** (String)
- **value** (String)


<a id="nestedblock--event_type_filters--custom_event_filter--custom_title_filter"></a>
### Nested Schema for `event_type_filters.custom_event_filter.custom_title_filter`

Optional:

- **case_insensitive** (Boolean)
- **enabled** (Boolean)
- **negate** (Boolean)
- **operator** (String)
- **value** (String)



<a id="nestedblock--event_type_filters--predefined_event_filter"></a>
### Nested Schema for `event_type_filters.predefined_event_filter`

Optional:

- **event_type** (String)
- **negate** (Boolean)



<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- **cluster_version** (String)
- **configuration_versions** (List of Number)
- **current_configuration_versions** (List of String)


<a id="nestedblock--rules"></a>
### Nested Schema for `rules`

Optional:

- **delay_in_minutes** (Number)
- **severity_level** (String)
- **tag_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--rules--tag_filter))

<a id="nestedblock--rules--tag_filter"></a>
### Nested Schema for `rules.tag_filter`

Optional:

- **include_mode** (String)
- **tag_filters** (Block List) (see [below for nested schema](#nestedblock--rules--tag_filter--tag_filters))

<a id="nestedblock--rules--tag_filter--tag_filters"></a>
### Nested Schema for `rules.tag_filter.tag_filters`

Optional:

- **context** (String)
- **key** (String)
- **value** (String)


