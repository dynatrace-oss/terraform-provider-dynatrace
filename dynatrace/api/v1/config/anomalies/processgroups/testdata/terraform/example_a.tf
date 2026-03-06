variable "PROCESS_GROUP_ID" {
  type = string
}

resource "dynatrace_pg_anomalies" "anomaly" {
  pg_id = var.PROCESS_GROUP_ID
  availability {
    method            = "OFF"
    minimum_threshold = 0
  }
}
