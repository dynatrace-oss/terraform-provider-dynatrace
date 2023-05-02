resource "dynatrace_vmware_anomalies" "#name#" {
  dropped_packets_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      dropped_packets_per_second = 4
    }
  }
  esxi_high_cpu_detection {
    enabled = false
  }
  esxi_high_memory_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      compression_decompression_rate = 104
    }
  }
  guest_cpu_limit_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      host_cpu_usage_percentage = 74
      vm_cpu_ready_percentage   = 13
      vm_cpu_usage_percentage   = 94
    }
  }
  low_datastore_space_detection {
    enabled = false
  }
  overloaded_storage_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      command_aborts_number = 4
    }
  }
  slow_physical_storage_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      avg_read_write_latency  = 204
      peak_read_write_latency = 304
    }
  }
  undersized_storage_detection {
    enabled = false
  }
}
