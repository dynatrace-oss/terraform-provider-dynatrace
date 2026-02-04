resource "dynatrace_calculated_service_metric" "metric_l" {
  depends_on = [time_sleep.wait_for_request_attributes]
  name             = "#name#"
  enabled          = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "SERVICE_REQUEST_ATTRIBUTE"
      comparison {
        negate = false
        string_request_attribute {
          case_sensitive       = true
          match_on_child_calls = false
          operator             = "BEGINS_WITH"
          request_attribute    = dynatrace_request_attribute.behavior_class.name
          value                = "\"kk"
        }
      }
    }
  }
  dimension_definition {
    name              = "jhj"
    dimension         = "{ESB:LibraryName}{acceptranges}"
    top_x             = 40
    top_x_aggregation = "SUM"
    top_x_direction   = "DESCENDING"
    placeholders {
      placeholder {
        name                 = "acceptranges"
        aggregation          = "FIRST"
        attribute            = "SERVICE_REQUEST_ATTRIBUTE"
        delimiter_or_regex   = "k"
        end_delimiter        = "l"
        kind                 = "BETWEEN_DELIMITER"
        normalization        = "TO_UPPER_CASE"
        request_attribute    = dynatrace_request_attribute.accept_ranges.name
        use_from_child_calls = true
        source {
          management_zone = dynatrace_management_zone_v2.mzone.name
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = dynatrace_request_attribute.attribute.name
  }
}
