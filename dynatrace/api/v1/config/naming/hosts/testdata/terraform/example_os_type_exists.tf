resource "dynatrace_host_naming" "host_naming_os_type_exists" {
  name    = "#name#"
  enabled = true
  format  = "{Host:DetectedName} - {HostGroup:Name}"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_OS_TYPE"
      }
      os_type {
        operator = "EXISTS"
      }
    }
  }
}
