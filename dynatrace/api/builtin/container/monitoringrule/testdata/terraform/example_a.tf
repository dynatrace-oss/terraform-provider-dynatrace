resource "dynatrace_container_rule" "#name#" {
  enabled  = true
  mode     = "MONITORING_ON"
  operator = "NOT_CONTAINS"
  property = "CONTAINER_NAME"
  value    = "Terraform"
}