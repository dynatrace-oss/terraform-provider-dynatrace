resource "dynatrace_host_naming" "#name#" {
  name = "#name#"
  enabled = true
  format = "{AwsAutoScalingGroup:Name}"
  conditions {
    condition {
      host_tech {
        negate = false
        operator = "EQUALS"
        value {
          type = "BOSH"
        }
      }
      key {
        type = "STATIC"
        attribute = "HOST_TECHNOLOGY"
      }
    }
    condition {
      integer {
        negate = false
        operator = "EQUALS"
        value = 3
      }
      key {
        type = "STATIC"
        attribute = "HOST_AIX_VIRTUAL_CPU_COUNT"
      }
    }
  }
}
