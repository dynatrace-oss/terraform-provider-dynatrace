resource "dynatrace_custom_app_crash_rate" "CUSTOM_APPLICATION-1234567890000000" {
  scope = "CUSTOM_APPLICATION-1234567890000000"
  crash_rate_increase {
    enabled        = true
    detection_mode = "fixed"
    crash_rate_increase_fixed {
      absolute_crash_rate = 25
      concurrent_users    = 200
    }
  }
}
