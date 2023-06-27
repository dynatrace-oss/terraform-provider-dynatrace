resource "dynatrace_mobile_app_crash_rate" "#name#" {
  crash_rate_increase {
    enabled        = true
    detection_mode = "fixed"
    crash_rate_increase_fixed {
      absolute_crash_rate = 45
      concurrent_users    = 300
    }
  }
}
