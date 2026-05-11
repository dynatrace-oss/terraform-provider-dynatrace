data "dynatrace_entity" "process_group" {
  entity_selector = "type(\"PROCESS_GROUP\")"
}

resource "dynatrace_devobs_agent_optin" "optin" {
  scope   = data.dynatrace_entity.process_group.id
  enabled = false
}
