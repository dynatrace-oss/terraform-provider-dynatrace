resource "dynatrace_calculated_web_metric" "#name#" {
  name           = "#name#"
  enabled        = true
  app_identifier = "APPLICATION-EA7C4B59F27D43EB"
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