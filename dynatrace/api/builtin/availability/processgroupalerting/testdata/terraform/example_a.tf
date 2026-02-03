variable "PROCESS_GROUP_ID" {
  type = string
}

resource "dynatrace_pg_alerting" "alert" {
  enabled                    = true
  alerting_mode              = "ON_INSTANCE_COUNT_VIOLATION"
  minimum_instance_threshold = 5
  process_group              = var.PROCESS_GROUP_ID
}
