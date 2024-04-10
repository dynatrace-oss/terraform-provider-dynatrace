resource "dynatrace_dashboards_presets" "#name#" {
  enable_dashboard_presets = true
  dashboard_presets_list {
    dashboard_presets {
      dashboard_preset = "41eae96d-4930-4f44-bbd8-3699f21a8bbf"
      user_group = "d0c2d3e3-c1b4-456a-b0ce-c560273f1488"
    }
  }
}