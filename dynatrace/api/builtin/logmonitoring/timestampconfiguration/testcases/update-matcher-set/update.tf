resource "dynatrace_log_timestamp" "timestamp" {
  enabled               = false
  config_item_title = "#name#"
  date_time_pattern = "%m/%d/%Y %I:%M:%S %p"
  scope                 = "environment"
  timezone              = "America/Detroit"
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = [ "TerraformTest" ]
    }
    # update => re-create due to set-hash change
    matcher {
      attribute = "k8s.namespace.name"
      operator  = "MATCHES"
      values    = [ "TerraformTestEdit" ]
    }
  }
}
