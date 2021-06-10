resource "dynatrace_calculated_service_metric" "#name#" {
  name             = "#name#"
  enabled          = true
  management_zones = ["AAAA"]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "IS_KEY_REQUEST"
      comparison {
        negate = false
        boolean {
          operator = "EQUALS"
          value    = true
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = "foo"
  }
}
