resource "dynatrace_calculated_service_metric" "#name#" {
  name = "#name#"
  enabled = true
  management_zones = ["AAAA"]
  metric_key = "calc:service.#name#"
  unit = "MICRO_SECOND"
  conditions {
    condition {
      attribute = "CPU_TIME"
      comparison {
        negate = false
        number {
          operator = "EQUALS_ANY_OF"
          values = [1,2]
        }
      }
    }
  }
  metric_definition {
    metric = "RESPONSE_TIME"
  }
}
