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

resource "dynatrace_iam_group" "group" {
  count = 2
  name = "#name#-${count.index}"
}

resource "dynatrace_dashboard_sharing" "sharing" {
  dashboard_id = dynatrace_dashboard.dashboard.id
  permissions {
    permission {
      level = "VIEW"
      type  = "GROUP"
      id = dynatrace_iam_group.group[0].id
    }
    # to update
    permission {
      level = "VIEW"
      type  = "GROUP"
      id = dynatrace_iam_group.group[1].id
    }
  }
}
