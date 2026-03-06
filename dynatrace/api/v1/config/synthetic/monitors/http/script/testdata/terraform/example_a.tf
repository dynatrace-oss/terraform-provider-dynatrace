data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

resource "dynatrace_http_monitor" "monitor" {
  name      = "#name#"
  frequency = 1
  locations = [data.dynatrace_synthetic_location.location.id]
  anomaly_detection {
    loading_time_thresholds {
    }
    outage_handling {
      global_outage = true
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  no_script = true
}

resource "dynatrace_http_monitor_script" "script" {
  http_id = dynatrace_http_monitor.monitor.id
  script {
    request {
      description     = "request1"
      method          = "GET"
      url             = "https://example.com"
      configuration {
        accept_any_certificate = true
      }
    }
    request {
      description     = "request2"
      method          = "GET"
      url             = "https://example.com"
      configuration {
        accept_any_certificate = true
      }
    }
  }
}
