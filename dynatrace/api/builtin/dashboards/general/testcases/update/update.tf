resource "dynatrace_dashboards_general" "general" {
  enable_public_sharing = false
  default_dashboard_list {
    default_dashboard {
      # updated to use a new just created dashboard
      dashboard = dynatrace_dashboard.dashboard_new.id
      user_group = dynatrace_iam_group.group.id
    }
  }
}

resource "dynatrace_iam_group" "group" {
  name = "#name#"
}

resource "dynatrace_dashboard" "dashboard_new" {
  dashboard_metadata {
    name   = "#name#-new"
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
