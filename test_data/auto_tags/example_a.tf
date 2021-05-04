resource "dynatrace_autotag" "#name#" {
  rules {
    conditions {
      service_topology {
        negate   = false
        operator = "EQUALS"
        value    = "EXTERNAL_SERVICE"
      }
      key {
        attribute = "SERVICE_TOPOLOGY"
      }
    }
    conditions {
      string {
        negate         = false
        operator       = "EQUALS"
        value          = "Requests to public networks"
        case_sensitive = true
      }
      key {
        attribute = "SERVICE_DETECTED_NAME"
      }
    }
    enabled      = true
    type         = "SERVICE"
    value_format = "{Service:EndpointPath}"
  }
  name = "#name#"
}
