data "dynatrace_entity" "process_group" {
  entity_selector = "type(\"PROCESS_GROUP\")"
}

resource "dynatrace_connectivity_alerts" "alert" {
  connectivity_alerts = false
  process_group_id    = data.dynatrace_entity.process_group.id
}
