resource "dynatrace_log_timestamp" "first-instance" {
  enabled           = false
  config_item_title = "#name#"
  date_time_pattern = "%m/%d/%Y %I:%M:%S %p"
  scope             = "environment"
  timezone          = "America/Detroit"
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = ["terraform"]
    }
  }
}

resource "dynatrace_log_timestamp" "second-instance" {
  enabled           = false
  config_item_title = "#name#-second"
  date_time_pattern = "%m/%d/%Y %I:%M:%S %p"
  scope             = "environment"
  timezone          = "America/Detroit"
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = ["terraform-second"]
    }
  }
  insert_after = dynatrace_log_timestamp.first-instance.id
}
