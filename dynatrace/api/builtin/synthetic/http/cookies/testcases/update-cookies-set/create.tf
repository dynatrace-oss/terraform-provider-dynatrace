resource "dynatrace_http_monitor_cookies" "cookies" {
  enabled = true
  scope   = "HTTP_CHECK-0000000000000000"
  cookies {
    cookie {
      name   = "SampleName1"
      value  = "SampleValue1"
      domain = "dynatrace.com"
    }
    # to update
    cookie {
      name   = "SampleName2"
      value  = "SampleValue2"
      domain = "dynatrace.com"
    }
  }
}
