resource "dynatrace_management_zone" "#name#" {
  name = "#name#"
  rules {
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
        type      = "STATIC"
      }
    }    
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
        type      = "STATIC"
      }
    }
    enabled           = true
    propagation_types = ["PROCESS_GROUP_TO_SERVICE", "PROCESS_GROUP_TO_HOST"]
    type              = "PROCESS_GROUP"
  }
}

