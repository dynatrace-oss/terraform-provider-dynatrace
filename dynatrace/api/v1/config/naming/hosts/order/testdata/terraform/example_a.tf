resource "dynatrace_host_naming_order" "this" {
  naming_rule_ids = [
    dynatrace_host_naming.first.id,
    dynatrace_host_naming.second.id,
  ]  
}

resource "dynatrace_host_naming" "first" {
  name = "${randomize}"
  enabled = true
  format = "{AwsAutoScalingGroup:Name} ${randomize}"
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
  }
}

resource "dynatrace_host_naming" "second" {
  name = "${randomize}"
  enabled = true
  format = "{AwsAutoScalingGroup:Name} ${randomize}"
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
  }
}
