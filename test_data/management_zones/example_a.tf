resource "dynatrace_management_zone" "#name#" {
  name = "#name#"
  rules {
    conditions {
      host_tech {
        negate = true
        value {
          type = "APPARMOR"
        }
        operator = "EQUALS"
      }
      key {
        attribute = "HOST_TECHNOLOGY"
      }
    }
    conditions {
      tag {
        #        negate = false
        value {
          context = "CONTEXTLESS"
          key     = "Asddf"
        }
        operator = "TAG_KEY_EQUALS"
      }
      key {
        attribute = "HOST_TAGS"
      }
    }
    enabled = true
    type    = "HOST"
  }
}
