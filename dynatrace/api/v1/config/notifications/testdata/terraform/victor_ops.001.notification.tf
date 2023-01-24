resource "dynatrace_notification" "#name#" {
  victor_ops {
    name             = "#name#"
    active           = false
    alerting_profile = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    api_key          = "#######"
    message          = "victor-ops-message"
    routing_key      = "victor-ops-routing-key"
  }
}
