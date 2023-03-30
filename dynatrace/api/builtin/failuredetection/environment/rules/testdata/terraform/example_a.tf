resource "dynatrace_failure_detection_rules" "#name#" {
  name         = ""
  enabled      = true
  parameter_id = "vu9U3hXa3q0AAAABADBidWlsdGluOmZhaWx1cmUtZGV0ZWN0aW9uLmVudmlyb25tZW50LnBhcmFtZXRlcnMABnRlbmFudAAGdGVuYW50ACQwYTIwM2UxYy1mYjYxLTMyMDEtOTc1NC00M2NiNDdhMWI1ODG-71TeFdrerQ"
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