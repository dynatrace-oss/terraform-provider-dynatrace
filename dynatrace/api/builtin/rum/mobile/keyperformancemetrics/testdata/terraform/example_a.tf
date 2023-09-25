resource "dynatrace_mobile_app_key_performance" "#name#" {
  frustrating_if_reported_or_web_request_error = true
  scope                                        = "DEVICE_APPLICATION_METHOD-1234567890000000"
  thresholds {
    frustrating_threshold_seconds = 12
    tolerable_threshold_seconds   = 3
  }
}