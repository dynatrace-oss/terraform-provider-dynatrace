resource "dynatrace_synthetic_v2_location" "example" {
  name                                  = "#name#"
  nam_execution_supported               = true
  browser_execution_supported           = true
  auto_update_chromium                  = true
  availability_location_outage          = true
  availability_node_outage              = true
  availability_notifications_enabled    = true
  use_new_kubernetes_version            = true
  fips_mode                             = "ENABLED"
  city                                  = "San Francisco de Asis"
  country_code                          = "VE"
  deployment_type                       = "STANDARD"
  region_code                           = "04"
  latitude                              = 10.0758
  location_node_outage_delay_in_minutes = 3
  longitude                             = -67.5442
  min_active_gate_count                 = 1
  max_active_gate_count                 = 2
  status                                = "ENABLED"
}
