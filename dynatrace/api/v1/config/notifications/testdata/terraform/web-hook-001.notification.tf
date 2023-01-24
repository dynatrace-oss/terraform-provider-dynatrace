resource "dynatrace_notification" "#name#" {
  web_hook {
    name                   = "#name#"
    accept_any_certificate = true
    active                 = false
    alerting_profile       = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    # notify_event_merges  = false
    payload                = "test6"
    url                    = "https://webhook.site/eba95d71-d7a2-4e52-9a5b-d8e52869fb6d"
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