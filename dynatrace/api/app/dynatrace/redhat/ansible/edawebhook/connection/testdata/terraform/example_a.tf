resource "dynatrace_event_driven_ansible_connections" "connection" {
  event_stream_enabled = true
  name    = "#name#"
  url     = "https://www.example.com"
  type    = "api-token"
  token   = "######"
}
