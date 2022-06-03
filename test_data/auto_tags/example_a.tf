resource "dynatrace_autotag" "#name#" {
  rules {
    conditions {
      string {
        negate         = false
        operator       = "EQUALS"
        value          = "Requests to public networks"
        case_sensitive = true
      }
      key {
        attribute = "SERVICE_DETECTED_NAME"
        type      = "STATIC"
      }
    }    
    conditions {
      service_topology {
        negate   = false
        operator = "EQUALS"
        value    = "EXTERNAL_SERVICE"
      }
      key {
        attribute = "SERVICE_TOPOLOGY"
        type      = "STATIC"
      }
    }
    enabled      = true
    type         = "SERVICE"
    value_format = "{Service:EndpointPath}"
    normalization     = "LEAVE_TEXT_AS_IS"
  }
  name = "#name#"
}
