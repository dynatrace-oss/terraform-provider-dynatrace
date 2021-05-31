resource "dynatrace_custom_anomalies" "#name#" {
  name                  = "#name#"
  description           = "The {metricname} value of {severity} was {alert_condition} the baseline of {baseline}."
  enabled               = true
  aggregation_type      = "AVG"
  disabled_reason       = "NONE"
  metric_id             = "ghputoutgoing"
  primary_dimension_key = "dt.entity.service"
  severity              = "PERFORMANCE"
  warning_reason        = "NONE"
  scopes {
    management_zone {
      id = "6734823652592292763"
    }
  }
  strategy {
    auto {
      alert_condition          = "ABOVE"
      alerting_on_missing_data = false
      dealerting_samples       = 5
      samples                  = 5
      signal_fluctuations      = 1
      violating_samples        = 3
    }
  }
}
