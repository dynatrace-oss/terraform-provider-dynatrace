resource "dynatrace_processgroup_naming" "team-hawaiian-pizza" {
  name    = "team-hawaiian-pizza"
  enabled = false
  format  = "{ProcessGroup:DetectedName} {ProcessGroup:CommandLineArgs}"
  conditions {
    condition {
      process_metadata {
        attribute   = "PROCESS_GROUP_PREDEFINED_METADATA"
        dynamic_key = "COMMAND_LINE_ARGS"
      }
      string {
        operator = "CONTAINS"
        value    = "team-hawaiian-pizza"
      }
    }
  }
}