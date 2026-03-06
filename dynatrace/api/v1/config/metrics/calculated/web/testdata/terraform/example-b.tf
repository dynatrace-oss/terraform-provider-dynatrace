resource "dynatrace_calculated_web_metric" "apdex" {
  name           = "#name#"
  enabled        = true
  app_identifier = dynatrace_web_application.application.id
  metric_key     = "calc:apps.web.#name#"
  metric_definition {
    metric = "Apdex"
  }
}
