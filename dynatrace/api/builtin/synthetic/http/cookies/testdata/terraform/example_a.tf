resource "dynatrace_http_monitor_cookies" "#name#" {
  enabled = true
  scope   = "HTTP_CHECK-1234567890000000"
  cookies {
    cookie {
      name   = "SampleName"
      domain = "google.com"
      value  = "SampleValue"
    }
  }
}
