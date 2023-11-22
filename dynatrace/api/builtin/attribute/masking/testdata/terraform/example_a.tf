resource "dynatrace_attribute_masking" "#name#" {
  enabled = true
  key = "attribute.example"
  masking = "MASK_ENTIRE_VALUE"
}