resource "dynatrace_jira_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active      = false
  name        = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile     = data.dynatrace_alerting_profile.Default.id
  url         = "https://localhost:9999/jira"
  username    = "jira-user-name"
  api_token   = "jira-api-token"
  project_key = "jira-project-key"
  issue_type  = "jira-issue-type"
  summary     = "jira-summary"
  description = "jira-description"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}