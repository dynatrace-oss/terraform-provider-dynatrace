variable "PROCESS_GROUP_ID" {
  type = string
}

resource "dynatrace_process_group_monitoring" "monitoring" {
  monitoring_state = "MONITORING_ON"
  process_group_id = var.PROCESS_GROUP_ID
}
