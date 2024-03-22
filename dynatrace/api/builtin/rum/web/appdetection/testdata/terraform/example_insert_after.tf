resource "dynatrace_application_detection_rule_v2" "first-instance" {
  application_id = "APPLICATION-90CE01F27D579187"
  matcher        = "DOMAIN_MATCHES"
  pattern        = "TerraformTest"
}

resource "dynatrace_application_detection_rule_v2" "second-instance" {
  application_id = "APPLICATION-90CE01F27D579187"
  matcher        = "DOMAIN_MATCHES"
  pattern        = "TerraformTest-2"
  insert_after = dynatrace_application_detection_rule_v2.first-instance.id
}