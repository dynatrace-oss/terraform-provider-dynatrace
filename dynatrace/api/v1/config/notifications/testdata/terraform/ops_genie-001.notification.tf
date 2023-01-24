resource "dynatrace_notification" "#name#" {
  ops_genie {
    name             = "#name#"
    active           = false
    alerting_profile = dynatrace_alerting_profile.Default.id
    api_key          = "#######"
    domain           = "#name#"
    message          = "{ProblemImpact} Problem {ProblemID}: {ProblemTitle}"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
