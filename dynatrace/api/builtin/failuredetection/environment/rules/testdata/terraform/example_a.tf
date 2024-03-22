resource "dynatrace_failure_detection_rules" "#name#" {
  name         = ""
  enabled      = true
  parameter_id = "${dynatrace_failure_detection_parameters.parameter.id}"
  conditions {
    condition {
      attribute = "SERVICE_NAME"
      predicate {
        case_sensitive = true
        predicate_type = "STRING_EQUALS"
        text_values    = [ "Terraform" ]
      }
    }
  }
}

resource "dynatrace_failure_detection_parameters" "parameter" {
  name        = "#name#"
  description = "Created by Terraform"
  broken_links {
    http_404_not_found_failures = false
  }
  exception_rules {
    ignore_all_exceptions         = false
    ignore_span_failure_detection = true
    custom_error_rules {
      custom_error_rule {
        request_attribute = "195b205c-5c01-4563-b29b-e33caf24ec7d"
        condition {
          compare_operation_type = "STRING_EXISTS"
        }
      }
    }
    custom_handled_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
    }
    ignored_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
    }
    success_forcing_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
    }
  }
  http_response_codes {
    client_side_errors                        = "400-599"
    fail_on_missing_response_code_client_side = false
    fail_on_missing_response_code_server_side = true
    server_side_errors                        = "500-599"
  }
}
