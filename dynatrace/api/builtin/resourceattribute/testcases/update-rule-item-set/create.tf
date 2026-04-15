resource "dynatrace_resource_attributes" "attributes" {
  keys {
    rule {
      enabled       = true
      attribute_key = "key1"
      masking       = "NOT_MASKED"
    }
    # to update
    rule {
      enabled       = true
      attribute_key = "key2"
      masking       = "NOT_MASKED"
    }
  }
}
