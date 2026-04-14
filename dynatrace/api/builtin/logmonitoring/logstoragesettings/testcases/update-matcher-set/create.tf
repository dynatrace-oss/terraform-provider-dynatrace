resource "dynatrace_log_storage" "storage" {
  name            = "#name#"
  enabled         = false
  scope           = "HOST_GROUP-0000000000000000"
  send_to_storage = false
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = [ "TerraformTest" ]
    }
    matcher {
      attribute = "k8s.container.name"
      operator  = "MATCHES"
      values    = [ "TerraformTest" ]
    }
  }
}
