resource "dynatrace_calculated_synthetic_metric" "metric_b" {
  name               = "#name#"
  enabled            = true
  metric             = "JavaScriptErrors"
  metric_key         = "calc:synthetic.browser.#name#"
  monitor_identifier = dynatrace_browser_monitor.monitor.id
}
