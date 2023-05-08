resource "dynatrace_web_app_key_performance_custom" "#name#" {
  scope = "APPLICATION_METHOD-1234567890000000"
  thresholds {
    frustrating_threshold_seconds = 12
    tolerated_threshold_seconds   = 3
  }
}
