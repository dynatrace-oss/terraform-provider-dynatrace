resource "dynatrace_disk_anomalies" "#name#" {
  name              = "#name#"
  enabled           = true
  metric            = "LOW_DISK_SPACE"
  samples           = 5
  threshold         = 10
  violating_samples = 5
  disk_name {
    operator = "CONTAINS"
    value    = "888"
  }
}
