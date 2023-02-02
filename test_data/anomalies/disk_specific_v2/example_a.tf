resource "dynatrace_disk_specific_anomalies_v2" "#name#" {
  disk_id                                  = "DISK-1234567890000000"
  override_disk_low_space_detection        = true
  override_low_inodes_detection            = true
  override_slow_writes_and_reads_detection = true
  disk_low_inodes_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      free_inodes_percentage = 1
    }
  }
  disk_low_space_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      free_space_percentage = 1
    }
  }
  disk_slow_writes_and_reads_detection {
    enabled        = true
    detection_mode = "custom"
    custom_thresholds {
      write_and_read_time = 300
    }
  }
}
