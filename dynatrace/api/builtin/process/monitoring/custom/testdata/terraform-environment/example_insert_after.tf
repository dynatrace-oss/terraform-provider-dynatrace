resource "dynatrace_process_monitoring_rule" "first-instance" {
  enabled = true
  mode    = "MONITORING_OFF"
  condition {
    item     = "APACHE_CONFIG_PATH"
    operator = "STARTS"
    value    = "foo-bar-x"
  }
}

resource "dynatrace_process_monitoring_rule" "second-instance" {
  enabled = true
  mode    = "MONITORING_OFF"
  condition {
    item     = "APACHE_CONFIG_PATH"
    operator = "STARTS"
    value    = "foo-bar-x-2"
  }
  insert_after = dynatrace_process_monitoring_rule.first-instance.id
}
