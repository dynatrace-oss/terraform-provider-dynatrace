resource "dynatrace_email_notification" "#name#" { # replace #name# with the name you would like your resource be known within your Terraform Module
  active                 = false
  name                   = "#name#" # replace #name# with the name you would like your entry to be displayed within the Dynatrace Web UI
  profile                = data.dynatrace_alerting_profile.Default.id
  subject                = "EMAIL-SUBJECT"
  to                     = ["she@home.com", "me@home.com", "you@home.com"]
  cc                     = ["she@home.org", "me@home.org", "you@home.org"]
  bcc                    = ["she@home.gov", "me@home.gov", "you@home.gov"]
  notify_closed_problems = true
  body                   = "{ProblemDetailsHTML}"
}

data "dynatrace_alerting_profile" "Default" {
  name = "Default"
}