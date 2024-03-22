resource "dynatrace_issue_tracking" "first-instance" {
  enabled            = true
  issuelabel         = "#name#"
  issuequery         = "{NAME}, {VERSION}"
  issuetheme         = "INFO"
  issuetrackersystem = "GITHUB"
  token              = "################"
  url                = "https://github.com/sampleorg/samplerepo"
  username           = "terraform-user"
}

resource "dynatrace_issue_tracking" "second-instance" {
  enabled            = true
  issuelabel         = "#name#-second"
  issuequery         = "{NAME}, {VERSION}"
  issuetheme         = "INFO"
  issuetrackersystem = "GITHUB"
  token              = "################-second"
  url                = "https://github.com/sampleorg/samplerepo-second"
  username           = "terraform-user"
  insert_after       = dynatrace_issue_tracking.first-instance.id
}
