resource "dynatrace_mobile_app_crash_rate" "#name#" {
  application_id = "MOBILE_APPLICATION-1234567890000000"
  crash_rate_increase {
    enabled        = true
    detection_mode = "fixed"
    crash_rate_increase_fixed {
      absolute_crash_rate = 25
      concurrent_users    = 200
    }
  }
}
