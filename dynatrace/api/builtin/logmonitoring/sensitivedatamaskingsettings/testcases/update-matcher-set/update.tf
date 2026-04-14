resource "dynatrace_log_sensitive_data_masking" "masking" {
  name    = "#name#"
  enabled = true
  scope   = "environment"
  masking {
    type = "SHA1"
    expression  = "FOO"
  }
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
