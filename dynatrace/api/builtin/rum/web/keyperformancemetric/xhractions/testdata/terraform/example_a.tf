resource "dynatrace_web_app_key_performance_xhr" "#name#" {
  kpm   = "VISUALLY_COMPLETE"
  scope = "APPLICATION_METHOD-1234567890000000"
  fallback_thresholds {
    frustrating_fallback_threshold_seconds = 12
    tolerated_fallback_threshold_seconds   = 3
  }
  thresholds {
    frustrating_threshold_seconds = 12
    tolerated_threshold_seconds   = 3
  }
}
