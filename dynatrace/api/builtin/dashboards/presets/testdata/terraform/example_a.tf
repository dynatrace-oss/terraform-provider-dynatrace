resource "dynatrace_dashboards_presets" "#name#" {
  enable_dashboard_presets = true
  dashboard_presets_list {
    dashboard_presets {
      dashboard_preset = "00000000-0000-0000-0000-000000000000"
      user_group = "00000000-0000-0000-0000-000000000000"
    }
  }
}