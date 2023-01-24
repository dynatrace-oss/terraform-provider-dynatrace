resource "dynatrace_management_zone" "#name#" {
  name = "#name#"
  rules {
    type              = "PROCESS_GROUP"
    enabled           = true
    propagation_types = ["PROCESS_GROUP_TO_HOST", "PROCESS_GROUP_TO_SERVICE"]
    conditions {
      key {
        type      = "STATIC"
        attribute = "PROCESS_GROUP_TAGS"
      }
      tag {
        # negate = false 
        operator = "TAG_KEY_EQUALS"
        value {
          context = "CONTEXTLESS"
          key     = "Environment"
        }
      }
    }
    conditions {
      key {
        type      = "STATIC"
        attribute = "PROCESS_GROUP_TAGS"
      }
      tag {
        # negate = false 
        operator = "TAG_KEY_EQUALS"
        value {
          context = "CONTEXTLESS"
          key     = "Team"
        }
      }
    }
  }
  rules {
    type              = "PROCESS_GROUP"
    enabled           = true
    propagation_types = ["PROCESS_GROUP_TO_HOST", "PROCESS_GROUP_TO_SERVICE"]
    conditions {
      key {
        type      = "STATIC"
        attribute = "PROCESS_GROUP_TAGS"
      }
      tag {
        # negate = false 
        operator = "TAG_KEY_EQUALS"
        value {
          context = "CONTEXTLESS"
          key     = "EnvironmentX"
        }
      }
    }
    conditions {
      key {
        type      = "STATIC"
        attribute = "PROCESS_GROUP_TAGS"
      }
      tag {
        # negate = false 
        operator = "TAG_KEY_EQUALS"
        value {
          context = "CONTEXTLESS"
          key     = "TeamX"
        }
      }
    }
  }
}
