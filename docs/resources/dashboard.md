---
layout: ""
page_title: dynatrace_dashboard Resource - terraform-provider-dynatrace"
description: |-
  The resource `dynatrace_dashboard` covers configuration for dashboards
---

# dynatrace_dashboard (Resource)

## Dynatrace Documentation

- Dashboards and reports - https://www.dynatrace.com/support/help/how-to-use-dynatrace/dashboards-and-charts

- Dashboards API - https://www.dynatrace.com/support/help/dynatrace-api/configuration-api/dashboards-api

## Export Example Usage

- `terraform-provider-dynatrace -export dynatrace_dashboard` downloads all existing dashboard configuration

The full documentation of the export feature is available [here](https://registry.terraform.io/providers/dynatrace-oss/dynatrace/latest/docs/guides/export-v2).

## Resource Example Usage

```terraform
resource "dynatrace_dashboard" "#name#" {
  dashboard_metadata {
    name   = "#name#"
    owner  = "Dynatrace"
    tags   = ["Kubernetes"]
    dynamic_filters {
      filters = ["KUBERNETES_CLUSTER"]
    }
  }
  tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 0
      width  = 684
      height = 38
      left   = 0
    }
    markdown = "## Cluster resource overview"
  }
  tile {
    filter_config {
      default_name = "Full-Stack Kubernetes nodes"
      chart_config {
        legend = true
        type   = "TIMESERIES"
      }
      filters {
        filter {
          entity_type = "HOST"
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
        }
      }
      type        = "HOST"
      custom_name = "Full-Stack Kubernetes nodes"
    }
    chart_visible = true
    name          = ""
    tile_type     = "HOSTS"
    configured    = true
    bounds {
      width  = 342
      height = 304
      left   = 342
      top    = 38
    }
    filter {
      timeframe = "-5m"
    }
  }
  tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      width  = 190
      height = 152
      left   = 190
      top    = 418
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "CPU available"
      default_name = "Custom chart"
      chart_config {
        type = "SINGLE_VALUE"
        series {
          aggregation_rate = "TOTAL"
          metric           = "builtin:cloud.kubernetes.cluster.cpuAvailable"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
        }
        legend = true
      }
    }
  }
  tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      height = 304
      left   = 684
      top    = 38
      width  = 304
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      custom_name  = "Pods"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "PIE"
        series {
          sort_column      = true
          aggregation_rate = "TOTAL"
          dimension {
            name             = "Pod phase"
            entity_dimension = false
            id               = "1"
          }
          metric         = "builtin:cloud.kubernetes.workload.pods"
          aggregation    = "SUM_DIMENSIONS"
          type           = "LINE"
          entity_type    = "CLOUD_APPLICATION"
          sort_ascending = false
        }
        result_metadata {
          config {
            last_modified = 1597237249882
            custom_color  = "#008cdb"
            key           = "null¦Pod phase»Succeeded»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            last_modified = 1597234642722
            custom_color  = "#64bd64"
            key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            last_modified = 1597234457744
            custom_color  = "#f5d30f"
            key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            custom_color  = "#ff0000"
            key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
            last_modified = 1597234118116
          }
        }
      }
      type = "MIXED"
    }
  }
  tile {
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      left   = 608
      top    = 418
      width  = 190
      height = 152
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      chart_config {
        legend = true
        type   = "SINGLE_VALUE"
        series {
          aggregation_rate = "TOTAL"
          metric           = "builtin:cloud.kubernetes.cluster.memoryAvailable"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
        }
      }
      type         = "MIXED"
      custom_name  = "Memory available"
      default_name = "Custom chart"
    }
    name = "Custom chart"
  }
  tile {
    configured = true
    bounds {
      left   = 0
      top    = 380
      width  = 1634
      height = 38
    }
    markdown  = "## Node resource usage"
    name      = "Markdown"
    tile_type = "MARKDOWN"
  }
  tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      left   = 0
      top    = 38
      width  = 342
      height = 304
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Cluster nodes"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "PIE"
        series {
          metric           = "builtin:cloud.kubernetes.cluster.nodes"
          aggregation      = "AVG"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
          dimension {
            id               = "0"
            name             = "dt.entity.kubernetes_cluster"
            entity_dimension = true
          }
        }
      }
    }
  }
  tile {
    filter_config {
      type         = "MIXED"
      custom_name  = "Disk available"
      default_name = "Custom chart"
      chart_config {
        series {
          metric           = "builtin:host.disk.avail"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "HOST"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        legend = true
        type   = "SINGLE_VALUE"
      }
      filters {
        filter {
          entity_type = "HOST"
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
        }
      }
    }
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      height = 152
      left   = 1026
      top    = 418
      width  = 190
    }
    filter {
      timeframe = "-5m"
    }
  }
  tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      height = 304
      left   = 0
      top    = 570
      width  = 418
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "CPU usage % "
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "TIMESERIES"
        series {
          metric           = "builtin:host.cpu.usage"
          aggregation      = "AVG"
          type             = "LINE"
          entity_type      = "HOST"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
          dimension {
            id               = "0"
            name             = "dt.entity.host"
            entity_dimension = true
          }
        }
      }
      filters {
        filter {
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
          entity_type = "HOST"
        }
      }
    }
  }
  tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      width  = 418
      height = 304
      left   = 418
      top    = 570
    }
    filter_config {
      chart_config {
        legend = true
        type   = "TIMESERIES"
        series {
          aggregation_rate = "TOTAL"
          dimension {
            entity_dimension = true
            id               = "0"
            name             = "dt.entity.host"
          }
          metric         = "builtin:host.mem.usage"
          aggregation    = "AVG"
          type           = "LINE"
          entity_type    = "HOST"
          sort_ascending = false
          sort_column    = true
        }
      }
      filters {
        filter {
          entity_type = "HOST"
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
        }
      }
      type         = "MIXED"
      custom_name  = "Memory usage % "
      default_name = "Custom chart"
    }
  }
  tile {
    bounds {
      left   = 836
      top    = 570
      width  = 418
      height = 304
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Disk usage % "
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "TIMESERIES"
        series {
          aggregation_rate = "TOTAL"
          dimension {
            id               = "0"
            name             = "dt.entity.host"
            entity_dimension = true
          }
          dimension {
            entity_dimension = true
            id               = "1"
            name             = "dt.entity.disk"
          }
          metric         = "builtin:host.disk.usedPct"
          aggregation    = "AVG"
          type           = "LINE"
          entity_type    = "HOST"
          sort_ascending = false
          sort_column    = true
        }
      }
      filters {
        filter {
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
          entity_type = "HOST"
        }
      }
    }
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
  }
  tile {
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      left   = 0
      top    = 418
      width  = 190
      height = 152
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      custom_name  = "Total CPU requests"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "SINGLE_VALUE"
        series {
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
          metric           = "builtin:cloud.kubernetes.cluster.cpuRequested"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
        }
      }
      type = "MIXED"
    }
    name = "Custom chart"
  }
  tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      height = 152
      left   = 418
      top    = 418
      width  = 190
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      chart_config {
        legend = true
        type   = "SINGLE_VALUE"
        series {
          aggregation_rate = "TOTAL"
          metric           = "builtin:cloud.kubernetes.cluster.memoryRequested"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
        }
      }
      type         = "MIXED"
      custom_name  = "Total memory requests"
      default_name = "Custom chart"
    }
  }
  tile {
    filter {
      timeframe = "-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Total disk used"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "SINGLE_VALUE"
        series {
          metric           = "builtin:host.disk.used"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "HOST"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
      filters {
        filter {
          entity_type = "HOST"
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
        }
      }
    }
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      left   = 836
      top    = 418
      width  = 190
      height = 152
    }
  }
  tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      width  = 380
      height = 304
      left   = 1254
      top    = 570
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Traffic in/out"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "TIMESERIES"
        series {
          aggregation_rate = "TOTAL"
          dimension {
            id               = "0"
            name             = "dt.entity.host"
            entity_dimension = true
          }
          metric         = "builtin:host.net.nic.trafficIn"
          aggregation    = "SUM_DIMENSIONS"
          type           = "LINE"
          entity_type    = "HOST"
          sort_ascending = false
          sort_column    = false
        }
        series {
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "HOST"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
          dimension {
            id               = "0"
            name             = "dt.entity.host"
            entity_dimension = true
          }
          metric = "builtin:host.net.nic.trafficOut"
        }
      }
      filters {
        filter {
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
          entity_type = "HOST"
        }
      }
    }
  }
  tile {
    filter {
      timeframe = "-5m"
    }
    filter_config {
      chart_config {
        legend = true
        type   = "SINGLE_VALUE"
        series {
          metric           = "builtin:host.net.nic.trafficOut"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "HOST"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
      filters {
        filter {
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
          entity_type = "HOST"
        }
      }
      type         = "MIXED"
      custom_name  = "Traffic out"
      default_name = "Custom chart"
    }
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      left   = 1444
      top    = 418
      width  = 190
      height = 152
    }
  }
  tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      height = 152
      left   = 1254
      top    = 418
      width  = 190
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      custom_name  = "Traffic in"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "SINGLE_VALUE"
        series {
          metric           = "builtin:host.net.nic.trafficIn"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "HOST"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
      filters {
        filter {
          entity_type = "HOST"
          match {
            key    = "HOST_SOFTWARE_TECH"
            values = ["KUBERNETES"]
          }
        }
      }
      type = "MIXED"
    }
  }
  tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      left   = 684
      top    = 0
      width  = 950
      height = 38
    }
    markdown = "## [Workloads overview](#dashboard;id=6b38732e-d26b-45c7-b107-ed85e87ff288)"
  }
  tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      width  = 304
      height = 304
      left   = 1330
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Workloads"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "PIE"
        series {
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
          dimension {
            entity_dimension = false
            id               = "1"
            name             = "Deployment type"
          }
          metric      = "builtin:cloud.kubernetes.namespace.workloads"
          aggregation = "SUM_DIMENSIONS"
          type        = "LINE"
        }
        result_metadata {
          config {
            last_modified = 1597237249882
            custom_color  = "#008cdb"
            key           = "null¦Pod phase»Succeeded»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            last_modified = 1597234642722
            custom_color  = "#64bd64"
            key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
            last_modified = 1597234457744
            custom_color  = "#f5d30f"
          }
          config {
            last_modified = 1597234118116
            custom_color  = "#ff0000"
            key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            last_modified = 1597858600132
            custom_color  = "#ffa86c"
            key           = "null¦Deployment type»DaemonSet»falsebuiltin:cloud.kubernetes.namespace.workloads|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          }
        }
      }
    }
  }
  tile {
    configured = true
    bounds {
      left   = 988
      top    = 38
      width  = 342
      height = 304
    }
    filter {
      timeframe = "-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Running pods"
      default_name = "Custom chart"
      chart_config {
        legend = true
        type   = "TOP_LIST"
        series {
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
          dimension {
            name             = "dt.entity.cloud_application_namespace"
            entity_dimension = true
            id               = "0"
          }
          metric      = "builtin:cloud.kubernetes.namespace.runningPods"
          aggregation = "SUM_DIMENSIONS"
          type        = "LINE"
        }
        result_metadata {
          config {
            last_modified = 1597237249882
            custom_color  = "#008cdb"
            key           = "null¦Pod phase»Succeeded»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            custom_color  = "#64bd64"
            key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
            last_modified = 1597234642722
          }
          config {
            last_modified = 1597234457744
            custom_color  = "#f5d30f"
            key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          }
          config {
            custom_color  = "#ff0000"
            key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
            last_modified = 1597234118116
          }
          config {
            custom_color  = "#ffa86c"
            key           = "null¦Deployment type»DaemonSet»falsebuiltin:cloud.kubernetes.namespace.workloads|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
            last_modified = 1597858600132
          }
        }
      }
    }
    name      = ""
    tile_type = "CUSTOM_CHARTING"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `dashboard_metadata` (Block List, Max: 1) contains parameters of a dashboard (see [below for nested schema](#nestedblock--dashboard_metadata))
- `metadata` (Block List, Max: 1, Deprecated) `metadata` exists for backwards compatibility but shouldn't get specified anymore (see [below for nested schema](#nestedblock--metadata))
- `tile` (Block List) the tiles this Dashboard consist of (see [below for nested schema](#nestedblock--tile))
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--dashboard_metadata"></a>
### Nested Schema for `dashboard_metadata`

Required:

- `name` (String) the name of the dashboard
- `owner` (String) the owner of the dashboard

Optional:

- `consistent_colors` (Boolean) The tile uses consistent colors when rendering its content
- `dynamic_filters` (Block List, Max: 1) Dashboard filter configuration of a dashboard (see [below for nested schema](#nestedblock--dashboard_metadata--dynamic_filters))
- `filter` (Block List, Max: 1) Global filter Settings for the Dashboard (see [below for nested schema](#nestedblock--dashboard_metadata--filter))
- `preset` (Boolean) the dashboard is a preset (`true`) or not (`false`). Default is `false`.
- `shared` (Boolean, Deprecated) the dashboard is shared (`true`) or private (`false`)
- `sharing_details` (Block List, Max: 1, Deprecated) represents sharing configuration of a dashboard (see [below for nested schema](#nestedblock--dashboard_metadata--sharing_details))
- `tags` (Set of String) a set of tags assigned to the dashboard
- `tiles_name_size` (Number) No documentation available
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider
- `valid_filter_keys` (Set of String) a set of all possible global dashboard filters that can be applied to dashboard

<a id="nestedblock--dashboard_metadata--dynamic_filters"></a>
### Nested Schema for `dashboard_metadata.dynamic_filters`

Required:

- `filters` (Set of String) A set of all possible global dashboard filters that can be applied to a dashboard 

Currently supported values are: 

	OS_TYPE,
	SERVICE_TYPE,
	DEPLOYMENT_TYPE,
	APPLICATION_INJECTION_TYPE,
	PAAS_VENDOR_TYPE,
	DATABASE_VENDOR,
	HOST_VIRTUALIZATION_TYPE,
	HOST_MONITORING_MODE,
	KUBERNETES_CLUSTER,
	RELATED_CLOUD_APPLICATION,
	RELATED_NAMESPACE,
	TAG_KEY:<tagname>

Optional:

- `generic_tag_filters` (Block List, Max: 1) A set of generic tag filters that can be applied to a dashboard (see [below for nested schema](#nestedblock--dashboard_metadata--dynamic_filters--generic_tag_filters))
- `tag_suggestion_types` (Set of String) A set of entities applied for tag filter suggestions. You can fetch the list of possible values with the [GET all entity types](https://dt-url.net/dw03s7h)request. 

Only applicable if the **filters** set includes `TAG_KEY:<tagname>`
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

<a id="nestedblock--dashboard_metadata--dynamic_filters--generic_tag_filters"></a>
### Nested Schema for `dashboard_metadata.dynamic_filters.generic_tag_filters`

Required:

- `filter` (Block Set, Min: 1) (see [below for nested schema](#nestedblock--dashboard_metadata--dynamic_filters--generic_tag_filters--filter))

<a id="nestedblock--dashboard_metadata--dynamic_filters--generic_tag_filters--filter"></a>
### Nested Schema for `dashboard_metadata.dynamic_filters.generic_tag_filters.filter`

Required:

- `entity_types` (Set of String) Entity types affected by tag

Optional:

- `name` (String) The display name used to identify this generic filter
- `suggestions_from_entity_type` (String) The entity type for which the suggestions should be provided.
- `tag_key` (String) The tag key for this filter




<a id="nestedblock--dashboard_metadata--filter"></a>
### Nested Schema for `dashboard_metadata.filter`

Optional:

- `management_zone` (Block List) the management zone this dashboard applies to (see [below for nested schema](#nestedblock--dashboard_metadata--filter--management_zone))
- `timeframe` (String) the default timeframe of the dashboard
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

<a id="nestedblock--dashboard_metadata--filter--management_zone"></a>
### Nested Schema for `dashboard_metadata.filter.management_zone`

Required:

- `id` (String) the ID of the Dynatrace entity

Optional:

- `description` (String) a short description of the Dynatrace entity
- `name` (String) the name of the Dynatrace entity
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider



<a id="nestedblock--dashboard_metadata--sharing_details"></a>
### Nested Schema for `dashboard_metadata.sharing_details`

Optional:

- `link_shared` (Boolean) If `true`, the dashboard is shared via link and authenticated users with the link can view
- `published` (Boolean) If `true`, the dashboard is published to anyone on this environment
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider



<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- `cluster_version` (String) Dynatrace server version
- `configuration_versions` (List of Number) A Sorted list of the version numbers of the configuration
- `current_configuration_versions` (List of String) A Sorted list of the version numbers of the configuration


<a id="nestedblock--tile"></a>
### Nested Schema for `tile`

Required:

- `name` (String) the name of the tile
- `tile_type` (String) the type of the tile. Must be either `APPLICATION_WORLDMAP`, `RESOURCES`, `THIRD_PARTY_MOST_ACTIVE`, `UEM_CONVERSIONS_PER_GOAL`, `PROCESS_GROUPS_ONE` or `HOST` .

Optional:

- `assigned_entities` (Set of String) The list of Dynatrace entities, assigned to the tile
- `auto_refresh_disabled` (Boolean) Auto Refresh is disabled (`true`)
- `bounds` (Block List, Max: 1) the position and size of a tile (see [below for nested schema](#nestedblock--tile--bounds))
- `chart_visible` (Boolean)
- `configured` (Boolean) The tile is configured and ready to use (`true`) or just placed on the dashboard (`false`)
- `custom_name` (String) The name of the tile, set by user
- `exclude_maintenance_windows` (Boolean) Include (`false') or exclude (`true`) maintenance windows from availability calculations
- `filter` (Block List, Max: 1) is filter applied to a tile. It overrides dashboard's filter (see [below for nested schema](#nestedblock--tile--filter))
- `filter_config` (Block List, Max: 1) the position and size of a tile (see [below for nested schema](#nestedblock--tile--filter_config))
- `limit` (Number) The limit of the results, if not set will use the default value of the system
- `markdown` (String) The markdown-formatted content of the tile
- `metric` (String) The metric assigned to the tile
- `name_size` (String) The size of the tile name. Possible values are `small`, `medium` and `large`.
- `query` (String) A [user session query](https://www.dynatrace.com/support/help/shortlink/usql-info) executed by the tile
- `time_frame_shift` (String) The comparison timeframe of the query. If specified, you additionally get the results of the same query with the specified time shift
- `type` (String) The attribute `type` exists for backwards compatibilty. Usage is discouraged. You should use `visualization` instead.
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider
- `visualization` (String) The visualization of the tile. Possible values are: `COLUMN_CHART`, `FUNNEL`, `LINE_CHART`, `PIE_CHART`, `SINGLE_VALUE`, `TABLE`
- `visualization_config` (Block List, Max: 1) Configuration of a User session query visualization tile (see [below for nested schema](#nestedblock--tile--visualization_config))

<a id="nestedblock--tile--bounds"></a>
### Nested Schema for `tile.bounds`

Required:

- `height` (Number) the height of the tile, in pixels
- `left` (Number) the horizontal distance from the top left corner of the dashboard to the top left corner of the tile, in pixels
- `top` (Number) the vertical distance from the top left corner of the dashboard to the top left corner of the tile, in pixels
- `width` (Number) the width of the tile, in pixels

Optional:

- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider


<a id="nestedblock--tile--filter"></a>
### Nested Schema for `tile.filter`

Optional:

- `management_zone` (Block List) the management zone this tile applies to (see [below for nested schema](#nestedblock--tile--filter--management_zone))
- `timeframe` (String) the default timeframe of the tile
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

<a id="nestedblock--tile--filter--management_zone"></a>
### Nested Schema for `tile.filter.management_zone`

Required:

- `id` (String) the ID of the Dynatrace entity

Optional:

- `description` (String) a short description of the Dynatrace entity
- `name` (String) the name of the Dynatrace entity
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider



<a id="nestedblock--tile--filter_config"></a>
### Nested Schema for `tile.filter_config`

Required:

- `custom_name` (String) The name of the tile, set by user
- `default_name` (String) The default name of the tile
- `type` (String) The type of the filter. Possible values are `ALB`, `APPLICATION`, `APPLICATION_METHOD`, `APPMON`, `ASG`, `AWS_CREDENTIALS`, `AWS_CUSTOM_SERVICE`, `AWS_LAMBDA_FUNCTION`, `CLOUD_APPLICATION`, `CLOUD_APPLICATION_INSTANCE`, `CLOUD_APPLICATION_NAMESPACE`, `CONTAINER_GROUP_INSTANCE`, `CUSTOM_APPLICATION`, `CUSTOM_DEVICES`, `CUSTOM_SERVICES`, `DATABASE`, `DATABASE_KEY_REQUEST`, `DCRUM_APPLICATION`, `DCRUM_ENTITY`, `DYNAMO_DB`, `EBS`, `EC2`, `ELB`, `ENVIRONMENT`, `ESXI`, `EXTERNAL_SYNTHETIC_TEST`, `GLOBAL_BACKGROUND_ACTIVITY`, `HOST`, `IOT`, `KUBERNETES_CLUSTER`, `KUBERNETES_NODE`, `MDA_SERVICE`, `MIXED`, `MOBILE_APPLICATION`, `MONITORED_ENTITY`, `NLB`, `PG_BACKGROUND_ACTIVITY`, `PROBLEM`, `PROCESS_GROUP_INSTANCE`, `RDS`, `REMOTE_PLUGIN`, `SERVICE`, `SERVICE_KEY_REQUEST`, `SYNTHETIC_BROWSER_MONITOR`, `SYNTHETIC_HTTPCHECK`, `SYNTHETIC_HTTPCHECK_STEP`, `SYNTHETIC_LOCATION`, `SYNTHETIC_TEST`, `SYNTHETIC_TEST_STEP`, `UI_ENTITY`, `VIRTUAL_MACHINE`, `WEB_CHECK`.

Optional:

- `chart_config` (Block List, Max: 1) Configuration of a custom chart (see [below for nested schema](#nestedblock--tile--filter_config--chart_config))
- `filters` (Block List, Max: 1) Configuration of a custom chart (see [below for nested schema](#nestedblock--tile--filter_config--filters))
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

<a id="nestedblock--tile--filter_config--chart_config"></a>
### Nested Schema for `tile.filter_config.chart_config`

Required:

- `type` (String) The type of the chart

Optional:

- `axis_limits` (Map of Number) The optional custom y-axis limits
- `left_axis_custom_unit` (String) Either one of `Bit`, `BitPerHour`, `BitPerMinute`, `BitPerSecond`, `Byte`, `BytePerHour`, `BytePerMinute`, `BytePerSecond`, `Cores`, `Count`, `Day`, `DecibelMilliWatt`, `GibiByte`, `Giga`, `GigaByte`, `Hour`, `KibiByte`, `KibiBytePerHour`, `KibiBytePerMinute`, `KibiBytePerSecond`, `Kilo`, `KiloByte`, `KiloBytePerHour`, `KiloBytePerMinute`, `KiloBytePerSecond`, `MebiByte`, `MebiBytePerHour`, `MebiBytePerMinute`, `MebiBytePerSecond`, `Mega`, `MegaByte`, `MegaBytePerHour`, `MegaBytePerMinute`, `MegaBytePerSecond`, `MicroSecond`, `MilliCores`, `MilliSecond`, `MilliSecondPerMinute`, `Minute`, `Month`, `NanoSecond`, `NanoSecondPerMinute`, `NotApplicable`, `PerHour`, `PerMinute`, `PerSecond`, `Percent`, `Pixel`, `Promille`, `Ratio`, `Second`, `State`, `Unspecified`, `Week`, `Year`
- `legend` (Boolean) Defines if a legend should be shown
- `result_metadata` (Block List) Additional information about charted metric (see [below for nested schema](#nestedblock--tile--filter_config--chart_config--result_metadata))
- `right_axis_custom_unit` (String) Either one of `Bit`, `BitPerHour`, `BitPerMinute`, `BitPerSecond`, `Byte`, `BytePerHour`, `BytePerMinute`, `BytePerSecond`, `Cores`, `Count`, `Day`, `DecibelMilliWatt`, `GibiByte`, `Giga`, `GigaByte`, `Hour`, `KibiByte`, `KibiBytePerHour`, `KibiBytePerMinute`, `KibiBytePerSecond`, `Kilo`, `KiloByte`, `KiloBytePerHour`, `KiloBytePerMinute`, `KiloBytePerSecond`, `MebiByte`, `MebiBytePerHour`, `MebiBytePerMinute`, `MebiBytePerSecond`, `Mega`, `MegaByte`, `MegaBytePerHour`, `MegaBytePerMinute`, `MegaBytePerSecond`, `MicroSecond`, `MilliCores`, `MilliSecond`, `MilliSecondPerMinute`, `Minute`, `Month`, `NanoSecond`, `NanoSecondPerMinute`, `NotApplicable`, `PerHour`, `PerMinute`, `PerSecond`, `Percent`, `Pixel`, `Promille`, `Ratio`, `Second`, `State`, `Unspecified`, `Week`, `Year`
- `series` (Block List) A list of charted metrics (see [below for nested schema](#nestedblock--tile--filter_config--chart_config--series))
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

<a id="nestedblock--tile--filter_config--chart_config--result_metadata"></a>
### Nested Schema for `tile.filter_config.chart_config.result_metadata`

Optional:

- `config` (Block List) Additional metadata for charted metric (see [below for nested schema](#nestedblock--tile--filter_config--chart_config--result_metadata--config))

<a id="nestedblock--tile--filter_config--chart_config--result_metadata--config"></a>
### Nested Schema for `tile.filter_config.chart_config.result_metadata.config`

Optional:

- `custom_color` (String) The color of the metric in the chart, hex format
- `key` (String) A generated key by the Dynatrace Server
- `last_modified` (Number) The timestamp of the last metadata modification, in UTC milliseconds
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider



<a id="nestedblock--tile--filter_config--chart_config--series"></a>
### Nested Schema for `tile.filter_config.chart_config.series`

Required:

- `aggregation` (String) The charted aggregation of the metric
- `entity_type` (String) The visualization of the timeseries chart
- `metric` (String) The name of the charted metric
- `type` (String) The visualization of the timeseries chart. Possible values are `AREA`, `BAR` and `LINE`.

Optional:

- `aggregation_rate` (String)
- `dimension` (Block List) Configuration of the charted metric splitting (see [below for nested schema](#nestedblock--tile--filter_config--chart_config--series--dimension))
- `percentile` (Number) The charted percentile. Only applicable if the **aggregation** is set to `PERCENTILE`
- `sort_ascending` (Boolean) Sort ascending (`true`) or descending (`false`)
- `sort_column` (Boolean) Sort the column (`true`) or (`false`)
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider

<a id="nestedblock--tile--filter_config--chart_config--series--dimension"></a>
### Nested Schema for `tile.filter_config.chart_config.series.dimension`

Required:

- `id` (String) The ID of the dimension by which the metric is split

Optional:

- `entity_dimension` (Boolean)
- `name` (String) The name of the dimension by which the metric is split
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider
- `values` (Set of String) The splitting value




<a id="nestedblock--tile--filter_config--filters"></a>
### Nested Schema for `tile.filter_config.filters`

Optional:

- `filter` (Block List) the tiles this Dashboard consist of (see [below for nested schema](#nestedblock--tile--filter_config--filters--filter))

<a id="nestedblock--tile--filter_config--filters--filter"></a>
### Nested Schema for `tile.filter_config.filters.filter`

Required:

- `entity_type` (String) The entity type (e.g. HOST, SERVICE, ...)

Optional:

- `match` (Block List) the tiles this Dashboard consist of (see [below for nested schema](#nestedblock--tile--filter_config--filters--filter--match))

<a id="nestedblock--tile--filter_config--filters--filter--match"></a>
### Nested Schema for `tile.filter_config.filters.filter.match`

Required:

- `key` (String) The entity type (e.g. HOST, SERVICE, ...)

Optional:

- `values` (Set of String) the tiles this Dashboard consist of





<a id="nestedblock--tile--visualization_config"></a>
### Nested Schema for `tile.visualization_config`

Optional:

- `has_axis_bucketing` (Boolean) The axis bucketing when enabled groups similar series in the same virtual axis
- `unknowns` (String) allows for configuring properties that are not explicitly supported by the current version of this provider
 