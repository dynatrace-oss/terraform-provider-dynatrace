data "dynatrace_synthetic_location" "location" {
  name = "Location"
}

resource "dynatrace_network_monitor" "TCP_Test" {
  name          = "#name#"
  description   = "This is an example TCP test"
  type          = "MULTI_PROTOCOL"
  enabled       = false
  frequency_min = 15
  locations     = [data.dynatrace_synthetic_location.location.id]
  outage_handling {
    global_consecutive_outage_count_threshold = 1
    global_outages                            = true
  }
  steps {
    step {
      name         = "TCP Test"
      request_type = "TCP"
      target_list  = ["8.8.8.8", "8.8.4.4"]
      properties = {
        "TCP_PORT_RANGES"   = "53"
        "EXECUTION_TIMEOUT" = "PT1S"
      }
      constraints {
        constraint {
          type = "SUCCESS_RATE_PERCENT"
          properties = {
            "value"    = "50"
            "operator" = ">="
          }
        }
      }
    }
  }
  # The API assigns source `USER` and context `CONTEXTLESS` to every tag it stores
  tags {
    tag {
      key   = "somekey"
      value = "somevalue"
    }
    tag {
      key = "test-key"
    }
    # removed
  }
}
