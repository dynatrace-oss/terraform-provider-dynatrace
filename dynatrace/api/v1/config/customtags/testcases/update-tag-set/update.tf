resource "dynatrace_custom_tags" "tags" {
  entity_selector = "type(\"HOST\")"
  tags {
    filter {
      context = "CONTEXTLESS"
      key     = "KeyExampleA"
    }
    filter {
      context = "CONTEXTLESS"
      key     = "KeyExampleEdit"
      value   = "KeyExampleEdit"
    }
    filter {
      context = "CONTEXTLESS"
      key     = "KeyExampleNew"
      value   = "ValueExampleNew"
    }
  }
}
