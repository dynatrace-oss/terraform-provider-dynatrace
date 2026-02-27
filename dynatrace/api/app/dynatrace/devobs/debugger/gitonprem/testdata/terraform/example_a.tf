resource "dynatrace_devobs_git_onprem" "onprem" {
  git_provider = "GithubOnPrem"
  url          = "https://example.com/test/#name#"
}
