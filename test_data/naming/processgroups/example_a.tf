resource "dynatrace_processgroup_naming" "#name#" {
  name    = "#name#"
  enabled = true
  format  = "{ProcessGroup:DetectedName} {ProcessGroup:CommandLineArgs}"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "PROCESS_GROUP_TECHNOLOGY"
      }
      tech {
        negate   = false
        operator = "EQUALS"
        value {
          type = "ADO_NET"
        }
      }
    }
    condition {
      process_metadata {
        attribute   = "PROCESS_GROUP_PREDEFINED_METADATA"
        dynamic_key = "COMMAND_LINE_ARGS"
      }
      string {
        case_sensitive = true
        negate         = false
        operator       = "CONTAINS"
        value          = "-config"
      }
    }
  }
}
