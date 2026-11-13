resource "dynatrace_autotag_v2" "network_availability" {
  name = "#name#"
  rules {
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "true"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type                 = "MULTIPROTOCOL_MONITOR"
        conditions {
          condition {
            case_sensitive = false
            key            = "NETWORK_AVAILABILITY_MONITOR_NAME"
            operator       = "REGEX_MATCHES"
            string_value   = ".*example.*"
          }
        }
      }
    }
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "prod"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type = "SYNTHETIC_TEST"
        conditions {
          condition {
            key      = "BROWSER_MONITOR_TAGS"
            operator = "TAG_KEY_EQUALS"
            tag      = "prod"
          }
        }
      }
    }
  }
}
