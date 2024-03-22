resource "dynatrace_container_rule" "first-instance" {
  enabled  = true
  mode     = "MONITORING_ON"
  operator = "NOT_CONTAINS"
  property = "CONTAINER_NAME"
  value    = "Terraform"
}

resource "dynatrace_container_rule" "second-instance" {
  enabled      = true
  mode         = "MONITORING_ON"
  operator     = "NOT_CONTAINS"
  property     = "CONTAINER_NAME"
  value        = "Terraform-second"
  insert_after = dynatrace_container_rule.first-instance.id
}
