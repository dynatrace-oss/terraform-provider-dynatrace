resource "dynatrace_host_process_group_monitoring" "#name#" {
  host_id          = "HOST-1234567890000000"
  monitoring_state = "MONITORING_ON"
  process_group    = "PROCESS_GROUP-1234567890000000"
}