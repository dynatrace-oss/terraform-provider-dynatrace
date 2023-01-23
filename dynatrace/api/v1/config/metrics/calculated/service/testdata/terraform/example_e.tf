resource "dynatrace_calculated_service_metric" "#name#" {
  name             = "#name#"
  enabled          = true
  management_zones = ["AAAA"]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "IO_TIME"
      comparison {
        negate = false
        number {
          operator = "GREATER_THAN"
          value    = 4
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = "foo"
  }
}
