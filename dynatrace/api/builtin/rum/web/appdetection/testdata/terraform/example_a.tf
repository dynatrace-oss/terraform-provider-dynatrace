resource "dynatrace_application_detection_rule_v2" "#name#" {
  application_id = "APPLICATION-90CE01F27D579187"
  matcher        = "DOMAIN_MATCHES"
  pattern        = "TerraformTest"
}