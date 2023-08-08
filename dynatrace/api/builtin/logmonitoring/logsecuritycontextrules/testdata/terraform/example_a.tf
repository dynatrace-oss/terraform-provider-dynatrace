resource "dynatrace_log_security_context" "#name#" {
  security_context_rule {
    query        = "matchesPhrase(content, \"#name#\")"
    rule_name    = "#name#"
    value_source_field = "#name#"
    value_source = "FIELD"
  }
}
