resource "dynatrace_custom_anomalies" "#name#" {
  name                  = "#name#"
  description           = "The {metricname} value of {severity} was {alert_condition} the baseline of {baseline}."
  metric_selector       = "ghputoutgoing:filter(existsKey(\"dt.entity.service\"),in(\"dt.entity.service\",entitySelector(\"type(SERVICE),mzId(6734823652592292763)\"))):avg"
  enabled               = true
  primary_dimension_key = "dt.entity.service"
  severity              = "PERFORMANCE"
  warning_reason        = "NONE"
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
