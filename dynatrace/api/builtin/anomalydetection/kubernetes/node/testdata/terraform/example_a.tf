resource "dynatrace_k8s_node_anomalies" "#name#" {
  scope = "environment"
  cpu_requests_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
      threshold                     = 95
    }
  }
  memory_requests_saturation {
    enabled = false
  }
  pods_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
      threshold                     = 95
    }
  }
  readiness_issues {
    enabled = false
  }
  node_problematic_condition {
    enabled = false
  }
}
