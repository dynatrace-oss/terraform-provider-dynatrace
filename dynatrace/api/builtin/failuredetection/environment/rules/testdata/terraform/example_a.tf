resource "dynatrace_failure_detection_rules" "#name#" {
  name         = ""
  enabled      = true
  parameter_id = "00000000-0000-0000-0000-000000000000"
  conditions {
    condition {
      attribute = "SERVICE_NAME"
      predicate {
        case_sensitive = true
        predicate_type = "STRING_EQUALS"
        text_values    = [ "Terraform" ]
      }
    }
  }
}