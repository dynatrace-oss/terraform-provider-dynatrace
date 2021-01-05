---
page_title: "dynatrace_maintenance_window Resource - terraform-provider-dynatrace-1"
subcategory: ""
description: |-
  
---

# Resource `dynatrace_maintenance_window`





## Schema

### Optional

- **description** (String)
- **id** (String) The ID of this resource.
- **metadata** (Block List, Max: 1) (see [below for nested schema](#nestedblock--metadata))
- **name** (String)
- **schedule** (Block List, Max: 1) (see [below for nested schema](#nestedblock--schedule))
- **scope** (Block List, Max: 1) (see [below for nested schema](#nestedblock--scope))
- **suppression** (String)
- **type** (String)

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- **cluster_version** (String)
- **configuration_versions** (List of Number)
- **current_configuration_versions** (List of String)


<a id="nestedblock--schedule"></a>
### Nested Schema for `schedule`

Optional:

- **end** (String)
- **recurrence** (Block List, Max: 1) (see [below for nested schema](#nestedblock--schedule--recurrence))
- **recurrence_type** (String)
- **start** (String)
- **zone_id** (String)

<a id="nestedblock--schedule--recurrence"></a>
### Nested Schema for `schedule.recurrence`

Optional:

- **day_of_month** (Number)
- **day_of_week** (String)
- **duration_minutes** (Number)
- **start_time** (String)



<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Optional:

- **entities** (List of String)
- **matches** (Block List) (see [below for nested schema](#nestedblock--scope--matches))

<a id="nestedblock--scope--matches"></a>
### Nested Schema for `scope.matches`

Optional:

- **mz_id** (String)
- **tag_combination** (String)
- **tags** (Block List) (see [below for nested schema](#nestedblock--scope--matches--tags))
- **type** (String)

<a id="nestedblock--scope--matches--tags"></a>
### Nested Schema for `scope.matches.tags`

Optional:

- **context** (String)
- **key** (String)
- **value** (String)


