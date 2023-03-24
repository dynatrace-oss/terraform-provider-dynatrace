resource "dynatrace_extension_execution_controller" "#name#" {
  enabled             = true
  ingest_active       = false
  performance_profile = "DEFAULT"
  scope               = "environment"
  statsd_active       = false
}