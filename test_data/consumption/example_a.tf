resource "dynatrace_ddu_pool" "#name#" {
  metrics {
    enabled = true
    type    = "MONTHLY"
    value   = 123
  }
  log_monitoring {
    enabled = true
    type    = "MONTHLY"
    value   = 124
  }
  events {
    enabled = true
    type    = "MONTHLY"
    value   = 125
  }
  serverless {
    enabled = true
    type    = "MONTHLY"
    value   = 126
  }
  traces {
    enabled = true
    type    = "MONTHLY"
    value   = 127
  }
}