resource "dynatrace_dashboard" "#name#" {
  metadata {
    cluster_version        = "1.206.95.20201116-094826"
    configuration_versions = [3]
  }
  dashboard_metadata {
    name   = "#name#"
    shared = true
    owner  = "reinhard.pilz@dynatrace.com"
    sharing_details {
      link_shared = true
      published   = true
    }
    dashboard_filter {
      timeframe = ""
    }
    tags = ["Configurator"]
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 0
      left   = 0
      width  = 342
      height = 38
    }
    tile_filter {
    }
    markdown = "###Cluster Overview   - [Cluster Insights](/ui/kubernetes/)"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 0
      left   = 342
      width  = 266
      height = 38
    }
    tile_filter {
    }
    markdown = "###Cluster Utilization"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 0
      left   = 1064
      width  = 266
      height = 38
    }
    tile_filter {
    }
    markdown = "###Pods phases "
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 0
      left   = 608
      width  = 456
      height = 38
    }
    tile_filter {
    }
    markdown = "###Cluster Workloads \u0026 Namespaces"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 0
      left   = 1330
      width  = 304
      height = 38
    }
    tile_filter {
    }
    markdown = "Running vs. desired pods"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 342
      left   = 0
      width  = 304
      height = 38
    }
    tile_filter {
    }
    markdown = "###[Cluster utilization](#dashboard;id=bbbbbbbb-0001-0000-0000-000000000001;)"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 342
      left   = 304
      width  = 304
      height = 38
    }
    tile_filter {
    }
    markdown = "###[Resource Quotas](#dashboard;id=bbbbbbbb-0001-0000-0000-000000000003;)\n"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 342
      left   = 608
      width  = 304
      height = 38
    }
    tile_filter {
    }
    markdown = "###[Container usage \u0026 health](#dashboard;id=bbbbbbbb-0001-0000-0000-000000000002;)"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 342
      left   = 912
      width  = 304
      height = 38
    }
    tile_filter {
    }
    markdown = "###[Performance Engineering](#dashboard;id=bbbbbbbb-0001-0000-0000-000000000004;)\n"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 342
      left   = 1216
      width  = 304
      height = 38
    }
    tile_filter {
    }
    markdown = "###[User Experience](#dashboard;id=bbbbbbbb-0001-0000-0000-000000000005;)\n"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 380
      left   = 0
      width  = 304
      height = 228
    }
    tile_filter {
    }
    markdown = "_____________________\nSee the Kubernetes cluster utilization. CPU and Memory Request and limits over time for all nodes and splitted by namespaces.\n"
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 380
      left   = 304
      width  = 304
      height = 228
    }
    tile_filter {
    }
    markdown = "_____________________\nGet an overview and understanding of the Kubernetes resource quotas (Memory and CPU) assigned to your namespaces and its usage. "
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 380
      left   = 608
      width  = 304
      height = 228
    }
    tile_filter {
    }
    markdown = "_____________________\nUnderstand the health and phases of your Pods in your clusters. Their memory and cpu usage, which pods are throttled, have failed or are pending to be scheduled. Also check if you have Out-of-memory killed containers."
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 380
      left   = 912
      width  = 304
      height = 228
    }
    tile_filter {
    }
    markdown = "_____________________\nGive your developers and SRE engineers all they need to understand and improve the performance of each app, pod and each transaction on your clusters. View the response time percentiles, slow transactions, database executions per microservice, its network usage and more. Filter the transactions by App label, namespace and much more."
  }
  markdown_tile {
    name       = "Markdown"
    tile_type  = "MARKDOWN"
    configured = true
    bounds {
      top    = 380
      left   = 1216
      width  = 304
      height = 228
    }
    tile_filter {
    }
    markdown = "_____________________\nAre your endusers satisfied? how is the engagement, experience and user behaviour of your applications? Get the insights of all your applications and users in an instance."
  }
  filterable_entity_tile {
    name       = ""
    tile_type  = "HOSTS"
    configured = true
    bounds {
      top    = 38
      left   = 152
      width  = 152
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "HOST"
      custom_name  = "Full-Stack Kubernetes nodes"
      default_name = "Full-Stack Kubernetes nodes"
      chart_config {
        legend_shown = true
        type         = "TIMESERIES"
      }
      filters_per_entity_type {
        key = "HOST"
        values {
          key    = "HOST_SOFTWARE_TECH"
          values = ["KUBERNETES"]
        }
      }
    }
    chart_visible = true
  }
  custom_charting_tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 456
      width  = 152
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "CPU available [last 5 min]"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "SINGLE_VALUE"
        series {
          metric           = "builtin:cloud.kubernetes.cluster.cpuAvailable"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
    }
  }
  custom_charting_tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 190
      left   = 456
      width  = 152
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Memory available [last 5 min]"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "SINGLE_VALUE"
        series {
          metric           = "builtin:cloud.kubernetes.cluster.memoryAvailable"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
    }
  }
  custom_charting_tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 304
      width  = 152
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Cores"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "SINGLE_VALUE"
        series {
          metric           = "builtin:cloud.kubernetes.cluster.cores"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
    }
  }
  custom_charting_tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 190
      left   = 304
      width  = 152
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Memory total"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "SINGLE_VALUE"
        series {
          metric           = "builtin:cloud.kubernetes.cluster.memory"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "KUBERNETES_CLUSTER"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
      }
    }
  }
  custom_charting_tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 1064
      width  = 266
      height = 304
    }
    tile_filter {
      timeframe = "now-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "🚦Pods phases"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "PIE"
        series {
          metric      = "builtin:cloud.kubernetes.workload.pods"
          aggregation = "SUM_DIMENSIONS"
          type        = "LINE"
          entity_type = "CLOUD_APPLICATION"
          dimensions {
            id               = "1"
            name             = "Pod phase"
            values           = []
            entity_dimension = false
          }
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594115586538
          custom_color  = "#f5d30f"
        }
        result_metadata {
          key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594141365181
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594025904142
          custom_color  = "#ff0000"
        }
      }
    }
  }
  custom_charting_tile {
    name       = "Custom chart"
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 0
      width  = 152
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Nodes / Cluster"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "PIE"
        series {
          metric      = "builtin:cloud.kubernetes.cluster.nodes"
          aggregation = "SUM_DIMENSIONS"
          type        = "LINE"
          entity_type = "KUBERNETES_CLUSTER"
          dimensions {
            id               = "0"
            name             = "dt.entity.kubernetes_cluster"
            values           = []
            entity_dimension = true
          }
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "KUBERNETES_CLUSTER-EA9EB67E9CBAE0CD¦KUBERNETES_CLUSTER»KUBERNETES_CLUSTER-EA9EB67E9CBAE0CD»truebuiltin:cloud.kubernetes.cluster.nodes|SUM_DIMENSIONS|TOTAL|LINE|KUBERNETES_CLUSTER"
          last_modified = 1594146638933
          custom_color  = "#b4e5f9"
        }
        result_metadata {
          key           = "KUBERNETES_CLUSTER-5E463319734AB4DD¦KUBERNETES_CLUSTER»KUBERNETES_CLUSTER-5E463319734AB4DD»truebuiltin:cloud.kubernetes.cluster.nodes|SUM_DIMENSIONS|TOTAL|LINE|KUBERNETES_CLUSTER"
          last_modified = 1594146636359
          custom_color  = "#008cdb"
        }
        result_metadata {
          key           = "KUBERNETES_CLUSTER-FA15B65ACE980EAD¦KUBERNETES_CLUSTER»KUBERNETES_CLUSTER-FA15B65ACE980EAD»truebuiltin:cloud.kubernetes.cluster.nodes|SUM_DIMENSIONS|TOTAL|LINE|KUBERNETES_CLUSTER"
          last_modified = 1594146627924
          custom_color  = "#2ab6f4"
        }
      }
    }
  }
  custom_charting_tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 190
      left   = 1330
      width  = 304
      height = 152
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Desired vs Running pods"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "TIMESERIES"
        series {
          metric           = "builtin:cloud.kubernetes.namespace.desiredPods"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = false
          aggregation_rate = "TOTAL"
        }
        series {
          metric           = "builtin:cloud.kubernetes.namespace.runningPods"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "nullbuiltin:cloud.kubernetes.namespace.desiredPods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          last_modified = 1594147642285
          custom_color  = "#ff0000"
        }
        result_metadata {
          key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594115586538
          custom_color  = "#f5d30f"
        }
        result_metadata {
          key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594141365181
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594025904142
          custom_color  = "#ff0000"
        }
        result_metadata {
          key           = "nullbuiltin:cloud.kubernetes.namespace.runningPods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          last_modified = 1594147621067
          custom_color  = "#64bd64"
        }
        left_axis_custom_unit = "Count"
      }
    }
  }
  custom_charting_tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 608
      width  = 228
      height = 304
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Workloads"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "TOP_LIST"
        series {
          metric      = "builtin:cloud.kubernetes.namespace.workloads"
          aggregation = "SUM_DIMENSIONS"
          type        = "LINE"
          entity_type = "CLOUD_APPLICATION_NAMESPACE"
          dimensions {
            id               = "1"
            name             = "Deployment type"
            values           = []
            entity_dimension = false
          }
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594115586538
          custom_color  = "#f5d30f"
        }
        result_metadata {
          key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594141365181
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594025904142
          custom_color  = "#ff0000"
        }
      }
    }
  }
  custom_charting_tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 836
      width  = 228
      height = 304
    }
    tile_filter {
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "Pods by namespace"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "TOP_LIST"
        series {
          metric      = "builtin:cloud.kubernetes.namespace.workloads"
          aggregation = "SUM_DIMENSIONS"
          type        = "LINE"
          entity_type = "CLOUD_APPLICATION_NAMESPACE"
          dimensions {
            id               = "0"
            name             = "dt.entity.cloud_application_namespace"
            values           = []
            entity_dimension = true
          }
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594115586538
          custom_color  = "#f5d30f"
        }
        result_metadata {
          key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594141365181
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594025904142
          custom_color  = "#ff0000"
        }
      }
    }
  }
  custom_charting_tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 1330
      width  = 152
      height = 152
    }
    tile_filter {
      timeframe = "now-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "desired pods"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "SINGLE_VALUE"
        series {
          metric           = "builtin:cloud.kubernetes.namespace.desiredPods"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = false
          aggregation_rate = "TOTAL"
        }
        series {
          metric           = "builtin:cloud.kubernetes.namespace.runningPods"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594025904142
          custom_color  = "#ff0000"
        }
        result_metadata {
          key           = "nullbuiltin:cloud.kubernetes.namespace.runningPods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          last_modified = 1594147621067
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "nullbuiltin:cloud.kubernetes.namespace.desiredPods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          last_modified = 1594147642285
          custom_color  = "#ff0000"
        }
        result_metadata {
          key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594115586538
          custom_color  = "#f5d30f"
        }
        result_metadata {
          key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594141365181
          custom_color  = "#64bd64"
        }
        left_axis_custom_unit = "Count"
      }
    }
  }
  custom_charting_tile {
    name       = ""
    tile_type  = "CUSTOM_CHARTING"
    configured = true
    bounds {
      top    = 38
      left   = 1482
      width  = 152
      height = 152
    }
    tile_filter {
      timeframe = "now-5m"
    }
    filter_config {
      type         = "MIXED"
      custom_name  = "running pods"
      default_name = "Custom chart"
      chart_config {
        legend_shown = true
        type         = "SINGLE_VALUE"
        series {
          metric           = "builtin:cloud.kubernetes.namespace.runningPods"
          aggregation      = "SUM_DIMENSIONS"
          type             = "LINE"
          entity_type      = "CLOUD_APPLICATION_NAMESPACE"
          sort_ascending   = false
          sort_column      = true
          aggregation_rate = "TOTAL"
        }
        result_metadata {
          key           = "null¦Pod phase»Running»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594141365181
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "null¦Pod phase»Failed»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594025904142
          custom_color  = "#ff0000"
        }
        result_metadata {
          key           = "nullbuiltin:cloud.kubernetes.namespace.runningPods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          last_modified = 1594147621067
          custom_color  = "#64bd64"
        }
        result_metadata {
          key           = "nullbuiltin:cloud.kubernetes.namespace.desiredPods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION_NAMESPACE"
          last_modified = 1594147642285
          custom_color  = "#ff0000"
        }
        result_metadata {
          key           = "null¦Pod phase»Pending»falsebuiltin:cloud.kubernetes.workload.pods|SUM_DIMENSIONS|TOTAL|LINE|CLOUD_APPLICATION"
          last_modified = 1594115586538
          custom_color  = "#f5d30f"
        }
        left_axis_custom_unit = "Count"
      }
    }
  }
  abstract_tile {
    name       = "Smartscape"
    tile_type  = "PURE_MODEL"
    configured = true
    bounds {
      top    = 190
      left   = 0
      width  = 304
      height = 152
    }
    tile_filter {
    }
  }
}

