resource "dynatrace_log_sensitive_data_masking" "first-instance" {
  name    = "#name#"
  enabled = true
  scope   = "environment"
  masking {
    type       = "SHA1"
    expression = "FOO"
  }
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = ["jlkjk"]
    }
  }
}

resource "dynatrace_log_sensitive_data_masking" "second-instance" {
  name    = "#name#-second"
  enabled = true
  scope   = "environment"
  masking {
    type       = "SHA1"
    expression = "FOO"
  }
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = ["jlkjk-second"]
    }
  }
  insert_after = dynatrace_log_sensitive_data_masking.first-instance.id
}
