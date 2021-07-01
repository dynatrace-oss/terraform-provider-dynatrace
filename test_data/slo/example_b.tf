resource "dynatrace_slo" "#name#" {
  name = "#name#"
  denominator = "builtin:apps.web.action.speedIndex.load.browser:splitBy()"
  evaluation = "AGGREGATE"
  filter = "type(\"APPLICATION_METHOD\")"
  numerator = "builtin:apps.web.action.speedIndex.load.browser:splitBy()"
  target = 99.98
  timeframe = "-5m"
  warning = 99.99
}
