resource "dynatrace_span_capture_rule" "first-instance" {
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

resource "dynatrace_span_capture_rule" "second-instance" {
  name = "#name#-2"
  action = "IGNORE"
  matches {
    match {
      comparison = "EQUALS"
      source = "SPAN_NAME"
      value = "foo"
    }
  }
  insert_after = dynatrace_span_capture_rule.first-instance.id
}
