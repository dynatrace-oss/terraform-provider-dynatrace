resource "dynatrace_failure_detection_rule_sets" "#name#" {
  enabled = false
  scope   = "environment"
  ruleset {
    description  = "This is a sample description"
    condition    = "matchesValue(k8s.cluster.name,\"#name#\")"
    ruleset_name = "#name#"
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