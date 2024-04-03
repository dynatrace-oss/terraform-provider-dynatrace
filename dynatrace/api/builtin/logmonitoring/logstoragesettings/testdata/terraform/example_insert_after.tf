resource "dynatrace_log_storage" "first-instance" {
  name            = "#name#"
  enabled         = false
  scope           = "HOST_GROUP-1234567890000000"
  send_to_storage = false
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = ["TerraformTest"]
    }
  }
}

resource "dynatrace_log_storage" "second-instance" {
  name            = "#name#-second"
  enabled         = false
  scope           = "HOST_GROUP-1234567890000000"
  send_to_storage = false
  matchers {
    matcher {
      attribute = "container.name"
      operator  = "MATCHES"
      values    = ["TerraformTest-second"]
    }
  }
  insert_after = dynatrace_log_storage.first-instance.id
}
