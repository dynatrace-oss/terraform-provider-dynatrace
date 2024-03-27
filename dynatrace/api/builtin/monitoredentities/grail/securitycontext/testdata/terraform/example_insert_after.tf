resource "dynatrace_grail_security_context" "first-instance" {
  entity_type          = "exampletype"
  destination_property = "exampleproperty"
}

resource "dynatrace_grail_security_context" "second-instance" {
  entity_type          = "exampletype"
  destination_property = "exampleproperty"
  insert_after         = dynatrace_grail_security_context.first-instance.id
}
