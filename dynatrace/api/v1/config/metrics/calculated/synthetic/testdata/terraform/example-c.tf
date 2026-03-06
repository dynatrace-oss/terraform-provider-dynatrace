resource "dynatrace_calculated_synthetic_metric" "metric_c" {
  name               = "#name#"
  enabled            = true
  metric             = "HttpErrors"
  metric_key         = "calc:synthetic.browser.#name#"
  monitor_identifier = dynatrace_browser_monitor.monitor.id
}
