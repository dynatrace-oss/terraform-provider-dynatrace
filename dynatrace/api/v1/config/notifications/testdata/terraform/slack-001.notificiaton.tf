resource "dynatrace_notification" "#name#" {
  slack {
    name             = "#name#"
    active           = true
    alerting_profile = dynatrace_alerting_profile.Default.id
    channel          = "#dynatrace-critical"
    title            = "Test2"
    url              = "https://www.google.at/#name#"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
