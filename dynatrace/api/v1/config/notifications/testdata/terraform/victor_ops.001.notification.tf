resource "dynatrace_notification" "#name#" {
  victor_ops {
    name             = "#name#"
    active           = false
    alerting_profile = dynatrace_alerting_profile.Default.id
    api_key          = "#######"
    message          = "victor-ops-message"
    routing_key      = "victor-ops-routing-key"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
