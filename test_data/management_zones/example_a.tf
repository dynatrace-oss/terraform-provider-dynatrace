resource "dynatrace_management_zone" "#name#" {
  name = "#name#"
  rules {
    conditions {
      simple_host_tech_comparison {
        type   = "SIMPLE_HOST_TECH"
        negate = false
        value {
          type = "APPARMOR"
        }
        operator = "EQUALS"
      }
      base_condition_key {
        attribute = "HOST_TECHNOLOGY"
      }
    }
    conditions {
      tag_comparison {
        type   = "TAG"
        negate = false
        value {
          context = "CONTEXTLESS"
          key     = "Asddf"
        }
        operator = "TAG_KEY_EQUALS"
      }
      base_condition_key {
        attribute = "HOST_TAGS"
      }
    }
    enabled = true
    type    = "HOST"
  }
}
