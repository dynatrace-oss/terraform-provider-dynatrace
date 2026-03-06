resource "dynatrace_calculated_service_metric" "metric_m" {
  depends_on = [time_sleep.wait_for_request_attributes]
  name             = "#name#"
  enabled          = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "SERVICE_DISPLAY_NAME"
      comparison {
        negate = false
        fast_string {
          case_sensitive = true
          operator       = "CONTAINS"
          value          = "\"kk"
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
  ignore_muted_requests = true
}
