resource "dynatrace_calculated_service_metric" "#name#" {
  name             = "#name#"
  enabled          = true
  management_zones = ["AAAA"]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "SERVICE_TYPE"
      comparison {
        negate = false
        service_type {
          operator = "EQUALS_ANY_OF"
          values   = ["DATABASE_SERVICE", "CUSTOM_SERVICE"]
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = "foo"
  }
}
