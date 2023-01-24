resource "dynatrace_browser_monitor" "#name#" {
  name      = "#name#"
  frequency = 15
  locations = ["GEOLOCATION-57F63BAD1C6A415C"]
  anomaly_detection {
    loading_time_thresholds {
      enabled = true
    }
    outage_handling {
      global_outage  = true
      retry_on_error = true
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  key_performance_metrics {
    load_action_kpm = "VISUALLY_COMPLETE"
    xhr_action_kpm  = "VISUALLY_COMPLETE"
  }
  script {
    type = "availability"
    configuration {
      device {
        name        = "Desktop"
        orientation = "landscape"
      }
    }
    events {
      event {
        description = "Loading of \"https://www.heise.de\""
        navigate {
          url = "https://www.heise.de"
          authentication {
            type  = "webform"
            creds = "CREDENTIALS_VAULT-6E12E52EC9718586"
          }
          wait {
            wait_for = "page_complete"
          }
        }
      }
    }
  }
}
