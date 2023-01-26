resource "dynatrace_pg_alerting" "#name#" {
  enabled                    = true
  alerting_mode              = "ON_INSTANCE_COUNT_VIOLATION"
  minimum_instance_threshold = 5
  process_group              = "PROCESS_GROUP-XXXXXXXXXXXXXXXX"
}