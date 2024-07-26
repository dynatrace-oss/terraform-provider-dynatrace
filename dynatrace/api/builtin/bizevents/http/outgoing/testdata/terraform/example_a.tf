resource "dynatrace_business_events_oneagent_outgoing" "#name#" {
  enabled      = true
  rule_name    = "#name#"
  scope        = "environment"
  event {
    category {
      source      = "Category 1"
      source_type = "constant.string"
    }
    data {
      event_data_field_complex {
        name = "Field 1"
        source {
          path        = "Path 1"
          source_type = "request.body"
        }
      }
    }
    provider {
      source      = "Provider 1"
      source_type = "constant.string"
    }
    type {
      source      = "Type 1"
      source_type = "constant.string"
    }
  }
  triggers {
    trigger {
      type           = "EQUALS"
      case_sensitive = false
      value          = "Terraform"
      source {
        data_source = "request.path"
      }
    }
  }
}