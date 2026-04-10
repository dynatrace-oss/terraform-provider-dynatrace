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
              key   = "keyEdit"
              value = "valueEdit"
            }
            filter {
              key   = "keyNew"
              value = "valueNew"
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
      tags             = ["EnvironmentF:production", "Team:test"]
      delay_in_minutes = 0
      severity_level   = "RESOURCE_CONTENTION"
    }
    rule {
      include_mode     = "INCLUDE_ALL"
      tags             = ["EnvironmentE:production", "Team:test"]
      delay_in_minutes = 0
      severity_level   = "PERFORMANCE"
    }
  }
}
