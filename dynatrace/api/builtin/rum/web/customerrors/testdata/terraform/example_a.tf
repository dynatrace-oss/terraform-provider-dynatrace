resource "dynatrace_web_app_custom_errors" "#name#" {
  ignore_custom_errors_in_apdex_calculation = true
  scope                                     = "APPLICATION-1234567890000000"
  error_rules {
    error_rule {
      key_matcher   = "EQUALS"
      key_pattern   = "hashicorp"
      value_matcher = "ALL"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = false
      }
    }
    error_rule {
      key_matcher   = "CONTAINS"
      key_pattern   = "TF"
      value_matcher = "ENDS_WITH"
      value_pattern = "EX"
      capture_settings {
        capture = false
      }
    }
    error_rule {
      key_matcher   = "BEGINS_WITH"
      key_pattern   = "terraform"
      value_matcher = "CONTAINS"
      value_pattern = "example"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = true
      }
    }
    error_rule {
      key_matcher   = "ALL"
      value_matcher = "ALL"
      capture_settings {
        capture         = true
        consider_for_ai = false
        impact_apdex    = true
      }
    }
  }
}