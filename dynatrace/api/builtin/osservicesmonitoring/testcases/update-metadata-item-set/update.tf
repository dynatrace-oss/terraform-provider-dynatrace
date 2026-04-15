resource "dynatrace_os_services" "services" {
  name                      = "#name#"
  enabled                   = true
  alert_activation_duration = 5
  alerting                  = true
  monitoring                = true
  not_installed_alerting    = true
  scope                     = "HOST-0000000000000000"
  status_condition_linux    = "$eq(failed)"
  system                    = "LINUX"
  detection_conditions_linux {
    linux_detection_condition {
      condition = "$contains(Terraform)"
      property  = "ServiceName"
    }
  }
  metadata {
    item {
      metadata_key   = "TerraformKey1"
      metadata_value = "TerraformValue1"
    }
    # updated
    item {
      metadata_key   = "TerraformKeyEdit"
      metadata_value = "TerraformValueEdit"
    }
  }
}

