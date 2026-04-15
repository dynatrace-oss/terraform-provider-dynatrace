resource "dynatrace_http_monitor_performance" "performance" {
  enabled = true
  scope   = "HTTP_CHECK-0000000000000000"
  thresholds {
    threshold {
      event     = "HTTP_CHECK-0000000000000000"
      threshold = 10
    }
    # updated
    threshold {
      event     = "HTTP_CHECK_STEP-0000000000000001"
      threshold = 20
    }
  }
}
