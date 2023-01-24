resource "dynatrace_pager_duty_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active  = false
  name    = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile = dynatrace_alerting.Default.id
  account = "pager-duty-account"
  service = "pager-duty-service"
  api_key = "pager-duty-api-key"
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
