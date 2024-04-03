resource "dynatrace_pg_anomalies" "#name#" {
  pg_id = "PROCESS_GROUP-AA665B81183B0D7A"
  availability {
    method            = "OFF"
    minimum_threshold = 0
  }
}