resource "dynatrace_event_driven_ansible_connections" "#name#"{
  event_stream_enabled = true
  name    = "#name#"
  url     = "https://www.google.com"
  type    = "api-token"
  token   = "######"
  }