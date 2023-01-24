resource "dynatrace_notification" "#name#" {
  jira {
    name             = "#name#"
    description      = "jira-description"
    active           = false
    alerting_profile = dynatrace_alerting_profile.Default.id
    issue_type       = "jira-issue-type"
    password         = "#######"
    project_key      = "#name#"
    summary          = "jira-summary"
    url              = "https://localhost:9999/jira"
    username         = "jira-user-name"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
