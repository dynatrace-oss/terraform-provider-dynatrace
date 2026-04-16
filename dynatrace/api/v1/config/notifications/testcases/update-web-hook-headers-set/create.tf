resource "dynatrace_notification" "notification" {
  web_hook {
    name                   = "#name#"
    accept_any_certificate = true
    active                 = false
    alerting_profile       = dynatrace_alerting_profile.Default.id
    payload = "test6"
    url     = "https://dynatrace.com/#name#"
    header {
      name  = "http-header-name-01"
      value = "http-header-value-01"
    }
    # to update
    header {
      name  = "http-header-name-02"
      value = "http-header-value-02"
    }
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
