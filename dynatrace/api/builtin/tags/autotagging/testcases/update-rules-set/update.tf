resource "dynatrace_autotag_v2" "tag" {
  name = "#name#"
  rules {
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "{ProcessGroup:Environment:keptn_stage}"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type                 = "SERVICE"
        service_to_host_propagation = false
        service_to_pgpropagation    = true
        conditions {
          condition {
            dynamic_key        = "keptn_stage"
            dynamic_key_source = "ENVIRONMENT"
            key                = "PROCESS_GROUP_CUSTOM_METADATA"
            operator           = "EXISTS"
          }
        }
      }
    }
    # contains updated conditions
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "sprint"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type = "SYNTHETIC_TEST"
        conditions {
          condition {
            key      = "BROWSER_MONITOR_TAGS"
            operator = "TAG_KEY_EQUALS"
            tag      = "sprint"
          }
          condition {
            key      = "BROWSER_MONITOR_TAGS"
            operator = "TAG_KEY_EQUALS"
            tag      = "hardening"
          }
        }
      }
    }
    # prod removed, new one added
    rule {
      type                = "ME"
      enabled             = true
      value_format        = "newOne"
      value_normalization = "Leave text as-is"
      attribute_rule {
        entity_type = "SYNTHETIC_TEST"
        conditions {
          condition {
            key      = "BROWSER_MONITOR_TAGS"
            operator = "TAG_KEY_EQUALS"
            tag      = "newOne"
          }
        }
      }
    }
  }
}
