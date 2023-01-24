resource "dynatrace_notification" "#name#" {
  jira {
    name             = "#name#"
    description      = "jira-description"
    active           = false
    alerting_profile = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    issue_type       = "jira-issue-type"
    password         = "#######"
    project_key      = "jira-project-key"
    summary          = "jira-summary"
    url              = "https://localhost:9999/jira"
    username         = "jira-user-name"
  }
}
