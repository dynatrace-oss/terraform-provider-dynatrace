resource "dynatrace_request_attribute" "accept_ranges" {
  name         = "Accept-Ranges-#name#"
  enabled      = true
  aggregation  = "FIRST"
  confidential = false

  data_type                  = "STRING"
  normalization              = "ORIGINAL"
  skip_personal_data_masking = false
  data_sources {
    enabled = true
    source = "RESPONSE_HEADER"
    parameter_name                 = "Accept-Ranges"
    capturing_and_storage_location = "CAPTURE_ON_CLIENT_STORE_ON_SERVER"
  }
}

resource "dynatrace_request_attribute" "behavior_class" {
  name = "BehaviorClass-#name#"
  enabled = true
  aggregation = "FIRST"
  confidential = false
  normalization = "ORIGINAL"
  skip_personal_data_masking = false
  data_type = "STRING"
  data_sources {
    enabled = false
    source  = "METHOD_PARAM"
    technology = "DOTNET"
    value_processing {
      value_condition {
        operator = "ENDS_WITH"
        negate   = false
        value    = "gh"
      }
      value_extractor_regex = "s(.*+)"
      split_at              = "t"
      trim                  = true
      extract_substring {
        position  = "BEFORE"
        delimiter = "h"
      }
    }
    methods {
      capture = "CLASS_NAME"
      method {
        visibility = "PUBLIC"
        class_name = "NServiceBus.Pipeline.Behavior`1"
        method_name = "Invoke"
        argument_types = [
          "!0",
          "System.Func`1<System.Threading.Tasks.Task>"
        ]
        return_type = "System.Threading.Tasks.Task"
      }
    }
  }
}

resource "time_sleep" "wait_for_request_attributes" {
  depends_on = [dynatrace_request_attribute.attribute, dynatrace_request_attribute.accept_ranges, dynatrace_request_attribute.behavior_class]
  create_duration = "10s"
}
