resource "dynatrace_jira_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active      = false
  name        = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile     = dynatrace_alerting.Default.id
  url         = "https://localhost:9999/jira/#name#"
  username    = "jira-user-name"
  api_token   = "jira-api-token"
  project_key = "#name#"
  issue_type  = "jira-issue-type"
  summary     = "jira-summary"
  description = "jira-description"
}

resource "dynatrace_alerting" "Default" {
  name = "#name#"
}
