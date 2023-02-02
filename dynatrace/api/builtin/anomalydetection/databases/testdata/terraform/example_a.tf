resource "dynatrace_database_anomalies_v2" "#name#" {
  scope = "environment"
  database_connections {
    enabled             = true
    max_failed_connects = 5
    time_period         = 5
  }
  failure_rate {
    enabled        = true
    detection_mode = "fixed"
    fixed_detection {
      sensitivity = "low"
      threshold   = 0
      over_alerting_protection {
        minutes_abnormal_state = 1
        requests_per_minute    = 10
      }
    }
  }
  load_drops {
    enabled = false
  }
  load_spikes {
    enabled = false
  }
  response_time {
    enabled        = true
    detection_mode = "fixed"
    fixed_detection {
      sensitivity = "high"
      over_alerting_protection {
        minutes_abnormal_state = 1
        requests_per_minute    = 15
      }
      response_time_all {
        degradation_milliseconds = 15
      }
      response_time_slowest {
        slowest_degradation_milliseconds = 23
      }
    }
  }
}
