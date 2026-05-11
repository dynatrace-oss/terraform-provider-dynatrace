data "dynatrace_entity" "process_group" {
  entity_selector = "type(\"PROCESS_GROUP\")"
}

resource "dynatrace_process_group_monitoring" "monitoring" {
  monitoring_state = "MONITORING_ON"
  process_group_id = data.dynatrace_entity.process_group.id
}
