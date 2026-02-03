data "dynatrace_mobile_application" "application" {
  name = "Application"
}

resource "dynatrace_mobile_app_request_errors" "request_errors" {
  scope = data.dynatrace_mobile_application.application.id
  error_rules {
    error_rule {
      error_codes = "409"
    }
    error_rule {
      error_codes = "410"
    }
  }
}
