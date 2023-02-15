resource "dynatrace_browser_monitor_performance" "#name#" {
  enabled = true
  scope   = "SYNTHETIC_TEST-1234567890000000"
  thresholds {
    threshold {
      event     = "SYNTHETIC_TEST-1234567890000000"
      threshold = 10
    }
  }
}