resource "dynatrace_calculated_synthetic_metric" "metric_e" {
  name               = "#name#"
  enabled            = true
  metric             = "DOMComplete"
  metric_key         = "calc:synthetic.browser.#name#"
  monitor_identifier = dynatrace_browser_monitor.monitor.id
  dimensions {
    dimension {
      dimension = "Location"
      top_x     = 100
    }
  }
  filter {
    action_type = "Load"
  }
}
