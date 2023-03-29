resource "dynatrace_service_full_web_service" "#name#" {
  name             = "#name#"
  description      = "Created by Terraform"
  enabled          = true
  management_zones = [ "000000000000000000" ]
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
