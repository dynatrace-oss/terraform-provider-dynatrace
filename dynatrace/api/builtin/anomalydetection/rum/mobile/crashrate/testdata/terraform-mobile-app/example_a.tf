data "dynatrace_mobile_application" "application" {
  name = "Application"
}

resource "dynatrace_mobile_app_crash_rate" "crash_rate" {
  application_id = data.dynatrace_mobile_application.application.id
  crash_rate_increase {
    enabled        = true
    detection_mode = "fixed"
    crash_rate_increase_fixed {
      absolute_crash_rate = 25
      concurrent_users    = 200
    }
  }
}
