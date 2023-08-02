resource "dynatrace_attack_allowlist" "#name#" {
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