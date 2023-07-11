resource "dynatrace_ops_genie_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active  = false
  name    = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile = dynatrace_alerting.Default.id
  domain  = "#name#"
  message = "ops-genie-message"
  api_key = "ops-genie-api-key"
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
