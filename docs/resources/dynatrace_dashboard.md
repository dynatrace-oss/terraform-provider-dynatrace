---
page_title: "dynatrace_dashboard Resource - terraform-provider-dynatrace-1"
subcategory: ""
description: |-
  
---

# Resource `dynatrace_dashboard`





## Schema

### Optional

- **abstract_tile** (Block List) (see [below for nested schema](#nestedblock--abstract_tile))
- **assigned_entities_tile** (Block List) (see [below for nested schema](#nestedblock--assigned_entities_tile))
- **assigned_entities_with_metric_tile** (Block List) (see [below for nested schema](#nestedblock--assigned_entities_with_metric_tile))
- **custom_charting_tile** (Block List) (see [below for nested schema](#nestedblock--custom_charting_tile))
- **dashboard_metadata** (Block List, Max: 1) (see [below for nested schema](#nestedblock--dashboard_metadata))
- **filterable_entity_tile** (Block List) (see [below for nested schema](#nestedblock--filterable_entity_tile))
- **id** (String) The ID of this resource.
- **markdown_tile** (Block List) (see [below for nested schema](#nestedblock--markdown_tile))
- **metadata** (Block List, Max: 1) (see [below for nested schema](#nestedblock--metadata))
- **synthetic_single_web_check_tile** (Block List) (see [below for nested schema](#nestedblock--synthetic_single_web_check_tile))
- **user_session_query_tile** (Block List) (see [below for nested schema](#nestedblock--user_session_query_tile))

<a id="nestedblock--abstract_tile"></a>
### Nested Schema for `abstract_tile`

Optional:

- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--abstract_tile--bounds))
- **configured** (Boolean)
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--abstract_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--abstract_tile--bounds"></a>
### Nested Schema for `abstract_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--abstract_tile--tile_filter"></a>
### Nested Schema for `abstract_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--abstract_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--abstract_tile--tile_filter--management_zone"></a>
### Nested Schema for `abstract_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--assigned_entities_tile"></a>
### Nested Schema for `assigned_entities_tile`

Optional:

- **assigned_entities** (List of String)
- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--assigned_entities_tile--bounds))
- **configured** (Boolean)
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--assigned_entities_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--assigned_entities_tile--bounds"></a>
### Nested Schema for `assigned_entities_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--assigned_entities_tile--tile_filter"></a>
### Nested Schema for `assigned_entities_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--assigned_entities_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--assigned_entities_tile--tile_filter--management_zone"></a>
### Nested Schema for `assigned_entities_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--assigned_entities_with_metric_tile"></a>
### Nested Schema for `assigned_entities_with_metric_tile`

Optional:

- **assigned_entities** (List of String)
- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--assigned_entities_with_metric_tile--bounds))
- **configured** (Boolean)
- **metric** (String)
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--assigned_entities_with_metric_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--assigned_entities_with_metric_tile--bounds"></a>
### Nested Schema for `assigned_entities_with_metric_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--assigned_entities_with_metric_tile--tile_filter"></a>
### Nested Schema for `assigned_entities_with_metric_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--assigned_entities_with_metric_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--assigned_entities_with_metric_tile--tile_filter--management_zone"></a>
### Nested Schema for `assigned_entities_with_metric_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--custom_charting_tile"></a>
### Nested Schema for `custom_charting_tile`

Optional:

- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--custom_charting_tile--bounds))
- **configured** (Boolean)
- **filter_config** (Block List, Max: 1) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config))
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--custom_charting_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--custom_charting_tile--bounds"></a>
### Nested Schema for `custom_charting_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--custom_charting_tile--filter_config"></a>
### Nested Schema for `custom_charting_tile.filter_config`

Optional:

- **chart_config** (Block List, Max: 1) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--chart_config))
- **custom_name** (String)
- **default_name** (String)
- **filters_per_entity_type** (Block List) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--filters_per_entity_type))
- **type** (String)

<a id="nestedblock--custom_charting_tile--filter_config--chart_config"></a>
### Nested Schema for `custom_charting_tile.filter_config.chart_config`

Optional:

- **axis_limits** (Block List) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--chart_config--axis_limits))
- **left_axis_custom_unit** (String)
- **legend_shown** (Boolean)
- **result_metadata** (Block List) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--chart_config--result_metadata))
- **right_axis_custom_unit** (String)
- **series** (Block List) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--chart_config--series))
- **type** (String)

<a id="nestedblock--custom_charting_tile--filter_config--chart_config--axis_limits"></a>
### Nested Schema for `custom_charting_tile.filter_config.chart_config.type`

Required:

- **key** (String)
- **values** (Number)


