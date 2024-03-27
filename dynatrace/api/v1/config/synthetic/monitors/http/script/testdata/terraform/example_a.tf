resource "dynatrace_http_monitor_script" "#name#" {
  http_id = "${dynatrace_http_monitor.monitor.id}"
  script {
    request {
      description     = "request1"
      method          = "GET"
      url             = "http://httpstat.us/200"
      configuration {
        accept_any_certificate = true
      }
    }
    request {
      description     = "request2"
      method          = "GET"
      url             = "http://httpstat.us/400"
      configuration {
        accept_any_certificate = true
      }
    }
  }
}

resource "dynatrace_http_monitor" "monitor" {
  name      = "#name#"
  frequency = 1
  locations = ["GEOLOCATION-F3E06A526BE3B4C4"]
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
