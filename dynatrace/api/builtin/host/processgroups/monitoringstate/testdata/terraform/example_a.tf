variable "PROCESS_GROUP_ID" {
  type = string
}

resource "dynatrace_host_process_group_monitoring" "monitoring" {
  host_id          = "HOST-1234567890000000"
  monitoring_state = "MONITORING_ON"
  process_group    = var.PROCESS_GROUP_ID
}
