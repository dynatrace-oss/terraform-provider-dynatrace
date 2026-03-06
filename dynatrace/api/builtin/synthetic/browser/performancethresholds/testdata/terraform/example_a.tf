data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

resource "dynatrace_browser_monitor" "monitor" {
  name                   = "#name#"
  frequency              = 15
  locations              = [data.dynatrace_synthetic_location.location.id]
  key_performance_metrics {
    load_action_kpm = "VISUALLY_COMPLETE"
    xhr_action_kpm  = "VISUALLY_COMPLETE"
  }
  anomaly_detection {
    loading_time_thresholds {
      enabled = true
      thresholds {
        threshold {
          event_index   = 0
          request_index = 0
          type          = "TOTAL"
          value_ms      = 10000
        }
      }
    }
    outage_handling {
      global_outage = true
      local_outage = false
      retry_on_error = true
      global_outage_policy {
        consecutive_runs = 1
      }
    }
  }
  script {
    type = "clickpath"
    configuration {
      bypass_csp = true
      user_agent = "Mozilla"
      device {
        name        = "Desktop"
        orientation = "landscape"
      }
    }
    events {
      event {
        description = "my description"
        navigate {
          url = "https://www.example.com"
        }
      }
    }
  }
}

resource "dynatrace_browser_monitor_performance" "performance" {
  enabled = true
  scope   = dynatrace_browser_monitor.monitor.id
  thresholds {
    threshold {
      event     = dynatrace_browser_monitor.monitor.id
      threshold = 10
    }
  }
}
