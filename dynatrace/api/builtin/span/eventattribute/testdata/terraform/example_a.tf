resource "dynatrace_span_events" "#name#" {
  key     = "exception.terraform"
  masking = "NOT_MASKED"
}