resource "dynatrace_k8s_cluster_anomalies" "#name#" {
  scope = "environment"
  cpu_requests_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 20
      sample_period_in_minutes      = 15
      threshold                     = 95
    }
  }
  memory_requests_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 20
      sample_period_in_minutes      = 15
      threshold                     = 95
    }
  }
  monitoring_issues {
    enabled = true
    configuration {
      observation_period_in_minutes = 35
      sample_period_in_minutes      = 20
    }
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
    enabled = true
    configuration {
      observation_period_in_minutes = 5
      sample_period_in_minutes      = 4
    }
  }
}