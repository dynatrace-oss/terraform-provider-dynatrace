resource "dynatrace_calculated_web_metric" "user_action_properties" {
  name           = "#name#"
  enabled        = true
  app_identifier = dynatrace_web_application.application.id
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
    continent                         = data.dynatrace_synthetic_location.location.geo_location_id
    target_view_group_name_match_type = "Equals"
    target_view_name_match_type       = "Equals"
    user_action_properties {
      property {
        key        = "manifest_fetch_status"
        match_type = "Equals"
        value      = "success"
      }
    }
  }
}
