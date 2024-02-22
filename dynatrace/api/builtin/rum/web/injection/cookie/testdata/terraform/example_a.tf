resource "dynatrace_web_app_injection_cookie" "#name#" {
  application_id              = "APPLICATION-1234567890000000"
  same_site_cookie_attribute  = "STRICT"
  use_secure_cookie_attribute = true
}