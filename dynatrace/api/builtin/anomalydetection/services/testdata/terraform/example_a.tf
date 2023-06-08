resource "dynatrace_service_anomalies_v2" "#name#" {
  scope = "SERVICE-1234567890000000"
  failure_rate {
    enabled        = true
    detection_mode = "fixed"
    fixed_detection {
      sensitivity = "high"
      threshold   = 5
      over_alerting_protection {
        minutes_abnormal_state = 1
        requests_per_minute    = 10
      }
    }
  }
  load_drops {
    enabled                = true
    load_drop_percent      = 50
    minutes_abnormal_state = 1
  }
  load_spikes {
    enabled                = true
    load_spike_percent     = 200
    minutes_abnormal_state = 1
  }
  response_time {
    enabled        = true
    detection_mode = "fixed"
    fixed_detection {
      sensitivity = "high"
      over_alerting_protection {
        minutes_abnormal_state = 1
        requests_per_minute    = 10
      }
      response_time_all {
        degradation_milliseconds = 100
      }
      response_time_slowest {
        slowest_degradation_milliseconds = 1000
      }
    }
  }
}
