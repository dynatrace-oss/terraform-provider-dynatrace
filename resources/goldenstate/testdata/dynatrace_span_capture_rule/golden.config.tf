resource "dynatrace_span_capture_rule" "team-hawaiian-toast" {
  name   = "team-hawaiian-toast"
  action = "IGNORE"
  matches {
    match {
      comparison = "EQUALS"
      source     = "SPAN_NAME"
      value      = "team-hawaiian-toast"
    }
  }
}