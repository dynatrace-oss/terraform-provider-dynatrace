resource "dynatrace_calculated_service_metric" "#name#" {
  name             = "#name#"
  enabled          = true
  management_zones = ["AAAA"]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "HTTP_REQUEST_METHOD"
      comparison {
        negate = false
        http_method {
          operator = "EQUALS_ANY_OF"
          values   = ["POST", "GET"]
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = "foo"
  }
}
