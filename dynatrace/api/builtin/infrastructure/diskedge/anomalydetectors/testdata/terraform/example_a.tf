resource "dynatrace_disk_edge_anomaly_detectors" "#name#" {
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
  }
  event_properties {
    event_propertie {
      metadata_key   = "ExampleKey"
      metadata_value = "ExampleValue"
    }
  }

  detection_conditions {
    detection_condition {
      rule_type = "RuleTypeHost"
      host_metadata_condition {
        host_metadata_condition {
          # key_must_exist   = true
          metadata_condition = "$contains(terraform)"
          metadata_key       = "ExampleKey"
        }
      }
    }
    detection_condition {
      local_disk_condition = "REMOTE"
      property             = "DiskType"
      rule_type            = "RuleTypeDisk"
    }
    detection_condition {
      disk_filesystem_condition = "$match(ext*)"
      property                  = "DiskFilesystem"
      rule_type                 = "RuleTypeDisk"
    }
    detection_condition {
      property  = "DiskTotalSpace"
      rule_type = "RuleTypeDisk"
      disk_total_condition {
        threshold_above = 10
        threshold_below = 1000
      }
    }
  }
}
