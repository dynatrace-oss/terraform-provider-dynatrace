resource "dynatrace_attack_rules" "#name#" {
  criteria {
    attack_type = "ANY"
  }
  enabled = true
  metadata {
    comment = "#name#"
  }
  attack_handling {
    blocking_strategy = "MONITOR"
  }
}