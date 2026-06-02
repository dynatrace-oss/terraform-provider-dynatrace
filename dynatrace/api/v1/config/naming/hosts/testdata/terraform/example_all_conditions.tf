resource "dynatrace_host_naming" "all_conditions" {
  name    = "#name#"
  enabled = true
  format  = "{AwsAutoScalingGroup:Name}"
  conditions {
    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_AZURE_COMPUTE_MODE"
      }
      azure_compute_mode {
        negate   = true
        operator = "EQUALS"
        value    = "SHARED"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_AZURE_SKU"
      }
      azure_sku {
        negate   = true
        operator = "EQUALS"
        value    = "BASIC"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_BITNESS"
      }
      bitness {
        negate   = true
        operator = "EQUALS"
        value    = "64"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_CLOUD_TYPE"
      }
      cloud_type {
        negate   = true
        operator = "EQUALS"
        value    = "AZURE"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_HYPERVISOR_TYPE"
      }
      hypervisor {
        negate   = true
        operator = "EQUALS"
        value    = "HYPER_V"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_IP_ADDRESS"
      }
      ipaddress {
        negate   = true
        operator = "EQUALS"
        value    = "value"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_OS_TYPE"
      }
      os_type {
        negate   = true
        operator = "EQUALS"
        value    = "LINUX"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_PAAS_TYPE"
      }
      paas_type {
        negate   = true
        operator = "EQUALS"
        value    = "KUBERNETES"
      }
    }

    condition {
      key {
        type      = "STATIC"
        attribute = "HOST_TAGS"
      }
      tag {
        negate   = true
        operator = "EQUALS"
        value {
          context = "CONTEXTLESS"
          key     = "key"
          value   = "value"
        }
      }
    }
  }
}

