resource "dynatrace_log_custom_source" "#name#" {
  name    = "#name#"
  enabled = false
  scope   = "HOST_GROUP-1234567890000000"
  custom_log_source {
    accept_binary = true
    type   = "LOG_PATH_PATTERN"
    values = [ "/terraform" ]
  }
}