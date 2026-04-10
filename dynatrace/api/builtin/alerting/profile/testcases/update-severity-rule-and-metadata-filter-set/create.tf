resource "dynatrace_alerting" "alerting" {
  name            = "#name#"
  management_zone = ""
  filters {
    filter {
      custom {
        metadata {
          items {
            filter {
              key   = "key"
              value = "value"
            }
            filter {
              key   = "key2"
              value = "value2"
            }
          }
        }
      }
    }
  }
  rules {
    rule {
      include_mode     = "INCLUDE_ALL"
      tags             = ["EnvironmentA:production", "Team:test"]
      delay_in_minutes = 0
      severity_level   = "AVAILABILITY"
    }
    rule {
      include_mode     = "INCLUDE_ALL"
      tags             = ["EnvironmentB:production", "Team:test"]
      delay_in_minutes = 0
      severity_level   = "CUSTOM_ALERT"
    }
  }
}
