resource "dynatrace_attack_allowlist" "first-instance" {
  enabled      = true
  insert_after = ""
  rule_name    = "#name#-first"
  attack_handling {
    blocking_strategy = "MONITOR"
  }
  metadata {
    comment = "Example"
  }
  resource_attribute_conditions {
    resource_attribute_condition {
      matcher                  = "STARTS_WITH"
      resource_attribute_key   = "AttributeKey2"
      resource_attribute_value = "AttributeValue2"
    }
    resource_attribute_condition {
      matcher                  = "EQUALS"
      resource_attribute_key   = "AttributeKey1"
      resource_attribute_value = "AttributeValue1"
    }
  }
  rules {
    rule {
      criteria_key                  = "DETECTION_TYPE"
      criteria_matcher              = "EQUALS"
      criteria_value_detection_type = "SSRF"
    }
    rule {
      criteria_key             = "ACTOR_IP"
      criteria_matcher         = "CONTAINS"
      criteria_value_free_text = "192.168.1.2"
    }
  }
}

resource "dynatrace_attack_allowlist" "second-instance" {
  enabled      = false
  rule_name    = "#name#-second"
  attack_handling {
    blocking_strategy = "OFF"
  }
  metadata {
    comment = ""
  }
  rules {
    rule {
      criteria_key             = "ACTOR_IP"
      criteria_matcher         = "CONTAINS"
      criteria_value_free_text = "192.168.1.1"
    }
  }
  insert_after = dynatrace_attack_allowlist.first-instance.id
}
