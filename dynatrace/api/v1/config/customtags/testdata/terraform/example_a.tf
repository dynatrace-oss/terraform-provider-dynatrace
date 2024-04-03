resource "dynatrace_custom_tags" "#name#" {
    entity_selector = "entityId(\"HOST-123456789000000\")"
    tags {
      filter {
        context = "CONTEXTLESS"
        key = "KeyExampleA"
      }
      filter {
        context = "CONTEXTLESS"
        key = "KeyExampleA"
        value = "ValueExample1"
      }
      filter {
        context = "CONTEXTLESS"
        key = "KeyExampleB"
      }
      filter {
        context = "CONTEXTLESS"
        key = "KeyExampleC"
        value = "ValueExample2"
      }
    }
}