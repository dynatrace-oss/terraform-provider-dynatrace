resource "dynatrace_disk_anomalies_v2" "#name#" {
  scope = "environment"
  disk {
    disk_low_inodes_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        free_inodes_percentage = 5
      }
    }
    disk_low_space_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        free_space_percentage = 3
      }
    }
    disk_slow_writes_and_reads_detection {
      enabled        = true
      detection_mode = "custom"
      custom_thresholds {
        write_and_read_time = 200
      }
    }
  }
}
