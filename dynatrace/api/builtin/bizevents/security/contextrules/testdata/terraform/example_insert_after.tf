resource "dynatrace_business_events_security_context" "first-instance" {
  security_context_rule {
    query              = "matchesPhrase(content, \"#name#\")"
    rule_name          = "#name#"
    value_source_field = "#name#"
    value_source       = "FIELD"
  }
}

resource "dynatrace_business_events_security_context" "second-instance" {
  security_context_rule {
    query              = "matchesPhrase(content, \"#name#\")"
    rule_name          = "#name#-second"
    value_source_field = "#name#-second"
    value_source       = "FIELD"
  }
  insert_after = dynatrace_business_events_security_context.first-instance.id
}
