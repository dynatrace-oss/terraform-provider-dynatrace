resource "dynatrace_pager_duty_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active  = false
  name    = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile = data.dynatrace_alerting_profile.Default.id
  account = "pager-duty-account"
  service = "pager-duty-service"
  api_key = "pager-duty-api-key"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}