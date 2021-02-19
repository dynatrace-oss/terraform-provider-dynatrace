resource "dynatrace_autotag" "#name#" {
  rules {
    conditions {
      service_topology_comparison {
        negate   = false
        type     = "SERVICE_TOPOLOGY"
        operator = "EQUALS"
        value    = "EXTERNAL_SERVICE"
      }
      base_condition_key {
        attribute = "SERVICE_TOPOLOGY"
      }
    }
    conditions {
      string_comparison {
        negate         = false
        type           = "STRING"
        operator       = "EQUALS"
        value          = "Requests to public networks"
        case_sensitive = true
      }
      base_condition_key {
        attribute = "SERVICE_DETECTED_NAME"
      }
    }
    enabled      = true
    type         = "SERVICE"
    value_format = "{Service:EndpointPath}"
  }
  metadata {
    cluster_version        = "1.211.90.20210216-141629"
    configuration_versions = [7]
  }
  name = "#name#"
}
