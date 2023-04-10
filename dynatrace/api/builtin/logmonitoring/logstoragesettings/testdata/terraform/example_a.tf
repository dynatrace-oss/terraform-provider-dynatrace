resource "dynatrace_log_storage" "#name#" {
  name            = "#name#"
  enabled         = false
  scope           = "HOST_GROUP-1234567890000000"
  send_to_storage = false
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = [ "TerraformTest" ]
    }
  }
}