resource "dynatrace_pg_anomalies" "#name#" {
  process_group = "PROCESS_GROUP-XXXXXXXXXXXXXXXX"
  availability {
    method            = "MINIMUM_THRESHOLD"
    minimum_threshold = 5
  }
}