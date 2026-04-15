resource "dynatrace_resource_attributes" "attributes" {
  keys {
    rule {
      enabled       = true
      attribute_key = "key1"
      masking       = "NOT_MASKED"
    }
    # updated
    rule {
      enabled       = true
      attribute_key = "keyEdit"
      masking       = "NOT_MASKED"
    }
  }
}
