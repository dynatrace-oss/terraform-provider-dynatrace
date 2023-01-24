resource "dynatrace_slack_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active  = false
  name    = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile = dynatrace_alerting.Default.id
  url     = "https://slack.home.com"
  channel = "#name#"
  message = "slack-message"
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
