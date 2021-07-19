resource "dynatrace_span_context_propagation" "#name#" {
  name = "#name#"
  action = "PROPAGATE"
  matches {
    match {
      comparison = "EQUALS"
      source = "SPAN_NAME"
      value = "asdf"
    }
  }
}