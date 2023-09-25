resource "dynatrace_calculated_mobile_metric" "#name#" {
  name           = "#name#"
  enabled        = true
  app_identifier = "MOBILE_APPLICATION-7F6AE72450E14F11"
  metric_key     = "calc:apps.mobile.#name#"
  metric_type    = "USER_ACTION_DURATION"
  dimensions {
    dimension {
      dimension = "MANUFACTURER"
      top_x     = 10
    }
  }
}
