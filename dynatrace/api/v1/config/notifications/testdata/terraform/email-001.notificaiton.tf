resource "dynatrace_notification" "#name#" {
  email {
    name             = "#name#"
    active           = false
    alerting_profile = "c21f969b-5f03-333d-83e0-4f8f136e7682"
    bcc_receivers    = [ "she@home.gov", "you@home.gov", "me@home.gov" ]
    body             = "{ProblemDetailsHTML}"
    cc_receivers     = [ "me@home.org", "she@home.org", "you@home.org" ]
    receivers        = [ "she@home.com", "you@home.com", "me@home.com" ]
    subject          = "EMAIL-SUBJECT"
  }
}
