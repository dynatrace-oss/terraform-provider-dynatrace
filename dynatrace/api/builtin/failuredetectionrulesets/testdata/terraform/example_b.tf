resource "dynatrace_failure_detection_rule_sets" "first-instance" {
  enabled = true
  scope   = "environment"
  ruleset {
    description  = "This is a sample description"
    condition    = "matchesValue(k8s.cluster.name,\"#name#\")"
    ruleset_name = "#name#-1"
    fail_on_custom_rules {
      fail_on_custom_rule {
        enabled       = true
        dql_condition = "matchesValue(terraform.test.2, \"#name#\")"
        rule_name     = "Test 2"
      }
      fail_on_custom_rule {
        enabled       = true
        dql_condition = "matchesValue(terraform.test.1, \"#name#\")"
        rule_name     = "Test 1"
      }
    }
    fail_on_exceptions {
      enabled = true
      ignored_exceptions {
        ignored_exception {
          type    = "Terraform2"
          enabled = true
          message = "#name#"
        }
        ignored_exception {
          type    = "Terraform1"
          enabled = true
          message = "#name#"
        }
      }
    }
    fail_on_grpc_status_codes {
      status_codes = "2,4,12,13,14,15"
    }
    fail_on_http_response_status_codes {
      status_codes = "500-599"
    }
    fail_on_span_status_error {
      enabled = true
    }
    overrides {
      force_success_on_exceptions {
        ignored_exception {
          type    = "Terraform2"
          enabled = true
          message = "#name#"
        }
        ignored_exception {
          type    = "Terraform1"
          enabled = true
          message = "#name#"
        }
      }
      force_success_on_grpc_response_status_codes {
        status_codes = "20"
      }
      force_success_on_http_response_status_codes {
        status_codes = "555"
      }
      force_success_on_span_status_ok {
        enabled = true
      }
      force_success_with_custom_rules {
        fail_on_custom_rule {
          enabled       = true
          dql_condition = "matchesValue(terraform.test.2, \"#name#\")"
          rule_name     = "Test 2"
        }
        fail_on_custom_rule {
          enabled       = true
          dql_condition = "matchesValue(terraform.test.1, \"#name#\")"
          rule_name     = "Test 1"
        }
      }
    }
  }
}

resource "dynatrace_failure_detection_rule_sets" "second-instance" {
  insert_after = dynatrace_failure_detection_rule_sets.first-instance.id
  enabled = false
  scope   = "environment"
  ruleset {
    description  = "This is a sample description"
    condition    = "matchesValue(k8s.cluster.name,\"#name#\")"
    ruleset_name = "#name#-2"
    fail_on_exceptions {
      enabled = true
    }
    fail_on_grpc_status_codes {
      status_codes = "2,4,12,13,14,15"
    }
    fail_on_http_response_status_codes {
      status_codes = "500-599"
    }
    fail_on_span_status_error {
      enabled = true
    }
    overrides {
      force_success_on_span_status_ok {
        enabled = false
      }
    }
  }
}