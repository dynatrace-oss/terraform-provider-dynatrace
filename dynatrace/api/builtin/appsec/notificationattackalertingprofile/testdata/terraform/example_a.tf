resource "dynatrace_appsec_attack_alerting" "#name#" {
  name                       = "#name#"
  enabled                    = true
  enabled_attack_mitigations = [ "NONE_ALLOWLISTED", "BLOCKED_WITH_EXCEPTION", "NONE_BLOCKING_DISABLED" ]
}