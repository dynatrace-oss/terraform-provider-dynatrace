resource "dynatrace_business_events_capturing_variants" "#name#" {
  content_type_matcher = "EQUALS"
  content_type_value   = "#name#"
  parser               = "Text"
  scope                = "environment"
}
