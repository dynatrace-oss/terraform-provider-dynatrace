resource "dynatrace_request_naming" "terraform-request-naming-global" {
  enabled        = true
  naming_pattern = "terraform-request-naming-global"
  conditions {
    condition {
      attribute = "ONE_AGENT_ATTRIBUTE"
      comparison {
        string_one_agent_attribute {
          case_sensitive = false
          operator = "CONTAINS"
          one_agent_attribute_key = "http.route"
          value = "/services/*"
        }
      }
    }
  }
}
