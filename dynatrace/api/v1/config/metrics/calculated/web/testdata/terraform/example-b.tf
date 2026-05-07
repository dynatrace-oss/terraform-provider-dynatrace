data "dynatrace_application" "web_application" {
  name = "Web Application"
}

resource "dynatrace_calculated_web_metric" "apdex" {
  name           = "#name#"
  enabled        = true
  app_identifier = data.dynatrace_application.web_application.id
  metric_key     = "calc:apps.web.#name#"
  metric_definition {
    metric = "Apdex"
  }
}
