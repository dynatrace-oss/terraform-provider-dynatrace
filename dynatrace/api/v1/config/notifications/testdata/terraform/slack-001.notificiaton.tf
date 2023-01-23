resource "dynatrace_notification" "#name#" {
  slack {
    name             = "#name#"
    active           = true
    alerting_profile = "f75e68ef-aca7-3a07-9c21-94eb00ecfc56"
    channel          = "#dynatrace-critical"
    title            = "Test2"
    url              = "https://www.google.at/f75e68ef-aca7-3a07-9c21-94eb00ecfc56"
  }
}