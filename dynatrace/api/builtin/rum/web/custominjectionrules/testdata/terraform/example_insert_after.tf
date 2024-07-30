resource "dynatrace_web_app_custom_injection" "first-entry" {
  enabled        = false
  application_id = "APPLICATION-1234567890000000"
  operator       = "Starts"
  url_pattern    = "/terraform/afterhtml"
  rule           = "AfterSpecificHtml"
  html_pattern   = "example"
}

resource "dynatrace_web_app_custom_injection" "second-entry" {
  enabled        = false
  application_id = "APPLICATION-1234567890000000"
  operator       = "Contains"
  rule           = "DoNotInject"
  url_pattern    = "/terraform/donotinject"
  insert_after = "${dynatrace_web_app_custom_injection.first-entry.id}"
}