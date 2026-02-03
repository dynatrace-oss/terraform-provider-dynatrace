resource "dynatrace_calculated_synthetic_metric" "metric_d" {
  name               = "#name#"
  enabled            = true
  metric             = "HTMLDownloaded"
  metric_key         = "calc:synthetic.browser.#name#"
  monitor_identifier = dynatrace_browser_monitor.monitor.id
  dimensions {
    dimension {
      dimension = "Location"
      top_x     = 100
    }
  }
}
