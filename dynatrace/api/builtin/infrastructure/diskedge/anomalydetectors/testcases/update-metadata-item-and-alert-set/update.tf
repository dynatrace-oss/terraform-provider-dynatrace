resource "dynatrace_disk_edge_anomaly_detectors" "detectors" {
  enabled           = true
  disk_name_filters = [ "$eq(/)" ]
  operating_system  = [ "WINDOWS", "LINUX" ]
  policy_name       = "#name#"
  scope             = "environment"
  alerts {
    alert {
      threshold_percent = 10
      trigger           = "AVAILABLE_DISK_SPACE_PERCENT_BELOW"
      sample_count_thresholds {
        dealerting_evaluation_window = 30
        dealerting_samples           = 24
        violating_evaluation_window  = 30
        violating_samples            = 18
      }
    }
    # update => re-create due to set-hash change
    alert {
      threshold_milliseconds = 1000
      trigger           = "WRITE_TIME_EXCEEDING"
      sample_count_thresholds {
        dealerting_evaluation_window = 30
        dealerting_samples           = 24
        violating_evaluation_window  = 30
        violating_samples            = 18
      }
    }
  }
  event_properties {
    event_propertie {
      metadata_key   = "ExampleKey"
      metadata_value = "ExampleValue"
    }
    # update => re-create due to set-hash change
    event_propertie {
      metadata_key   = "ExampleKeyEdit"
      metadata_value = "ExampleValueEdit"
    }
  }
}