<a id="nestedblock--custom_charting_tile--filter_config--chart_config--result_metadata"></a>
### Nested Schema for `custom_charting_tile.filter_config.chart_config.type`

Required:

- **key** (String)

Optional:

- **custom_color** (String)
- **last_modified** (Number)


<a id="nestedblock--custom_charting_tile--filter_config--chart_config--series"></a>
### Nested Schema for `custom_charting_tile.filter_config.chart_config.type`

Optional:

- **aggregation** (String)
- **aggregation_rate** (String)
- **dimensions** (Block List) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--chart_config--type--dimensions))
- **entity_type** (String)
- **metric** (String)
- **percentile** (Number)
- **sort_ascending** (Boolean)
- **sort_column** (Boolean)
- **type** (String)

<a id="nestedblock--custom_charting_tile--filter_config--chart_config--type--dimensions"></a>
### Nested Schema for `custom_charting_tile.filter_config.chart_config.type.dimensions`

Optional:

- **entity_dimension** (Boolean)
- **id** (String) The ID of this resource.
- **name** (String)
- **values** (List of String)




<a id="nestedblock--custom_charting_tile--filter_config--filters_per_entity_type"></a>
### Nested Schema for `custom_charting_tile.filter_config.filters_per_entity_type`

Required:

- **key** (String)
- **values** (Block List, Min: 1) (see [below for nested schema](#nestedblock--custom_charting_tile--filter_config--filters_per_entity_type--values))

<a id="nestedblock--custom_charting_tile--filter_config--filters_per_entity_type--values"></a>
### Nested Schema for `custom_charting_tile.filter_config.filters_per_entity_type.values`

Required:

- **key** (String)
- **values** (List of String)




<a id="nestedblock--custom_charting_tile--tile_filter"></a>
### Nested Schema for `custom_charting_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--custom_charting_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--custom_charting_tile--tile_filter--management_zone"></a>
### Nested Schema for `custom_charting_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--dashboard_metadata"></a>
### Nested Schema for `dashboard_metadata`

Optional:

- **dashboard_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--dashboard_metadata--dashboard_filter))
- **name** (String)
- **owner** (String)
- **preset** (Boolean)
- **shared** (Boolean)
- **sharing_details** (Block List, Max: 1) (see [below for nested schema](#nestedblock--dashboard_metadata--sharing_details))
- **tags** (List of String)
- **valid_filter_keys** (List of String)

<a id="nestedblock--dashboard_metadata--dashboard_filter"></a>
### Nested Schema for `dashboard_metadata.dashboard_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--dashboard_metadata--dashboard_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--dashboard_metadata--dashboard_filter--management_zone"></a>
### Nested Schema for `dashboard_metadata.dashboard_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)



<a id="nestedblock--dashboard_metadata--sharing_details"></a>
### Nested Schema for `dashboard_metadata.sharing_details`

Optional:

- **link_shared** (Boolean)
- **published** (Boolean)



<a id="nestedblock--filterable_entity_tile"></a>
### Nested Schema for `filterable_entity_tile`

Optional:

- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--filterable_entity_tile--bounds))
- **chart_visible** (Boolean)
- **configured** (Boolean)
- **filter_config** (Block List, Max: 1) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config))
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--filterable_entity_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--filterable_entity_tile--bounds"></a>
### Nested Schema for `filterable_entity_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--filterable_entity_tile--filter_config"></a>
### Nested Schema for `filterable_entity_tile.filter_config`

Optional:

- **chart_config** (Block List, Max: 1) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--chart_config))
- **custom_name** (String)
- **default_name** (String)
- **filters_per_entity_type** (Block List) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--filters_per_entity_type))
- **type** (String)

<a id="nestedblock--filterable_entity_tile--filter_config--chart_config"></a>
### Nested Schema for `filterable_entity_tile.filter_config.chart_config`

Optional:

- **axis_limits** (Block List) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--chart_config--axis_limits))
- **left_axis_custom_unit** (String)
- **legend_shown** (Boolean)
- **result_metadata** (Block List) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--chart_config--result_metadata))
- **right_axis_custom_unit** (String)
- **series** (Block List) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--chart_config--series))
- **type** (String)

<a id="nestedblock--filterable_entity_tile--filter_config--chart_config--axis_limits"></a>
### Nested Schema for `filterable_entity_tile.filter_config.chart_config.type`

Required:

- **key** (String)
- **values** (Number)


<a id="nestedblock--filterable_entity_tile--filter_config--chart_config--result_metadata"></a>
### Nested Schema for `filterable_entity_tile.filter_config.chart_config.type`

Required:

- **key** (String)

Optional:

- **custom_color** (String)
- **last_modified** (Number)


<a id="nestedblock--filterable_entity_tile--filter_config--chart_config--series"></a>
### Nested Schema for `filterable_entity_tile.filter_config.chart_config.type`

