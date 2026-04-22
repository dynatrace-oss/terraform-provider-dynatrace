resource "dynatrace_host_monitoring_mode" "mode" {
  host_id         = data.dynatrace_entities.hosts.entities[0].entity_id
  monitoring_mode = "FULL_STACK"
}

data "dynatrace_entities" "hosts" {
  type = "HOST"
}
