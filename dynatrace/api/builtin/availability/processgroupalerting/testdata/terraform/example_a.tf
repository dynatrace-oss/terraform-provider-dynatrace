data "dynatrace_entity" "process_group" {
  entity_selector = "type(\"PROCESS_GROUP\")"
}

resource "dynatrace_pg_alerting" "alert" {
  enabled                    = true
  alerting_mode              = "ON_INSTANCE_COUNT_VIOLATION"
  minimum_instance_threshold = 5
  process_group              = data.dynatrace_entity.process_group.id
}
