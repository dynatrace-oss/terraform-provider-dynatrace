resource "dynatrace_dashboards_general" "#name#" {
  enable_public_sharing = false
  default_dashboard_list {
    default_dashboard {
      dashboard = "00000000-0000-0000-0000-000000000000"
      user_group = "terraform"
    }
  }
}