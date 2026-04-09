resource "dynatrace_custom_tags" "tags" {
  entity_selector = "type(\"HOST\")"
  tags {
    filter {
      context = "CONTEXTLESS"
      key     = "KeyExampleA"
    }
    filter {
      context = "CONTEXTLESS"
      key     = "KeyExampleA1"
      value   = "ValueExample1"
    }
  }
}
