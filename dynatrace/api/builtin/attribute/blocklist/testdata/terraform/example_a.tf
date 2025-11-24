resource "dynatrace_attribute_block_list" "#name#" {
  enabled = true
  key = "attribute.#name#"
}