Optional:

- **aggregation** (String)
- **aggregation_rate** (String)
- **dimensions** (Block List) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--chart_config--type--dimensions))
- **entity_type** (String)
- **metric** (String)
- **percentile** (Number)
- **sort_ascending** (Boolean)
- **sort_column** (Boolean)
- **type** (String)

<a id="nestedblock--filterable_entity_tile--filter_config--chart_config--type--dimensions"></a>
### Nested Schema for `filterable_entity_tile.filter_config.chart_config.type.dimensions`

Optional:

- **entity_dimension** (Boolean)
- **id** (String) The ID of this resource.
- **name** (String)
- **values** (List of String)




<a id="nestedblock--filterable_entity_tile--filter_config--filters_per_entity_type"></a>
### Nested Schema for `filterable_entity_tile.filter_config.filters_per_entity_type`

Required:

- **key** (String)
- **values** (Block List, Min: 1) (see [below for nested schema](#nestedblock--filterable_entity_tile--filter_config--filters_per_entity_type--values))

<a id="nestedblock--filterable_entity_tile--filter_config--filters_per_entity_type--values"></a>
### Nested Schema for `filterable_entity_tile.filter_config.filters_per_entity_type.values`

Required:

- **key** (String)
- **values** (List of String)




<a id="nestedblock--filterable_entity_tile--tile_filter"></a>
### Nested Schema for `filterable_entity_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--filterable_entity_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--filterable_entity_tile--tile_filter--management_zone"></a>
### Nested Schema for `filterable_entity_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--markdown_tile"></a>
### Nested Schema for `markdown_tile`

Optional:

- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--markdown_tile--bounds))
- **configured** (Boolean)
- **markdown** (String)
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--markdown_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--markdown_tile--bounds"></a>
### Nested Schema for `markdown_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--markdown_tile--tile_filter"></a>
### Nested Schema for `markdown_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--markdown_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--markdown_tile--tile_filter--management_zone"></a>
### Nested Schema for `markdown_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- **cluster_version** (String)
- **configuration_versions** (List of Number)
- **current_configuration_versions** (List of String)


<a id="nestedblock--synthetic_single_web_check_tile"></a>
### Nested Schema for `synthetic_single_web_check_tile`

Optional:

- **assigned_entities** (List of String)
- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--synthetic_single_web_check_tile--bounds))
- **configured** (Boolean)
- **exclude_maintenance_windows** (Boolean)
- **name** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--synthetic_single_web_check_tile--tile_filter))
- **tile_type** (String)

<a id="nestedblock--synthetic_single_web_check_tile--bounds"></a>
### Nested Schema for `synthetic_single_web_check_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--synthetic_single_web_check_tile--tile_filter"></a>
### Nested Schema for `synthetic_single_web_check_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--synthetic_single_web_check_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--synthetic_single_web_check_tile--tile_filter--management_zone"></a>
### Nested Schema for `synthetic_single_web_check_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)




<a id="nestedblock--user_session_query_tile"></a>
### Nested Schema for `user_session_query_tile`

Optional:

- **bounds** (Block List, Max: 1) (see [below for nested schema](#nestedblock--user_session_query_tile--bounds))
- **configured** (Boolean)
- **custom_name** (String)
- **limit** (Number)
- **name** (String)
- **query** (String)
- **tile_filter** (Block List, Max: 1) (see [below for nested schema](#nestedblock--user_session_query_tile--tile_filter))
- **tile_type** (String)
- **time_frame_shift** (String)
- **type** (String)
- **visualization_config** (Block List, Max: 1) (see [below for nested schema](#nestedblock--user_session_query_tile--visualization_config))

<a id="nestedblock--user_session_query_tile--bounds"></a>
### Nested Schema for `user_session_query_tile.bounds`

Optional:

- **height** (Number)
- **left** (Number)
- **top** (Number)
- **width** (Number)


<a id="nestedblock--user_session_query_tile--tile_filter"></a>
### Nested Schema for `user_session_query_tile.tile_filter`

Optional:

- **management_zone** (Block List, Max: 1) (see [below for nested schema](#nestedblock--user_session_query_tile--tile_filter--management_zone))
- **timeframe** (String)

<a id="nestedblock--user_session_query_tile--tile_filter--management_zone"></a>
### Nested Schema for `user_session_query_tile.tile_filter.management_zone`

Optional:

- **description** (String)
- **id** (String) The ID of this resource.
- **name** (String)



<a id="nestedblock--user_session_query_tile--visualization_config"></a>
### Nested Schema for `user_session_query_tile.visualization_config`

Optional:

- **has_axis_bucketing** (Boolean)


