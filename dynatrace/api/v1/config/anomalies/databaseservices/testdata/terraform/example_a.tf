resource "dynatrace_database_anomalies" "#name#" {
  db_connect_failures {
    connection_fails_count = 5
    eval_period            = 5
  }
  failure_rate {
    thresholds {
      sensitivity = "LOW"
      threshold   = 0
    }
  }
  response_time {
    thresholds {
      load                 = "FIFTEEN_REQUESTS_PER_MINUTE"
      milliseconds         = 15
      sensitivity          = "HIGH"
      slowest_milliseconds = 23
    }
  }
}
