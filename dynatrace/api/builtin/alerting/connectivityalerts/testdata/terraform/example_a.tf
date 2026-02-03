variable "PROCESS_GROUP_ID" {
  type = string
}

resource "dynatrace_connectivity_alerts" "alert" {
  connectivity_alerts = false
  process_group_id    = var.PROCESS_GROUP_ID
}
