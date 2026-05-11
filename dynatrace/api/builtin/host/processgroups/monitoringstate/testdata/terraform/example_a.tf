data "dynatrace_entity" "process_group" {
  entity_selector = "type(\"PROCESS_GROUP\")"
}

resource "dynatrace_host_process_group_monitoring" "monitoring" {
  host_id          = "HOST-1234567890000000"
  monitoring_state = "MONITORING_ON"
  process_group    = data.dynatrace_entity.process_group.id
}
