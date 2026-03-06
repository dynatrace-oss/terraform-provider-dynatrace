resource "dynatrace_request_attribute" "attribute" {
  name         = "#name#"
  enabled      = true
  aggregation  = "FIRST"
  confidential = false

  data_type                  = "INTEGER"
  normalization              = "ORIGINAL"
  skip_personal_data_masking = false
  data_sources {
    enabled = true
    source = "SERVER_VARIABLE"
    server_variable_technology = "ASP_NET"
    parameter_name             = "param"
  }
}

resource "dynatrace_management_zone_v2" "mzone" {
  name = "#name#"
  rules {
    rule {
      type            = "ME"
      enabled         = true
      entity_selector = ""
      attribute_rule {
        entity_type = "CLOUD_APPLICATION_NAMESPACE"
        attribute_conditions {
          condition {
            case_sensitive = false
            key            = "KUBERNETES_CLUSTER_NAME"
            operator       = "EQUALS"
            string_value   = "extensions"
          }
        }
      }
    }
  }
}

resource "time_sleep" "wait_for_request_attribute" {
  depends_on = [dynatrace_request_attribute.attribute]
  create_duration = "10s"
}

resource "dynatrace_calculated_service_metric" "metric" {
  depends_on = [time_sleep.wait_for_request_attribute]
  name             = "#name#"
  enabled          = true
  management_zones = [dynatrace_management_zone_v2.mzone.name]
  metric_key       = "calc:service.#name#"
  unit             = "MILLI_SECOND_PER_MINUTE"
  conditions {
    condition {
      attribute = "HTTP_REQUEST_METHOD"
      comparison {
        negate = false
        http_method {
          operator = "EQUALS_ANY_OF"
          values   = ["POST", "GET"]
        }
      }
    }
  }
  metric_definition {
    metric            = "REQUEST_ATTRIBUTE"
    request_attribute = dynatrace_request_attribute.attribute.name
  }
}
