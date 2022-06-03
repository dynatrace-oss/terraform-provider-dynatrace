resource "dynatrace_calculated_service_metric" "#name#" {
  name             = "#name#"
  enabled          = true
  management_zones = ["AAAA"]
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
        request_attribute    = "Accept-Ranges"
        use_from_child_calls = true
        source {
          management_zone = "AAAA"
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = "foo"
  }
}
