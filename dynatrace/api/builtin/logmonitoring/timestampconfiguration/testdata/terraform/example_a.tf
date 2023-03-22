resource "dynatrace_log_timestamp" "#name#" {
  enabled               = false
  config_item_title = "#name#"
  date_time_pattern = "%m/%d/%Y %I:%M:%S %p"
  scope                 = "environment"
  timezone              = "America/Detroit"
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = [ "terraform" ]
    }
  }
}