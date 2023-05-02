resource "dynatrace_disk_anomaly_rules" "#name#" {
  name              = "#name#"
  enabled           = false
  host_group_id     = "HOST_GROUP-1234567890000000"
  metric            = "LOW_DISK_SPACE"
  threshold_percent = 10
  disk_name_filter {
    operator = "CONTAINS"
    value    = "#name#"
  }
  sample_limit {
    samples           = 3
    violating_samples = 3
  }
}