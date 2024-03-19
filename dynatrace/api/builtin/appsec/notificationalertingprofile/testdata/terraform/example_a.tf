resource "dynatrace_management_zone_v2" "my-mgmz" {
  name = "terraform-my-mgmz-001"
  rules {
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION_NAMESPACE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "extensions"
          }
        }
      }
    }
  }
}

resource "dynatrace_vulnerability_alerting" "#name#" {
  name                   = "#name#"
  enabled                = true
  enabled_risk_levels    = [ "LOW", "MEDIUM", "HIGH", "CRITICAL" ]
  enabled_trigger_events = [ "SECURITY_PROBLEM_OPENED" ]
  management_zone        = dynatrace_management_zone_v2.my-mgmz.id
}
