resource "dynatrace_dashboard" "dashboard" {
  dashboard_metadata {
    name   = "#name#"
    owner  = "Dynatrace"
  }
  tile {
    name      = "my tile"
    tile_type = "HOST"
  }
}

resource "dynatrace_report" "report" {
  type                = "DASHBOARD"
  dashboard_id        = dynatrace_dashboard.dashboard.id
  email_notifications = true
  subscriptions {
    month = ["terraform1@dynatrace.com", "terraform2@dynatrace.com"]
    week = ["terraform3@dynatrace.com", "terraform4@dynatrace.com", "terraform5@dynatrace.com"]
  }
}
