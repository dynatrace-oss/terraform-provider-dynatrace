data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

resource "dynatrace_network_monitor" "DNS_Test" {
  name          = "#name#"
  description   = "This is an example DNS test"
  type          = "MULTI_PROTOCOL"
  enabled       = false
  frequency_min = 15
  locations     = [data.dynatrace_synthetic_location.location.id]
  outage_handling {
    global_consecutive_outage_count_threshold = 1
    global_outages                            = true
  }
  performance_thresholds {
    enabled = false
    thresholds {
      threshold {
        aggregation = "AVG"
        threshold   = 100
        step_index  = 0
      }
      # removed
    }
  }
  steps {
    step {
      name         = "DNS Test"
      request_type = "DNS"
      target_list  = ["google.com", "yahoo.com"]
      properties = {
        "DNS_RECORD_TYPES"  = "A"
        "EXECUTION_TIMEOUT" = "PT2S"
      }
      constraints {
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "90"
            "operator" = ">="
          }
        }
        # removed
      }
      request_configurations {
        request_configuration {
          constraints {
            constraint {
              type = "DNS_STATUS_CODE"
              properties = {
                "operator"   = "="
                "statusCode" = "0"
              }
            }
          }
        }
        # removed
      }
    }
    step {
      name         = "DNS Test 2"
      request_type = "DNS"
      target_list  = ["google.com", "yahoo.com"]
      properties = {
        "DNS_RECORD_TYPES"  = "A"
        "EXECUTION_TIMEOUT" = "PT2S"
      }
      constraints {
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "90"
            "operator" = ">="
          }
        }
      }
      request_configurations {
        request_configuration {
          constraints {
            constraint {
              type = "DNS_STATUS_CODE"
              properties = {
                "operator"   = "="
                "statusCode" = "0"
              }
            }
          }
        }
      }
    }
  }
  tags {
    tag {
      context = "CONTEXTLESS"
      key     = "Key1"
      source  = "USER"
      value   = "Value1"
    }
    # removed
  }
}
