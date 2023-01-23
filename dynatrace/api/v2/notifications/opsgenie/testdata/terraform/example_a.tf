resource "dynatrace_ops_genie_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active  = false
  name    = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile = data.dynatrace_alerting_profile.Default.id
  domain  = "ps-genie-domain"
  message = "ops-genie-message"
  api_key = "ops-genie-api-key"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}