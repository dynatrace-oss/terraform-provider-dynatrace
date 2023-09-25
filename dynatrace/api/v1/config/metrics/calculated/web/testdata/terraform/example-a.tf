resource "dynatrace_calculated_web_metric" "#name#" {
  name           = "#name#"
  enabled        = true
  app_identifier = "APPLICATION-EA7C4B59F27D43EB"
  metric_key     = "calc:apps.web.#name#"
  dimensions {
    dimension {
      dimension    = "StringProperty"
      property_key = "web_utm_campaign"
      top_x        = 10
    }
  }
  metric_definition {
    metric = "VisuallyComplete"
  }
  user_action_filter {
    continent                         = "GEOLOCATION-970B6D0A98F55995"
    target_view_group_name_match_type = "Equals"
    target_view_name_match_type       = "Equals"
  }
}
