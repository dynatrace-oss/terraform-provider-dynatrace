resource "dynatrace_process_monitoring_rule" "test" {
  enabled = true
  mode    = "MONITORING_OFF"
  condition {
    item     = "APACHE_CONFIG_PATH"
    operator = "STARTS"
    value    = "foo-bar-x"
  }
}
