data "dynatrace_mobile_application" "application" {
  name = "Application"
}

resource "dynatrace_calculated_mobile_metric" "metric" {
  name           = "#name#"
  enabled        = true
  app_identifier = data.dynatrace_mobile_application.application.id
  metric_key     = "calc:apps.mobile.#name#"
  metric_type    = "USER_ACTION_DURATION"
  dimensions {
    dimension {
      dimension = "APP_VERSION"
      top_x     = 10
    }
  }
}
