resource "dynatrace_mobile_app_request_errors" "#name#" {
  scope = "MOBILE_APPLICATION-1234567890000000"
  error_rules {
    error_rule {
      error_codes = "409"
    }
    error_rule {
      error_codes = "410"
    }
  }
}
