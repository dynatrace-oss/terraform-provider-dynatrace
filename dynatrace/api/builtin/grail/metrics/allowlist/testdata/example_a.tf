resource "dynatrace_grail_metrics_allowlist" "#name#" {
  allow_rules {
    allow_rule {
      enabled = false
      metric_key = "terraform"
      pattern = "CONTAINS"
    }
  }
}