resource "dynatrace_notification" "#name#" {
  email {
    name             = "#name#"
    active           = true
    alerting_profile = "f75e68ef-aca7-3a07-9c21-94eb00ecfc56"
    bcc_receivers    = ["they@home.com", "me@home.com", "you@home.com", "we@home.com", "us@home.com"]
    body             = "{ProblemDetailsHTML}"
    cc_receivers     = ["they@home.com", "you@home.com", "me@home.com", "we@home.com"]
    receivers        = ["they@home.com", "you@home.com", "me@home.com", "we@home.com"]
    subject          = "{State} Problem {ProblemID}: {ImpactedEntity}"
  }
}
