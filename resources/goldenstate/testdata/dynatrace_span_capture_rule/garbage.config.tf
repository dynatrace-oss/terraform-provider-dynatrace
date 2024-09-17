resource "dynatrace_span_capture_rule" "team-hawaiian-deleteme" {
  name   = "team-hawaiian-deleteme"
  action = "IGNORE"
  matches {
    match {
      comparison = "EQUALS"
      source     = "SPAN_NAME"
      value      = "team-hawaiian-deleteme"
    }
  }
}