resource "dynatrace_k8s_workload_anomalies" "#name#" {
  scope = "environment"
  container_restarts {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
      threshold                     = 2
    }
  }
  deployment_stuck {
    enabled = true
    configuration {
      observation_period_in_minutes = 5
      sample_period_in_minutes      = 4
    }
  }
  not_all_pods_ready {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
    }
  }
  pending_pods {
    enabled = true
    configuration {
      observation_period_in_minutes = 16
      sample_period_in_minutes      = 11
      threshold                     = 2
    }
  }
  pod_stuck_in_terminating {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
    }
  }
  workload_without_ready_pods {
    enabled = true
    configuration {
      observation_period_in_minutes = 6
      sample_period_in_minutes      = 4
    }
  }
}
