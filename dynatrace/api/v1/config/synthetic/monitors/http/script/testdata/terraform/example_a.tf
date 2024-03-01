resource "dynatrace_http_monitor_script" "#name#" {
  http_id = "HTTP_CHECK-1234567890000000"
  script {
    request {
      description     = "request1"
      method          = "GET"
      url             = "http://httpstat.us/200"
      configuration {
        accept_any_certificate = true
      }
    }
    request {
      description     = "request2"
      method          = "GET"
      url             = "http://httpstat.us/400"
      configuration {
        accept_any_certificate = true
      }
    }
  }
}