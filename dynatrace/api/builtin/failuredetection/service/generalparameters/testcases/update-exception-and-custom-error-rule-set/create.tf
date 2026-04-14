resource "dynatrace_service_failure" "failure" {
  enabled    = true
  service_id = "SERVICE-D892CFE7DFAB0D08"
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
      custom_error_rule {
        request_attribute = "00000000-0000-0000-0000-000000000001"
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
      custom_handled_exception {
        class_pattern   = "ClassPattern2"
        message_pattern = "ExceptionPattern"
      }
    }
    ignored_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
      custom_handled_exception {
        class_pattern   = "ClassPattern2"
        message_pattern = "ExceptionPattern"
      }
    }
    success_forcing_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPattern"
        message_pattern = "ExceptionPattern"
      }
      custom_handled_exception {
        class_pattern   = "ClassPattern2"
        message_pattern = "ExceptionPattern"
      }
    }
  }
}
