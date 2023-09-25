resource "dynatrace_calculated_web_metric" "#name#" {
  name           = "#name#"
  enabled        = true
  app_identifier = "APPLICATION-EA7C4B59F27D43EB"
  metric_key     = "calc:apps.web.#name#"
  metric_definition {
    metric = "Apdex"
  }
}
