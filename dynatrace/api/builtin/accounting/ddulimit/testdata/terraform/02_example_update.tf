resource "dynatrace_ddu_pool" "#name#" {
  metrics {
    enabled = true
    type    = "MONTHLY"
    value   = 600
  }
  log_monitoring {
    enabled = true
    type    = "ANNUAL"
    value   = 124
  }
  events {
    enabled = false
  }
  serverless {
    enabled = false
  }
  traces {
    enabled = true
    type    = "MONTHLY"
    value   = 127
  }
}