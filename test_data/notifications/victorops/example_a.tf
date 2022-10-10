resource "dynatrace_victor_ops_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active      = false
  name        = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile     = data.dynatrace_alerting_profile.Default.id
  api_key     = "victor-ops-api-key"
  routing_key = "victor-ops-routing-key"
  message     = "victor-ops-message"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}