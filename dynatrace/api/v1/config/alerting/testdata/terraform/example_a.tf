resource "dynatrace_alerting_profile" "#name#" {
  display_name = "#name#"
  mz_id = ""
  rules {
    tag_filter {
      include_mode = "INCLUDE_ALL"
      tag_filters {
        context = "CONTEXTLESS"
        key = "EnvironmentA"
        value = "production"
      }
      tag_filters {
        context = "CONTEXTLESS"
        key = "Team"
        value = "test"
      }
    }
    delay_in_minutes = 0
    severity_level   = "AVAILABILITY"
  }
  rules {
    tag_filter {
      include_mode = "INCLUDE_ALL"
      tag_filters {
        context = "CONTEXTLESS"
        key = "EnvironmentB"
        value = "production"
      }
      tag_filters {
        context = "CONTEXTLESS"
        key = "Team"
        value = "test"
      }
    }
    delay_in_minutes = 0
    severity_level   = "CUSTOM_ALERT"
  }
  rules {
    tag_filter {
      include_mode = "INCLUDE_ALL"
      tag_filters {
        context = "CONTEXTLESS"
        key = "EnvironmentC"
        value = "production"
      }
      tag_filters {
        context = "CONTEXTLESS"
        key = "Team"
        value = "test"
      }
    }
    delay_in_minutes = 0
    severity_level   = "ERROR"
  }
  rules {
    tag_filter {
      include_mode = "INCLUDE_ALL"
      tag_filters {
        context = "CONTEXTLESS"
        key = "EnvironmentD"
        value = "production"
      }
      tag_filters {
        context = "CONTEXTLESS"
        key = "Team"
        value = "test"
      }
    }
    delay_in_minutes = 0
    severity_level   = "MONITORING_UNAVAILABLE"
  }
  rules {
    tag_filter {
      include_mode = "INCLUDE_ALL"
      tag_filters {
        context = "CONTEXTLESS"
        key = "EnvironmentE"
        value = "production"
      }
      tag_filters {
        context = "CONTEXTLESS"
        key = "Team"
        value = "test"
      }
    }
    delay_in_minutes = 0
    severity_level   = "PERFORMANCE"
  }
  rules {
    tag_filter {
      include_mode = "INCLUDE_ALL"
      tag_filters {
        context = "CONTEXTLESS"
        key = "EnvironmentF"
        value = "production"
      }
      tag_filters {
        context = "CONTEXTLESS"
        key = "Team"
        value = "test"
      }
    }
    delay_in_minutes = 0
    severity_level   = "RESOURCE_CONTENTION"
  }
}
