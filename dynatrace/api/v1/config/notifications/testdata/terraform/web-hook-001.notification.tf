resource "dynatrace_notification" "#name#" {
  web_hook {
    name                   = "#name#"
    accept_any_certificate = true
    active                 = false
    alerting_profile       = dynatrace_alerting_profile.Default.id
    # notify_event_merges  = false
    payload = "test6"
    url     = "https://webhook.site/#name#"
    header {
      name  = "http-header-name-01"
      value = "http-header-value-01"
    }
    header {
      name  = "http-header-name-02"
      value = "asdfadsf"
    }
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
