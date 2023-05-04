resource "dynatrace_business_events_oneagent" "#name#" {
  enabled   = true
  rule_name = "#name#"
  scope     = "environment"
  event {
    category {
      source_type = "request.path"
    }
    data {
      event_data_field_complex {
        name = "rsbody"
        source {
          path        = "*"
          source_type = "response.body"
        }
      }
      event_data_field_complex {
        name = "req content-type"
        source {
          path        = "content-type"
          source_type = "request.headers"
        }
      }
      event_data_field_complex {
        name = "technology"
        source {
          source      = "java"
          source_type = "constant.string"
        }
      }
      event_data_field_complex {
        name = "rqbody"
        source {
          path        = "*"
          source_type = "request.body"
        }
      }
      event_data_field_complex {
        name = "res content-type"
        source {
          path        = "content-type"
          source_type = "response.headers"
        }
      }
    }
    provider {
      source      = "www.easytravel.com"
      source_type = "constant.string"
    }
    type {
      source      = "com.easytravel.search-journey"
      source_type = "constant.string"
    }
  }
  triggers {
    trigger {
      type           = "STARTS_WITH"
      case_sensitive = false
      value          = "/easytravel/rest/journeys"
      source {
        data_source = "request.path"
      }
    }
  }
}
