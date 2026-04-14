resource "dynatrace_failure_detection_rule_sets" "rule-set" {
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
      # update => re-create due to set-hash change
      fail_on_custom_rule {
        enabled       = true
        dql_condition = "matchesValue(terraform.test.1, \"#name#-edit\")"
        rule_name     = "Test edit"
      }
      fail_on_custom_rule {
        enabled       = true
        dql_condition = "matchesValue(terraform.test.1, \"#name#-new\")"
        rule_name     = "Test new"
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
        # update => re-create due to set-hash change
        ignored_exception {
          type    = "Terraform1-edit"
          enabled = true
          message = "#name#-edit"
        }
        ignored_exception {
          type    = "Terraform1-new"
          enabled = true
          message = "#name#-new"
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
        # update => re-create due to set-hash change
        ignored_exception {
          type    = "TerraformEdit"
          enabled = true
          message = "#name#-edit"
        }
        ignored_exception {
          type    = "TerraformNew"
          enabled = true
          message = "#name#-new"
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
        # update => re-create due to set-hash change
        fail_on_custom_rule {
          enabled       = true
          dql_condition = "matchesValue(terraform.test.1, \"#name#-edit\")"
          rule_name     = "Test edit"
        }
      }
    }
  }
}
