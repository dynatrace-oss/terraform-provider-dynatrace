resource "dynatrace_span_entry_point" "#name#" {
  name = "#name#"
  action = "CREATE_ENTRYPOINT"
  matches {
    match {
      comparison = "EQUALS"
      source = "ATTRIBUTE"
      key = "asdf"
      value = "CLIENT"
    }
  }
}

