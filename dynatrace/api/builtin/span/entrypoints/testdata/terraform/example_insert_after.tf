resource "dynatrace_span_entry_point" "first-instance" {
  name = "#name#"
  action = "DONT_CREATE_ENTRYPOINT"
  matches {
    match {
      case_sensitive = true
      comparison = "DOES_NOT_CONTAIN"
      source = "SPAN_NAME"
      value = "foo"
    }
  }
}

resource "dynatrace_span_entry_point" "second-instance" {
  name = "#name#-2"
  action = "DONT_CREATE_ENTRYPOINT"
  matches {
    match {
      case_sensitive = true
      comparison = "DOES_NOT_CONTAIN"
      source = "SPAN_NAME"
      value = "foo-2"
    }
  }
  insert_after = dynatrace_span_entry_point.first-instance.id
}
