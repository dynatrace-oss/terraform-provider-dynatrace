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
    # update => re-create due to set-hash change
    matcher {
      attribute = "k8s.namespace.name"
      operator  = "MATCHES"
      values    = [ "TerraformTestEdit" ]
    }
  }
}
