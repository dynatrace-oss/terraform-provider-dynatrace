resource "dynatrace_synthetic_v2_configuration" "configuration" {
  bm_monitor_timeout                             = 180000
  bm_step_timeout                                = 90000
  maintenance_windows_excluded_from_availability = false
}
