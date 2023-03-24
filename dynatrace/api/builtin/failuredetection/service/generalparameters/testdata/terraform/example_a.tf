resource "dynatrace_service_failure" "#name#" {
  enabled    = true
  service_id = "SERVICE-D892CFE7DFAB0D08"
  exception_rules {
    ignore_all_exceptions         = false
    ignore_span_failure_detection = true
    custom_error_rules {
      custom_error_rule {
        request_attribute = "00000000-0000-0000-0000-000000000000"
        condition {
          compare_operation_type = "STARTS_WITH"
          case_sensitive = false
          text_value = "terraform"
        }
      }
    }
    custom_handled_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPatternExample"
        message_pattern = "ExceptionMessagePatternExample"
      }
    }
    ignored_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPatternExample"
        message_pattern = "ExceptionMessagePatternExample"
      }
    }
    success_forcing_exceptions {
      custom_handled_exception {
        class_pattern   = "ClassPatternExample"
        message_pattern = "ExceptionMessagePatternExample"
      }
    }
  }
}