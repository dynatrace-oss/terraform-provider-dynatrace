resource "dynatrace_network_monitor_outage" "#name#" {
  global_outages = true
  global_consecutive_outage_count_threshold = 5
  local_outages = true
  local_consecutive_outage_count_threshold = 3
  local_location_outage_count_threshold = 1
}