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
}
