resource "dynatrace_service_naming" "naming" {
  name    = "#name#"
  enabled = true
  format  = "format"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_DATABASE_TOPOLOGY"
      }
      database_topology {
        negate   = true
        operator = "EQUALS"
        value    = "CLUSTER"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TOPOLOGY"
      }
      service_topology {
        negate   = true
        operator = "EQUALS"
        value    = "FULLY_MONITORED "
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TYPE"
      }
      service_type {
        negate   = true
        operator = "EQUALS"
        value    = "CUSTOM_SERVICE"
      }
    }
  }
}
