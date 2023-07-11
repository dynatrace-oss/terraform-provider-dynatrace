resource "dynatrace_webhook_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active                 = false
  name                   = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
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
      name         = "http-header-name-03"
      secret_value = "http-header-value-03"
    }
    header {
      name         = "http-header-name-04"
      secret_value = "http-header-value-04"
    }
    header {
      name         = "http-header-name-05"
      secret_value = "http-header-value-05"
    }
    header {
      name         = "http-header-name-06"
      secret_value = "http-header-value-06"
    }
  }
  use_oauth_2            = true
  oauth_2_credentials {
    access_token_url = "https://www.google.com"
    client_id        = "terraform"
    client_secret    = "#######"
  }
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
