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

resource "dynatrace_service_external_web_request" "first-instance" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ dynatrace_management_zone_v2.my-mgmz.id ]
  conditions {
    condition {
      attribute              = "ApplicationId"
      compare_operation_type = "StringEquals"
      ignore_case            = false
      text_values            = [ "Terraform" ]
    }
  }
  id_contributors {
    port_for_service_id = true
    application_id {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OriginalValue"
      }
    }
    context_root {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OverrideValue"
        value_override {
          value = "Terraform"
        }
      }
    }
    public_domain_name {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type   = "TransformValue"
        copy_from_host_name = true
        transformations {
          transformation {
            transformation_type = "REMOVE_IPS"
          }
        }
      }
    }
  }
}

resource "dynatrace_service_external_web_request" "second-instance" {
  name             = "#name#-second"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ dynatrace_management_zone_v2.my-mgmz.id ]
  conditions {
    condition {
      attribute              = "ApplicationId"
      compare_operation_type = "StringEquals"
      ignore_case            = false
      text_values            = [ "Terraform-2" ]
    }
  }
  id_contributors {
    port_for_service_id = true
    application_id {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OriginalValue"
      }
    }
    context_root {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type = "OverrideValue"
        value_override {
          value = "Terraform"
        }
      }
    }
    public_domain_name {
      enable_id_contributor = true
      service_id_contributor {
        contribution_type   = "TransformValue"
        copy_from_host_name = true
        transformations {
          transformation {
            transformation_type = "REMOVE_IPS"
          }
        }
      }
    }
  }
  insert_after = dynatrace_service_external_web_request.first-instance.id
}
