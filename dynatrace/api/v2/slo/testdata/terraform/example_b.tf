resource "dynatrace_slo" "#name#" {
  name = "#name#"
  metric_expression = "(100)*((builtin:apps.web.action.speedIndex.load.browser:splitBy())/(builtin:apps.web.action.speedIndex.load.browser:splitBy()))"
  evaluation = "AGGREGATE"
  filter = "type(\"APPLICATION_METHOD\")"
  target = 99.98
  timeframe = "-5m"
  warning = 99.99
}
