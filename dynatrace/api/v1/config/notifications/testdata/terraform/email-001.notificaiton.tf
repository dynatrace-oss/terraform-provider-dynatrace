resource "dynatrace_notification" "#name#" {
  email {
    name             = "#name#"
    active           = false
    alerting_profile = dynatrace_alerting_profile.Default.id
    bcc_receivers    = ["she@home.gov", "you@home.gov", "me@home.gov"]
    body             = "{ProblemDetailsHTML}"
    cc_receivers     = ["me@home.org", "she@home.org", "you@home.org"]
    receivers        = ["she@home.com", "you@home.com", "me@home.com"]
    subject          = "EMAIL-SUBJECT"
  }
}

resource "dynatrace_alerting_profile" "Default" {
  display_name = "#name#"
}
