resource "dynatrace_web_app_anomalies" "#name#" {
  scope = "APPLICATION_METHOD-1234567890000000"
  error_rate {
    enabled                   = true
    error_rate_detection_mode = "auto"
    error_rate_auto {
      absolute_increase = 5
      relative_increase = 55
      over_alerting_protection {
        actions_per_minute     = 15
        minutes_abnormal_state = 5
      }
    }
  }
  response_time {
    enabled        = true
    detection_mode = "fixed"
    response_time_fixed {
      sensitivity = "high"
      over_alerting_protection {
        actions_per_minute     = 15
        minutes_abnormal_state = 5
      }
      response_time_all {
        degradation_milliseconds = 100
      }
      response_time_slowest {
        slowest_degradation_milliseconds = 1500
      }
    }
  }
  traffic_drops {
    enabled = true
    traffic_drops {
      abnormal_state_abnormal_state = 1
      traffic_drop_percentage       = 55
    }
  }
  traffic_spikes {
    enabled = true
    traffic_spikes {
      minutes_abnormal_state   = 5
      traffic_spike_percentage = 250
    }
  }
}
