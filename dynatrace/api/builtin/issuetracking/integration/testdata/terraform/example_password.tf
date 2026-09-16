resource "dynatrace_issue_tracking" "jira_password" {
  enabled            = true
  issuelabel         = "#name#"
  issuequery         = "{NAME}, {VERSION}"
  issuetheme         = "INFO"
  issuetrackersystem = "JIRA"
  password           = "################"
  url                = "https://www.atlassian.com/"
  username           = "terraform-user"
}
