resource "dynatrace_slo" "#name#" {
  name              = "#name#"
  metric_name       = "#name#"
  evaluation        = "AGGREGATE"
  metric_expression = "(100)*((builtin:apps.web.action.speedIndex.load.browser:splitBy())/(builtin:apps.web.action.speedIndex.load.browser:splitBy()))"
  filter            = "type(\"APPLICATION_METHOD\")"
  target            = 99.98
  timeframe         = "-1w"
  warning           = 99.99
  error_budget_burn_rate {
    burn_rate_visualization_enabled = true
    fast_burn_threshold             = 15
  }
}