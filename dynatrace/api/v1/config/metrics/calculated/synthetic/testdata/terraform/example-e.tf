# ID calc:synthetic.browser.easytravelbooknow.domcompletesplitbygeolocation
resource "dynatrace_calculated_synthetic_metric" "#name#" {
  name               = "#name#"
  enabled            = true
  metric             = "DOMComplete"
  metric_key         = "calc:synthetic.browser.#name#"
  monitor_identifier = "SYNTHETIC_TEST-74EEC98A3855C3DD"
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
