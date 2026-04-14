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
    alert {
      threshold_milliseconds = 1000
      trigger           = "READ_TIME_EXCEEDING"
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
    event_propertie {
      metadata_key   = "ExampleKey2"
      metadata_value = "ExampleValue2"
    }
  }
}
