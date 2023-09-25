# ID calc:synthetic.browser.targethomepageco.htmldownloaded_responseend_splitbygeolocation
resource "dynatrace_calculated_synthetic_metric" "#name#" {
  name               = "#name#"
  enabled            = true
  metric             = "HTMLDownloaded"
  metric_key         = "calc:synthetic.browser.#name#"
  monitor_identifier = "SYNTHETIC_TEST-147CFF44DDB25C05"
  dimensions {
    dimension {
      dimension = "Location"
      top_x     = 100
    }
  }
}
