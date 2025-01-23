resource "dynatrace_business_events_capturing_variants" "first-instance" {
  content_type_matcher = "EQUALS"
  content_type_value   = "TerraformFirst"
  parser               = "Text"
  scope                = "environment"
}


resource "dynatrace_business_events_capturing_variants" "second-instance" {
  content_type_matcher = "EQUALS"
  content_type_value   = "TerraformSecond"
  parser               = "Raw"
  scope                = "environment"
  insert_after = dynatrace_business_events_capturing_variants.first-instance.id
}
