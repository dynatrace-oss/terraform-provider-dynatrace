resource "dynatrace_web_app_custom_injection" "APPLICATION-1234567890000000" {
  enabled        = false
  application_id = "APPLICATION-1234567890000000"
  operator       = "Starts"
  url_pattern    = "/terraform"
  rule           = "AfterSpecificHtml"
  html_pattern   = "example"
}