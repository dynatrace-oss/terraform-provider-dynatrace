resource "dynatrace_service_naming" "#name#" {
  name    = "#name#"
  enabled = true
  format  = "ABCD"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TYPE"
      }
      service_type {
        negate   = false
        operator = "EQUALS"
        value    = "WEB_REQUEST_SERVICE"
      }
    }
    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TECHNOLOGY"
      }
      tech {
        negate   = false
        operator = "EQUALS"
        value {
          type = "APACHE_HTTP_SERVER"
        }
      }
    }
    condition {
      key {
        type      = "STATIC"
        attribute = "SERVICE_TOPOLOGY"
      }
      service_topology {
        negate   = false
        operator = "EQUALS"
        value    = "FULLY_MONITORED"
      }
    }
    condition {
      key {
        type      = "STATIC"
        attribute = "PROCESS_GROUP_TAGS"
      }
      tag {
        negate   = true
        operator = "TAG_KEY_EQUALS"
        value {
          context = "CONTEXTLESS"
          key     = "dfoo"
        }
      }
    }
  }
}
