resource "dynatrace_xmatters_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active   = false
  name     = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile  = data.dynatrace_alerting_profile.Default.id
  url      = "https://webhook.site/40bf4d43-1a50-4ebd-913d-bf50ce7c3a1e"
  insecure = true
  payload  = "x-matters-payload"
  headers {
    header {
      name  = "http-header-name-01"
      value = "http-header-value-01"
    }
    header {
      name         = "http-header-name-02"
      secret_value = "http-header-value-02"
    }
  }
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}