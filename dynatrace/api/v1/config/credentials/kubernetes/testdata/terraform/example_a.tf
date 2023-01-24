resource "dynatrace_k8s_credentials" "#name#" {
  label                        = "#name#"
  auth_token                   = "XXXXXXXX"
  active                       = false
  certificate_check_enabled    = false
  endpoint_url                 = "https://${randomize}.${randomize}.com:6443/"
  events_integration_enabled   = true
  hostname_verification        = false
  prometheus_exporters         = false
  workload_integration_enabled = true
  event_analysis_and_alerting_enabled = true

  events_field_selectors {
    active         = true
    field_selector = "involvedObject.kind=Node"
    label          = "Node events"
  }
}
