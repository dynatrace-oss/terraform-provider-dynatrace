resource "dynatrace_log_sensitive_data_masking" "#name#" {
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
      values    = ["jlkjk"]
    }
  }
}
