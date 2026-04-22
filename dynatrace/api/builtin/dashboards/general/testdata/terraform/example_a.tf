resource "dynatrace_dashboards_general" "general" {
  depends_on = [time_sleep.dashboard_create]
  enable_public_sharing = false
  default_dashboard_list {
    default_dashboard {
      dashboard = dynatrace_dashboard.dashboard.id
      user_group = dynatrace_iam_group.group.id
    }
  }
}

# the datasource of dynatrace_dashboards_general may not be updated yet.
resource "time_sleep" "dashboard_create" {
  depends_on = [dynatrace_dashboard.dashboard]
  create_duration = "5s"
}

resource "dynatrace_iam_group" "group" {
  name = "#name#"
}

resource "dynatrace_dashboard" "dashboard" {
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
}
