resource "dynatrace_pg_anomalies" "#name#" {
  pg_id = "PROCESS_GROUP-XXXXXXXXXXXXXXXX"
  availability {
    method            = "MINIMUM_THRESHOLD"
    minimum_threshold = 5
  }
}