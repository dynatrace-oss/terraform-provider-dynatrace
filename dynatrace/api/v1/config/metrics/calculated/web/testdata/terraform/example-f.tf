data "dynatrace_application" "web_application" {
  name = "Web Application"
}

resource "dynatrace_calculated_web_metric" "dimensions" {
  name           = "#name#"
  enabled        = true
  app_identifier = data.dynatrace_application.web_application.id
  metric_key     = "calc:apps.web.#name#"
  dimensions {
    dimension {
      dimension    = "StringProperty"
      property_key = "web_utm_source"
      top_x        = 10
    }
    dimension {
      dimension = "GeoLocation"
      top_x     = 10
    }
  }
  metric_definition {
    metric = "UserActionDuration"
  }
}
