resource "dynatrace_user_action_metrics" "#name#" {
  enabled    = true
  dimensions = [ "application" ]
  metric_key = "uacm.TerraformTest"
  filters {
    filter {
      field_name = "type"
      operator   = "EQUALS"
      value      = "Xhr"
    }
  }
  value {
    type = "COUNTER"
  }
}