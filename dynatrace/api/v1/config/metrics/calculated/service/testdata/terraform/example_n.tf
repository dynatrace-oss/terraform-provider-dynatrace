resource "dynatrace_calculated_service_metric" "metric_n" {
  name = "#name#"
  enabled = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
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
