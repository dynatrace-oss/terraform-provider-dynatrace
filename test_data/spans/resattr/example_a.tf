resource "dynatrace_resource_attributes" "dynatrace_resource_attributes" {
  keys {
    rule {
      enabled       = true
      attribute_key = "gffgf"
      masking       = "NOT_MASKED"
    }
    rule {
      enabled       = true
      attribute_key = "cdefgh"
      masking       = "NOT_MASKED"
    }
    rule {
      enabled       = true
      attribute_key = "jjhhj"
      masking       = "NOT_MASKED"
    }
  }
}