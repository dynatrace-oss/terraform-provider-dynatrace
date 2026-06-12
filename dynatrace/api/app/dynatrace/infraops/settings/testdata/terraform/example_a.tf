resource "dynatrace_infraops_app_settings" "example" {
  show_monitoring_candidates     = true
  show_standalone_hosts          = true
  interface_saturation_threshold = 0.95
  invex_dql_query_limit          = 1000
}
