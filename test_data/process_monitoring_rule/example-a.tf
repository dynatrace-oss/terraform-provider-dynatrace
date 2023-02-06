resource "dynatrace_process_monitoring_rule" "test" {
  enabled = true
  mode    = "MONITORING_OFF"
  host_group_id   = "HOST_GROUP-0000000000000000"
  condition {
    item     = "APACHE_CONFIG_PATH"
    operator = "STARTS"
    value    = "foo-bar-x"
  }
}