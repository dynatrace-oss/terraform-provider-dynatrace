resource "dynatrace_k8s_namespace_anomalies" "#name#" {
  scope = "environment"
  cpu_limits_quota_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 15
      sample_period_in_minutes      = 10
      threshold                     = 90
    }
  }
  cpu_requests_quota_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 15
      sample_period_in_minutes      = 10
      threshold                     = 90
    }
  }
  memory_limits_quota_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 15
      sample_period_in_minutes      = 10
      threshold                     = 90
    }
  }
  memory_requests_quota_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 15
      sample_period_in_minutes      = 10
      threshold                     = 90
    }
  }
  pods_quota_saturation {
    enabled = true
    configuration {
      observation_period_in_minutes = 15
      sample_period_in_minutes      = 10
      threshold                     = 90
    }
  }
}
