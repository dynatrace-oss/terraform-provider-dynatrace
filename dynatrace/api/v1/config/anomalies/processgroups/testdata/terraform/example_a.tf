data "dynatrace_entity" "process_group" {
  entity_selector = "type(\"PROCESS_GROUP\")"
}

resource "dynatrace_pg_anomalies" "anomaly" {
  pg_id = data.dynatrace_entity.process_group.id
  availability {
    method            = "OFF"
    minimum_threshold = 0
  }
}
