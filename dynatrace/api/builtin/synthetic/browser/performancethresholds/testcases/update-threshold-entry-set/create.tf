resource "dynatrace_browser_monitor_performance" "performance" {
  enabled = true
  scope   = "SYNTHETIC_TEST-0000000000000000"
  thresholds {
    threshold {
      event     = "SYNTHETIC_TEST-0000000000000000"
      threshold = 10
    }
    # to update
    threshold {
      event     = "SYNTHETIC_TEST_STEP-0000000000000000"
      threshold = 10
    }
  }
}
