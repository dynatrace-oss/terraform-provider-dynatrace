resource "dynatrace_service_anomalies" "#name#" {
  failure_rates {
    auto {
      absolute = 0
      relative = 50
    }
  }
  load {
    drops {
      minutes = 1
      percent = 50
    }
    spikes {
      minutes = 1
      percent = 200
    }
  }
  load_drops {
    minutes = 1
    percent = 50
  }
  response_times {
    auto {
      load                 = "TEN_REQUESTS_PER_MINUTE"
      milliseconds         = 100
      percent              = 50
      slowest_milliseconds = 1000
      slowest_percent      = 100
    }
  }
}
