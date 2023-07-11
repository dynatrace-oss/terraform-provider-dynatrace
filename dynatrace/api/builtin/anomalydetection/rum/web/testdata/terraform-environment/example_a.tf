resource "dynatrace_web_app_anomalies" "#name#" {
  scope = "environment"
  error_rate {
    enabled                   = true
    error_rate_detection_mode = "auto"
    error_rate_auto {
      absolute_increase = 5
      relative_increase = 50
      over_alerting_protection {
        actions_per_minute     = 10
        minutes_abnormal_state = 1
      }
    }
  }
  response_time {
    enabled        = true
    detection_mode = "auto"
    response_time_auto {
      over_alerting_protection {
        actions_per_minute     = 10
        minutes_abnormal_state = 1
      }
      response_time_all {
        degradation_milliseconds = 1000
        degradation_percent      = 100
      }
      response_time_slowest {
        slowest_degradation_milliseconds = 2000
        slowest_degradation_percent      = 10
      }
    }
  }
  traffic_drops {
    enabled = true
    traffic_drops {
      abnormal_state_abnormal_state = 1
      traffic_drop_percentage       = 75
    }
  }
  traffic_spikes {
    enabled = false
  }
}
