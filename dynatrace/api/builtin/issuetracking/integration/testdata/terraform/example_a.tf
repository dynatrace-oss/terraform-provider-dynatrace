resource "dynatrace_issue_tracking" "#name#" {
  enabled            = true
  issuelabel         = "#name#"
  issuequery         = "{NAME}, {VERSION}"
  issuetheme         = "INFO"
  issuetrackersystem = "GITHUB"
  token              = "################"
  url                = "https://github.com/sampleorg/samplerepo"
  username           = "terraform-user"
}
