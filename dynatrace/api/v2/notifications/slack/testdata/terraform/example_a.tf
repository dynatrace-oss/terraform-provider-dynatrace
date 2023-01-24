resource "dynatrace_slack_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active  = false
  name    = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile = data.dynatrace_alerting_profile.Default.id
  url     = "https://slack.home.com"
  channel = "slack-channel"
  message = "slack-message"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}