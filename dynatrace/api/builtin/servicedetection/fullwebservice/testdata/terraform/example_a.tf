resource "dynatrace_management_zone_v2" "my-mgmz" {
  name = "#name#"
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

resource "dynatrace_service_full_web_service" "#name#" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = true
  management_zones = [ dynatrace_management_zone_v2.my-mgmz.id ]
  conditions {
    condition {
      attribute              = "HostName"
      compare_operation_type = "StringEndsWith"
      ignore_case            = true
      text_values            = [ "Terraform" ]
    }
  }
  id_contributors {
    detect_as_web_request_service = true
  }
}
