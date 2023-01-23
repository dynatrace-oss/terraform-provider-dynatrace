resource "dynatrace_notification" "#name#" {
  ops_genie {
    name             = "#name#"
    active           = false
    alerting_profile = "f75e68ef-aca7-3a07-9c21-94eb00ecfc56"
    api_key          = "#######"
    domain           = "jhjh"
    message          = "{ProblemImpact} Problem {ProblemID}: {ProblemTitle}"
  }
}
