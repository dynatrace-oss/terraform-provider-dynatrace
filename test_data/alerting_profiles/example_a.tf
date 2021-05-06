resource "dynatrace_alerting_profile" "#name#" {
  mz_id = ""
  rules {
    tag_filter {
      include_mode = "NONE"
    }
    delay_in_minutes = 0
    severity_level   = "AVAILABILITY"
  }
  rules {
    tag_filter {
      include_mode = "NONE"
    }
    delay_in_minutes = 0
    severity_level   = "CUSTOM_ALERT"
  }
  rules {
    tag_filter {
      include_mode = "NONE"
    }
    delay_in_minutes = 0
    severity_level   = "ERROR"
  }
  rules {
    tag_filter {
      include_mode = "NONE"
    }
    delay_in_minutes = 0
    severity_level   = "MONITORING_UNAVAILABLE"
  }
  rules {
    tag_filter {
      include_mode = "NONE"
    }
    delay_in_minutes = 0
    severity_level   = "PERFORMANCE"
  }
  rules {
    tag_filter {
      include_mode = "NONE"
    }
    delay_in_minutes = 0
    severity_level   = "RESOURCE_CONTENTION"
  }
  display_name = "#name#"
}
