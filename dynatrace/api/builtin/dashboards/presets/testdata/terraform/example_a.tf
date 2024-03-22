resource "dynatrace_dashboards_presets" "#name#" {
  enable_dashboard_presets = true
  dashboard_presets_list {
    dashboard_presets {
      dashboard_preset = "${dynatrace_dashboard.dashboard.id}"
      user_group = "d0c2d3e3-c1b4-456a-b0ce-c560273f1488"
    }
  }
}

resource "dynatrace_dashboard" "dashboard" {
  dashboard_metadata {
    name   = "#name#"
    shared = false
    preset = true
    owner  = "Dynatrace"
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
    markdown = "## Terraform Test"
  }
}