resource "dynatrace_request_attribute" "#name#" {
  name         = "#name#"
  enabled      = true
  aggregation  = "FIRST"
  confidential = false

  data_type                  = "STRING"
  normalization              = "ORIGINAL"
  skip_personal_data_masking = false
  data_sources {
    enabled    = true
    source     = "METHOD_PARAM"
    technology = "DOTNET"
    methods {
      capture = "CLASS_NAME"
      method {
        argument_types = ["!0", "System.Func`2\u003c!0,System.Threading.Tasks.Task\u003e"]
        class_name     = "NServiceBus.Pipeline.Behavior`1"
        method_name    = "Invoke"
        return_type    = "System.Threading.Tasks.Task"
        visibility     = "PUBLIC"
      }
    }
    value_processing {
      split_at              = "t"
      trim                  = true
      value_extractor_regex = "s(.*+)"
      extract_substring {
        delimiter = "h"
        position  = "BEFORE"
      }
      value_condition {
        negate   = false
        operator = "ENDS_WITH"
        value    = "gh"
      }
    }
  }
  data_sources {
    enabled    = true
    source     = "METHOD_PARAM"
    technology = "DOTNET"
    methods {
      capture = "CLASS_NAME"
      method {
        argument_types = ["!0", "System.Func`1\u003cSystem.Threading.Tasks.Task\u003e"]
        class_name     = "NServiceBus.Pipeline.Behavior`1"
        method_name    = "Invoke"
        return_type    = "System.Threading.Tasks.Task"
        visibility     = "PUBLIC"
      }
    }
    value_processing {
      split_at              = "t"
      trim                  = true
      value_extractor_regex = "s(.*+)"
      value_condition {
        operator = "ENDS_WITH"
        value    = "gh"
      }
    }
  }
  data_sources {
    enabled    = false
    source     = "METHOD_PARAM"
    technology = "DOTNET"
    methods {
      capture = "CLASS_NAME"
      method {
        argument_types = ["!0", "System.Func`1\u003cSystem.Threading.Tasks.Task\u003e"]
        class_name     = "NServiceBus.Pipeline.Behavior`1"
        method_name    = "Invoke"
        return_type    = "System.Threading.Tasks.Task"
        visibility     = "PUBLIC"
      }
    }
  }
}
