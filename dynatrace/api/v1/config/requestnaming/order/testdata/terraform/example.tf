resource "dynatrace_request_namings" "order" {
  ids = [
    dynatrace_request_naming.terraform-request-naming-global[0].id,
    dynatrace_request_naming.terraform-request-naming-global[1].id
  ]
}

resource "dynatrace_request_naming" "terraform-request-naming-global" {
  enabled        = true
  count          = 2
  naming_pattern = "terraform-request-naming-global-#name#-${count.index}"
  conditions {
    condition {
      attribute = "ONE_AGENT_ATTRIBUTE"
      comparison {
        string_one_agent_attribute {
          case_sensitive          = false
          operator                = "CONTAINS"
          one_agent_attribute_key = "http.route"
          value                   = "/services/*"
        }
      }
    }
  }
}
