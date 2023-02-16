resource "dynatrace_http_monitor_performance" "#name#" {
  enabled = true
  scope   = "HTTP_CHECK-1234567890000000"
  thresholds {
    threshold {
      event     = "HTTP_CHECK-1234567890000000"
      threshold = 10
    }
  }
}