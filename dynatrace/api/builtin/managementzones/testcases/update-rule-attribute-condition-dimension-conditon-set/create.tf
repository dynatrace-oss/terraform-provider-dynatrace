resource "dynatrace_management_zone_v2" "zone" {
  name = "#name#"
  rules {
    # not updated
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

    # rule will be directly updated
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "linkerd"
          }
        }
      }
    }

    # will be updated on attribute_rule
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CUSTOM_DEVICE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "CUSTOM_DEVICE_NAME"
            operator       = "CONTAINS"
            string_value   = "gcp"
          }
          condition {
            case_sensitive = true
            key            = "CUSTOM_DEVICE_NAME"
            operator       = "CONTAINS"
            string_value   = "aws"
          }
        }
      }
    }

    # will be updated on dimension_rule
    rule {
      type            = "DIMENSION"
      enabled         = true
      entity_selector = ""
      dimension_rule {
        applies_to = "ANY"
        dimension_conditions {
          condition {
            condition_type = "DIMENSION"
            key            = "cloud.gcp"
            rule_matcher   = "BEGINS_WITH"
            value          = "cloud.gcp."
          }
          condition {
            condition_type = "DIMENSION"
            key            = "cloud.aws"
            rule_matcher   = "BEGINS_WITH"
            value          = "cloud.aws."
          }
        }
      }
    }
  }
}
