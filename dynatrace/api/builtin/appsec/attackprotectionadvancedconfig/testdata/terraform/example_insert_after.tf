resource "dynatrace_attack_rules" "first-instance" {
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

resource "dynatrace_attack_rules" "second-instance" {
  criteria {
    attack_type = "ANY"
  }
  enabled = true
  metadata {
    comment = "#name#-second"
  }
  attack_handling {
    blocking_strategy = "MONITOR"
  }
  insert_after = "${dynatrace_attack_rules.first-instance.id}"
}