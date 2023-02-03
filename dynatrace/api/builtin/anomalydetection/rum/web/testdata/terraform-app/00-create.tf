resource "dynatrace_web_app_anomalies" "#name#" {
  scope = "APPLICATION-1234567890000000"
  error_rate {
    enabled                   = true
    error_rate_detection_mode = "auto"
    error_rate_auto {
      absolute_increase = 10
      relative_increase = 70
      over_alerting_protection {
        actions_per_minute     = 12
        minutes_abnormal_state = 2
      }
    }
  }
  response_time {
    enabled        = true
    detection_mode = "fixed"
    response_time_fixed {
      sensitivity = "low"
      over_alerting_protection {
        actions_per_minute     = 10
        minutes_abnormal_state = 1
      }
      response_time_all {
        degradation_milliseconds = 100
      }
      response_time_slowest {
        slowest_degradation_milliseconds = 1000
      }
    }
  }
  traffic_drops {
    enabled = false
  }
  traffic_spikes {
    enabled = true
    traffic_spikes {
      minutes_abnormal_state   = 2
      traffic_spike_percentage = 250
    }
  }
}
