data "dynatrace_application" "web_application" {
  name = "Web Application"
}

resource "dynatrace_application_detection_rule_v2" "detection_rule" {
  application_id = data.dynatrace_application.web_application.id
  matcher        = "DOMAIN_MATCHES"
  pattern        = "TerraformTest"
}
