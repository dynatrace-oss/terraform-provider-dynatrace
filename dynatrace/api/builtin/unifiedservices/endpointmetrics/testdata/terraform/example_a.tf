resource "dynatrace_unified_services_metrics" "#name#" {
  enable_endpoint_metrics = true
  service_id              = "environment"
}
