resource "dynatrace_notification" "#name#" {
  xmatters {
    name                   = "#name#"
    accept_any_certificate = true
    active                 = false
    alerting_profile       = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    payload                = "x-matters-payload"
    url                    = "https://webhook.site/40bf4d43-1a50-4ebd-913d-bf50ce7c3a1e"
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