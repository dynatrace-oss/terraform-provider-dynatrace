resource "dynatrace_application_anomalies" "#name#" {
  failure_rate {
    auto {
      absolute = 5
      relative = 50
    }
  }
  response_time {
    auto {
      load                 = "TEN_REQUESTS_PER_MINUTE"
      milliseconds         = 100
      percent              = 50
      slowest_milliseconds = 1000
      slowest_percent      = 100
    }
  }
  traffic {
    drops {
      enabled = true
      percent = 50
    }
  }
}
