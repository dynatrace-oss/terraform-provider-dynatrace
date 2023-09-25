resource "dynatrace_attack_settings" "#name#" {
  enabled = true
  default_attack_handling {
    blocking_strategy_java = "MONITOR"
  }
}