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

resource "dynatrace_service_external_web_service" "first-instance" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ dynatrace_management_zone_v2.my-mgmz.id ]
  conditions {
    condition {
      attribute              = "UrlPath"
      compare_operation_type = "StringContains"
      ignore_case            = true
      text_values            = [ "Terraform" ]
    }
  }
  id_contributors {
    detect_as_web_request_service = false
    port_for_service_id           = true
    url_path {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OverrideValue"
        value_override {
          value = "Terraform"
        }
      }
    }
  }
}

resource "dynatrace_service_external_web_service" "second-instance" {
  name             = "#name#-second"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ dynatrace_management_zone_v2.my-mgmz.id ]
  conditions {
    condition {
      attribute              = "UrlPath"
      compare_operation_type = "StringContains"
      ignore_case            = true
      text_values            = [ "Terraform-3" ]
    }
  }
  id_contributors {
    detect_as_web_request_service = false
    port_for_service_id           = true
    url_path {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OverrideValue"
        value_override {
          value = "Terraform"
        }
      }
    }
  }
  insert_after = dynatrace_service_external_web_service.first-instance.id
}