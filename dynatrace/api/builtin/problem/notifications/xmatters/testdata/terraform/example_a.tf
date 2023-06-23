resource "dynatrace_xmatters_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active   = false
  name     = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile  = dynatrace_alerting.Default.id
  url      = "https://webhook.site/#name#"
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

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
