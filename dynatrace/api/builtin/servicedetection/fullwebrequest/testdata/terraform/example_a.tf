resource "dynatrace_service_full_web_request" "#name#" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ "000000000000000000" ]
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
