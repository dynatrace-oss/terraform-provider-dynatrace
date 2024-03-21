resource "dynatrace_failure_detection_rules" "first-instance" {
  name         = "#name#"
  enabled      = true
  parameter_id = "00000000-0000-0000-0000-000000000000"
  conditions {
    condition {
      attribute = "SERVICE_NAME"
      predicate {
        case_sensitive = true
        predicate_type = "STRING_EQUALS"
        text_values    = ["Terraform"]
      }
    }
  }
}

resource "dynatrace_failure_detection_rules" "second-instance" {
  name         = "#name#-second"
  enabled      = true
  parameter_id = "00000000-0000-0000-0000-000000000000"
  conditions {
    condition {
      attribute = "SERVICE_NAME"
      predicate {
        case_sensitive = true
        predicate_type = "STRING_EQUALS"
        text_values    = ["Terraform-second"]
      }
    }
  }
  insert_after = dynatrace_failure_detection_rules.first-instance.id
}
