resource "dynatrace_aws_anomalies" "#name#" {
  ec_2_candidate_high_cpu_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      cpu_usage = 94
    }
  }
  elb_high_connection_errors_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      connection_errors_per_minute = 9
    }
  }
  lambda_high_error_rate_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      failed_invocations_rate = 4
    }
  }
  rds_high_cpu_detection {
    enabled = false
  }
  rds_high_memory_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      free_memory = 94
      swap_usage  = 4
    }
  }
  rds_high_write_read_latency_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      read_write_latency = 194
    }
  }
  rds_low_storage_detection {
    enabled = false
  }
  rds_restarts_sequence_detection {
    enabled = false
  }
}
