resource "dynatrace_span_context_propagation" "first-instance" {
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

resource "dynatrace_span_context_propagation" "second-instance" {
  name = "#name#-2"
  action = "PROPAGATE"
  matches {
    match {
      comparison = "EQUALS"
      source = "SPAN_NAME"
      value = "asdf-2"
    }
  }
  insert_after = dynatrace_span_context_propagation.first-instance.id
}