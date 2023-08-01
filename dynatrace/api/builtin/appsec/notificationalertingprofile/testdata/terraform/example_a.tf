resource "dynatrace_appsec_vulnerability_alerting" "#name#" {
  name                   = "#name#"
  enabled                = true
  enabled_risk_levels    = [ "LOW", "MEDIUM", "HIGH", "CRITICAL" ]
  enabled_trigger_events = [ "SECURITY_PROBLEM_OPENED" ]
  management_zone        = "000000000000000000"
}
