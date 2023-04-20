resource "dynatrace_application_detection_rule_v2" "#name#" {
  application_id = "APPLICATION-1234567890000000"
  matcher        = "DOMAIN_MATCHES"
  pattern        = "TerraformTest"
}