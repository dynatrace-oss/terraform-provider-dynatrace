resource "dynatrace_attack_allowlist" "first-instance" {
  criteria {
    source_ip = "192.168.0.1"
  }
  enabled = false
  attack_handling {
    blocking_strategy = "MONITOR"
  }
  metadata {
    comment = ""
  }
}

resource "dynatrace_attack_allowlist" "second-instance" {
  criteria {
    source_ip = "192.168.0.2"
  }
  enabled = false
  attack_handling {
    blocking_strategy = "MONITOR"
  }
  metadata {
    comment = ""
  }
  insert_after = dynatrace_attack_allowlist.first-instance.id
}
