resource "dynatrace_calculated_service_metric" "metric_k" {
  depends_on = [time_sleep.wait_for_request_attributes]
  name             = "#name#"
  enabled          = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "PROCESS_GROUP_TAG"
      comparison {
        negate = false
        tag {
          operator = "TAG_KEY_EQUALS_ANY_OF"
          values {
            value {
              context = "KUBERNETES"
              key     = "dynatrace"
            }
            value {
              context = "KUBERNETES"
              key     = "name"
            }
          }
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
