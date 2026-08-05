resource "dynatrace_synthetic_primary_grail_tags" "example" {
  scope = dynatrace_http_monitor.monitor.id
  tags {
    tag {
      key   = "my-key-1"
      value = "my-value-1"
    }
    // removed
  }
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

data "dynatrace_synthetic_location" "location" {
  name = "Location"
}
