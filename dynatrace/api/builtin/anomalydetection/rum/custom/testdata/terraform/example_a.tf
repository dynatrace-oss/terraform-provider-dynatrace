resource "dynatrace_custom_app_anomalies" "#name#" {
  scope = "CUSTOM_APPLICATION-1234567890000000"
  error_rate_increase {
    enabled        = true
    detection_mode = "fixed"
    error_rate_increase_fixed {
      sensitivity        = "low"
      threshold_absolute = 5
    }
  }
  slow_user_actions {
    enabled = false
  }
  unexpected_high_load {
    enabled              = true
    threshold_percentage = 300
  }
  unexpected_low_load {
    enabled              = true
    threshold_percentage = 80
  }
}
