resource "dynatrace_processgroup_naming_order" "process_group_naming_order" {
  naming_rule_ids = [
    dynatrace_processgroup_naming.first.id,
    dynatrace_processgroup_naming.second.id,
  ]  
}


resource "dynatrace_processgroup_naming" "first" {
  name    = "${randomize}"
  enabled = true
  format  = "{ProcessGroup:DetectedName} ${randomize} {ProcessGroup:CommandLineArgs}"
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

resource "dynatrace_processgroup_naming" "second" {
  name    = "${randomize}"
  enabled = true
  format  = "{ProcessGroup:DetectedName} ${randomize} {ProcessGroup:CommandLineArgs}"
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
