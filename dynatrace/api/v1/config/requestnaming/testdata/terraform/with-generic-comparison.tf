resource "dynatrace_request_naming" "terraform-request-naming-global" {
  enabled        = true
  naming_pattern = "terraform-request-naming-global"
  conditions {
    condition {
      attribute = "ONE_AGENT_ATTRIBUTE"
      comparison {
        generic {
          type     = "STRING_ONE_AGENT_ATTRIBUTE"
          unknowns = jsonencode({
            "caseSensitive": false,
            "comparison": "CONTAINS",
            "oneAgentAttributeKey": "http.route",
            "value": "/services/*",
            "values": null
          })
        }
      }
    }
  }
}
