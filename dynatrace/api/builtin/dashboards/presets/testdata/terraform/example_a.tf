resource "dynatrace_dashboards_presets" "presets" {
  enable_dashboard_presets = true
  dashboard_presets_list {
    dashboard_presets {
      dashboard_preset = dynatrace_dashboard.dashboard.id
      user_group = dynatrace_iam_group.group.id
    }
  }
}

resource "dynatrace_iam_group" "group" {
  name = "#name#"
}

resource "dynatrace_dashboard" "dashboard" {
  dashboard_metadata {
    preset = true
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
}
