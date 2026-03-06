resource "dynatrace_calculated_web_metric" "user_action_duration" {
  name           = "#name#"
  enabled        = true
  app_identifier = dynatrace_web_application.application.id
  metric_key     = "calc:apps.web.#name#"
  metric_definition {
    metric = "UserActionDuration"
  }
  user_action_filter {
    target_view_group_name_match_type = "Equals"
    target_view_name_match_type       = "Equals"
    user_action_name                  = "Loading of page /easytravel/login"
  }
}
