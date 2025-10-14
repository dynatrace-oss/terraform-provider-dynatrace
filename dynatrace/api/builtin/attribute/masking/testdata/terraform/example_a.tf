resource "dynatrace_attribute_masking" "#name#" {
  enabled = true
  key = "attribute.#name#"
  masking = "MASK_ENTIRE_VALUE"
}
