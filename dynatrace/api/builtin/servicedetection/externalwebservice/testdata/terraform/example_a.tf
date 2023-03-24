resource "dynatrace_service_external_web_service" "#name#" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = false
  management_zones = [ "000000000000000000" ]
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