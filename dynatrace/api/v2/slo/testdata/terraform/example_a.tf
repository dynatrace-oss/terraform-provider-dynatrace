resource "dynatrace_slo" "#name#" {
  name = "#name#"
  evaluation = "AGGREGATE"
  filter = "type(\"APPLICATION_METHOD\")"
  metric_expression = "builtin:apps.web.action.speedIndex.load.browser:splitBy()"
  target = 99.58
  timeframe = "-5m"
  warning = 99.99
}
