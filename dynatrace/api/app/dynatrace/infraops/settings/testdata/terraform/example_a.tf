resource "dynatrace_infraops_app_settings" "#name#" {
  show_monitoring_candidates = true
  show_standalone_hosts      = true
  interface_saturation_threshold = 0.95
}