resource "dynatrace_web_app_beacon_endpoint" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  type           = "ONEAGENT"
  url            = "/beacon"
  use_cors       = true
}