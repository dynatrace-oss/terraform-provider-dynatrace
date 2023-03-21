resource "dynatrace_web_app_request_errors" "#name#" {
  ignore_request_errors_in_apdex_calculation = false
  scope                                      = "APPLICATION-1234567890000000"
  error_rules {
    error_rule {
      consider_csp_violations = true
      consider_failed_images  = true
      error_codes             = "400"
      capture_settings {
        capture = false
      }
      filter_settings {
        filter = "ENDS_WITH"
        url    = "hashicorp"
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "404"
      capture_settings {
        capture         = true
        consider_for_ai = false
        impact_apdex    = false
      }
      filter_settings {
        filter = "CONTAINS"
        url    = "terraform"
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "404"
      capture_settings {
        capture = false
      }
      filter_settings {
        filter = "ENDS_WITH"
        url    = "favicon.ico"
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "404"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = true
      }
      filter_settings {
        filter = "ENDS_WITH"
        url    = ".js"
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "404"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = true
      }
      filter_settings {
        filter = "ENDS_WITH"
        url    = ".css"
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "400-499"
      capture_settings {
        capture         = true
        consider_for_ai = false
        impact_apdex    = true
      }
      filter_settings {
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "500-599"
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = true
      }
      filter_settings {
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = false
      error_codes             = "970-979"
      capture_settings {
        capture         = true
        consider_for_ai = false
        impact_apdex    = true
      }
      filter_settings {
      }
    }
    error_rule {
      consider_csp_violations = false
      consider_failed_images  = true
      capture_settings {
        capture         = true
        consider_for_ai = false
        impact_apdex    = false
      }
      filter_settings {
      }
    }
    error_rule {
      consider_csp_violations = true
      consider_failed_images  = false
      capture_settings {
        capture         = true
        consider_for_ai = true
        impact_apdex    = true
      }
      filter_settings {
      }
    }
  }
}
