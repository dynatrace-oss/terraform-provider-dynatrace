resource "dynatrace_calculated_service_metric" "metric_e" {
  depends_on = [time_sleep.wait_for_request_attributes]
  name             = "#name#"
  enabled          = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
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
    request_attribute = dynatrace_request_attribute.attribute.name
  }
}
