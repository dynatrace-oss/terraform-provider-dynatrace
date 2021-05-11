resource "dynatrace_dashboard" "#name#" {
  dashboard_metadata {
    name   = "#name#"
    shared = false
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
