resource "dynatrace_span_capture_rule" "#name#" {
  name = "#name#"
  action = "IGNORE"
  matches {
    match {
      comparison = "EQUALS"
      source = "SPAN_NAME"
      value = "foo"
    }
  }
}
