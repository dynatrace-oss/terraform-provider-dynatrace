resource "dynatrace_network_monitor" "TCP_Test" {
  name          = "TCP Test"
  description   = "This is an example TCP Test"
  type          = "MULTI_PROTOCOL"
  enabled       = false
  frequency_min = 15
  locations     = [ "SYNTHETIC_LOCATION-39F97465BE7BF644" ]
  outage_handling {
    global_consecutive_outage_count_threshold = 1
    global_outages                            = true
  }
  steps {
    step {
      name         = "TCP Test"
      request_type = "TCP"
      target_list  = [ "8.8.8.8", "8.8.4.4" ]
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
    step {
      name         = "TCP Test 2"
      request_type = "TCP"
      target_list  = [ "8.8.8.8", "8.8.4.4" ]
      properties = {
        "EXECUTION_TIMEOUT" = "PT1S"
        "TCP_PORT_RANGES"   = "53"
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
  tags {
    tag {
      context = "CONTEXTLESS"
      key     = "KeyA"
      source  = "USER"
      value   = "ValueA"
    }
    tag {
      context = "CONTEXTLESS"
      key     = "KeyB"
      source  = "USER"
      value   = "ValueB"
    }
    tag {
      context = "CONTEXTLESS"
      key     = "KeyC"
      source  = "USER"
    }
  }
}