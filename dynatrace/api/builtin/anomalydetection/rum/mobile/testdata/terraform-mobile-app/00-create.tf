resource "dynatrace_mobile_app_anomalies" "#name#" {
  scope = "MOBILE_APPLICATION-1234567890000000"
  error_rate_increase {
    enabled        = true
    detection_mode = "fixed"
    error_rate_increase_fixed {
      sensitivity        = "medium"
      threshold_absolute = 6
    }
  }
  slow_user_actions {
    enabled        = true
    detection_mode = "fixed"
    slow_user_actions_fixed {
      sensitivity = "high"
      duration_avoid_overalerting {
        min_action_rate = 12
      }
      duration_threshold_all_fixed {
        duration_threshold = 200
      }
      duration_threshold_slowest {
        duration_threshold = 900
      }
    }
  }
  unexpected_high_load {
    enabled              = true
    threshold_percentage = 300
  }
  unexpected_low_load {
    enabled = false
  }
}
