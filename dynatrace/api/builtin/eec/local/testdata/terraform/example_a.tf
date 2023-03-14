resource "dynatrace_extension_execution_controller" "#name#" {
  enabled             = true
  ingest_active       = true
  performance_profile = "DEFAULT"
  scope               = "environment"
  statsd_active       = true
}