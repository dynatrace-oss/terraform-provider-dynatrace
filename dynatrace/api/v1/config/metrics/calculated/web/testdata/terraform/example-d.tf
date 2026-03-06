resource "dynatrace_calculated_web_metric" "application_cache" {
  name           = "#name#"
  enabled        = true
  app_identifier = dynatrace_web_application.application.id
  metric_key     = "calc:apps.web.#name#"
  metric_definition {
    metric = "ApplicationCache"
  }
  user_action_filter {
    load_action                       = true
    target_view_group                 = "/easytravel/home"
    target_view_group_name_match_type = "Equals"
    target_view_name_match_type       = "Equals"
  }
}
