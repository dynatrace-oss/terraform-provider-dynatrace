resource "dynatrace_webhook_notification" "notification" {
  active                 = false
  name                   = "#name#"
  profile                = dynatrace_alerting.Default.id
  url                    = "https://webhook.site/#name#"
  insecure               = true
  notify_event_merges    = true
  notify_closed_problems = true
  payload                = "web-hook-payload"
  headers {
    header {
      name  = "http-header-name-01"
      value = "http-header-value-01"
    }
    header {
      name         = "http-header-name-02"
      secret_value = "http-header-value-02"
    }
    header {
      name  = "http-header-name-03"
      value = "http-header-value-03"
    }
  }
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
