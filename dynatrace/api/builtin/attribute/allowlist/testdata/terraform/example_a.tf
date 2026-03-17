resource "dynatrace_attribute_allow_list" "list" {
  enabled = true
  key = "attribute.#name#"
}
