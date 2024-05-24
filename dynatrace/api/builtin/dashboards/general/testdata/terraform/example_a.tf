resource "dynatrace_dashboards_general" "#name#" {
  enable_public_sharing = false
  default_dashboard_list {
    default_dashboard {
      dashboard = "41eae96d-4930-4f44-bbd8-3699f21a8bbf"
      user_group = "d0c2d3e3-c1b4-456a-b0ce-c560273f1488"
    }
  }
}