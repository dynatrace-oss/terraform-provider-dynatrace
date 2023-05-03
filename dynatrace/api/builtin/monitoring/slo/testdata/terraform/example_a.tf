resource "dynatrace_slo_v2" "#name#" {
  name               = "#name#"
  enabled            = true
  custom_description = "Terraform Test"
  evaluation_type    = "AGGREGATE"
  evaluation_window  = "-1w"
  filter             = "type(SERVICE),serviceType(WEB_SERVICE,WEB_REQUEST_SERVICE)"
  metric_expression  = "100*(builtin:service.requestCount.server:splitBy())/(builtin:service.requestCount.server:splitBy())"
  metric_name        = "terraform_test"
  target_success     = 95
  target_warning     = 98
  error_budget_burn_rate {
    burn_rate_visualization_enabled = false
  }
}