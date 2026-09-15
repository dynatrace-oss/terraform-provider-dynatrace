resource "dynatrace_user_action_metrics" "metric" {
  enabled    = true
  dimensions = [ "application" ]
  metric_key = "uacm.#name#"
  filters {
    filter {
      field_name = "type"
      operator   = "EQUALS"
      value      = "Xhr"
    }
  }
  value {
    type = "COUNTER"
  }
}