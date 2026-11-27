resource "dynatrace_synthetic_v2_configuration" "#name#" {
  bm_monitor_timeout                             = 120000
  bm_step_timeout                                = 60000
  maintenance_windows_excluded_from_availability = true
}
