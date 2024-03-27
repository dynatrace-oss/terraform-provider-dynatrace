resource "dynatrace_os_services" "first-instance" {
  name                      = "#name#"
  enabled                   = true
  alert_activation_duration = 5
  alerting                  = true
  monitoring                = true
  not_installed_alerting    = true
  scope                     = "HOST-1234567890000000"
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
      metadata_key   = "TerraformKey"
      metadata_value = "TerraformValue"
    }
  }
}

resource "dynatrace_os_services" "second-instance" {
  name                      = "#name#-second"
  enabled                   = true
  alert_activation_duration = 5
  alerting                  = true
  monitoring                = true
  not_installed_alerting    = true
  scope                     = "HOST-1234567890000000"
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
      metadata_key   = "TerraformKey"
      metadata_value = "TerraformValue"
    }
  }
  insert_after = dynatrace_os_services.first-instance.id
}

