resource "dynatrace_calculated_service_metric" "metric_i" {
  depends_on = [time_sleep.wait_for_request_attributes]
  name             = "#name#"
  enabled          = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "ESB_INPUT_TYPE"
      comparison {
        negate = false
        esb_input_node_type {
          operator = "EQUALS_ANY_OF"
          values   = ["EVENT_INPUT", "MQ_INPUT_NODE"]
        }
      }
    }
  }
  dimension_definition {
    name              = "jhj"
    dimension         = "{ESB:LibraryName}"
    top_x             = 40
    top_x_aggregation = "SUM"
    top_x_direction   = "DESCENDING"
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = dynatrace_request_attribute.attribute.name
  }
}
