resource "dynatrace_span_entry_point" "#name#" {
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
