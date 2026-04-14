resource "dynatrace_failure_detection_parameters" "params" {
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
        request_attribute = "00000000-0000-0000-0000-000000000000"
        condition {
          compare_operation_type = "STRING_EXISTS"
        }
      }
      # update => re-create due to set-hash change
      custom_error_rule {
        request_attribute = "00000000-0000-0000-0000-000000000002"
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
      # update => re-create due to set-hash change
      custom_handled_exception {
        class_pattern   = "ClassPatternEdit"
        message_pattern = "ExceptionPatternEdit"
      }
    }
    ignored_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
      # update => re-create due to set-hash change
      custom_handled_exception {
        class_pattern   = "ClassPatternEdit"
        message_pattern = "ExceptionPatternEdit"
      }
    }
    success_forcing_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
      # update => re-create due to set-hash change
      custom_handled_exception {
        class_pattern   = "ClassPatternEdit"
        message_pattern = "ExceptionPatternEdit"
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
