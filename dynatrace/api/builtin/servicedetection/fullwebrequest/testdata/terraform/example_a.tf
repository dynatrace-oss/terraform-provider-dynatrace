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

resource "dynatrace_service_full_web_request" "#name#" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ dynatrace_management_zone_v2.my-mgmz.id ]
  conditions {
    condition {
      attribute              = "UrlPath"
      compare_operation_type = "StringStartsWith"
      ignore_case            = true
      text_values            = [ "Terraform" ]
    }
  }
  id_contributors {
    application_id {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OverrideValue"
        value_override {
          value = "Terraform"
        }
      }
    }
    context_root {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "TransformURL"
        segment_count     = 2
        transformations {
          transformation {
            include_hex_numbers = true
            min_digit_count     = 1
            transformation_type = "REMOVE_NUMBERS"
          }
        }
      }
    }
    server_name {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OriginalValue"
      }
    }
  }
}
