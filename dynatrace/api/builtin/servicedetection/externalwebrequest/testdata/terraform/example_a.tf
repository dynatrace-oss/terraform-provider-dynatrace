resource "dynatrace_service_external_web_request" "#name#" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ "000000000000000000" ]
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
